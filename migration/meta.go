package migration

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/bfg-dev/pig/db"
	"github.com/bfg-dev/pig/file"
)

// Meta migration metadata
type Meta struct {
	ID           uint64
	Requirements []*Meta
	Children     []*Meta
	Name         string
	Applied      bool
	Pending      bool
	TStamp       time.Time
	Note         string
	GITinfo      string
	Filename     string
	DBShortRec   *db.RecShort  // can not be nil on DOWN
	DBFullRec    *db.RecFull   // can not be nil on DOWN
	FileFullRec  *file.RecFull // can be nil on DOWN
}

// Statements - sql statements from SQLData
type Statements struct {
	Lines         []string
	Transactional bool
}

// checkLoop - check for loop in requirements
func (mm *Meta) checkLoop(total int, cur int) error {
	if cur > total {
		return &LoopError{LastMigration: mm.Name}
	}

	for _, m := range mm.Requirements {
		if err := m.checkLoop(total, cur+1); err != nil {
			return err
		}
	}

	return nil
}

// getUnappliedRequirements - return unapplied requirements
func (mm *Meta) getUnappliedRequirements() []*Meta {
	var ans []*Meta

	for _, r := range mm.Requirements {
		if !r.Applied {
			ans = append(ans, r)
		}
	}

	return ans
}

// getAppliedChildren - return applied requirements
func (mm *Meta) getAppliedChildren() []*Meta {
	var ans []*Meta

	for _, r := range mm.Children {
		if r.Applied {
			ans = append(ans, r)
		}
	}

	return ans
}

// getUpPlan
func (mm *Meta) getUpPlan() *Migrations {
	var ans Migrations

	if mm.Applied {
		return &ans
	}

	reqs := mm.getUnappliedRequirements()

	if len(reqs) == 0 {
		ans.Items = append(ans.Items, mm)
		return &ans
	}

	for _, r := range reqs {
		ans.Items = append(ans.Items, r.GetUpPlan().Items...)
	}

	ans.Items = append(ans.Items, mm)

	return &ans
}

// getDownPlan
func (mm *Meta) getDownPlan() *Migrations {
	var ans Migrations

	if !mm.Applied {
		return &ans
	}

	children := mm.getAppliedChildren()

	if len(children) == 0 {
		ans.Items = append(ans.Items, mm)
		return &ans
	}

	for _, c := range children {
		ans.Items = append(ans.Items, c.getDownPlan().Items...)
	}

	ans.Items = append(ans.Items, mm)

	return &ans
}

func (mm *Meta) String() string {
	return mm.Name
}

// AddRequirement - add requirement to the migration
func (mm *Meta) AddRequirement(req *Meta) error {
	if req == nil {
		return &NullRequirement{Migration: mm}
	}

	for _, r := range mm.Requirements {
		if r.Name == req.Name {
			return &RequirementDuplicateError{Migration: mm, Requirement: req}
		}
	}
	mm.Requirements = append(mm.Requirements, req)

	for _, c := range req.Children {
		if c.Name == mm.Name {
			return &RequirementDuplicateError{Migration: req, Requirement: mm}
		}
	}
	req.Children = append(req.Children, mm)

	return nil
}

// GetUpPlan - return up plan
func (mm *Meta) GetUpPlan() *Migrations {
	return mm.getUpPlan().RemoveDuplicates()
}

// GetDownPlan - return up plan
func (mm *Meta) GetDownPlan() *Migrations {
	return mm.getDownPlan().RemoveDuplicates()
}

// Checks the line to see if the line has a statement-ending semicolon
// or if the line contains a double-dash comment.
func (mm *Meta) endsWithSemicolon(line string) bool {

	prev := ""
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, "--") {
			break
		}
		prev = word
	}

	return strings.HasSuffix(prev, ";")
}

// GetSQLStatements - get UP or DOWN sql statements (code from goose)
func (mm *Meta) GetSQLStatements(up bool) (*Statements, error) {
	var (
		buf     bytes.Buffer
		scanner *bufio.Scanner
		result  Statements
	)

	if up {
		if mm.FileFullRec == nil {
			return nil, fmt.Errorf("No sql in files for migration %v", mm.Name)
		}
		scanner = bufio.NewScanner(strings.NewReader(*mm.FileFullRec.SQLData))
	} else {
		if mm.DBFullRec == nil {
			return nil, fmt.Errorf("No sql in database for migration %v", mm.Name)
		}
		scanner = bufio.NewScanner(strings.NewReader(*mm.DBFullRec.SQLData))
	}

	// track the count of each section
	// so we can diagnose scripts with no annotations
	upSections := 0
	downSections := 0

	statementEnded := false
	ignoreSemicolons := false
	directionIsActive := false
	result.Transactional = true

	for scanner.Scan() {

		line := scanner.Text()

		// handle any goose-specific commands
		if strings.HasPrefix(line, file.CMDPrefix) {
			cmd := strings.TrimSpace(line[len(file.CMDPrefix):])
			switch cmd {
			case "Up":
				directionIsActive = (up == true)
				upSections++
				break

			case "Down":
				directionIsActive = (up == false)
				downSections++
				break

			case "StatementBegin":
				if directionIsActive {
					ignoreSemicolons = true
				}
				break

			case "StatementEnd":
				if directionIsActive {
					statementEnded = (ignoreSemicolons == true)
					ignoreSemicolons = false
				}
				break

			case "NO TRANSACTION":
				result.Transactional = false
				break
			}
		}

		if !directionIsActive {
			continue
		}

		if _, err := buf.WriteString(line + "\n"); err != nil {
			return nil, fmt.Errorf("io error: %v", err)
		}

		// Wrap up the two supported cases: 1) basic with semicolon; 2) psql statement
		// Lines that end with semicolon that are in a statement block
		// do not conclude statement.
		if (!ignoreSemicolons && mm.endsWithSemicolon(line)) || statementEnded {
			statementEnded = false
			result.Lines = append(result.Lines, buf.String())
			buf.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning migration: %v", err)
	}

	if upSections == 0 && downSections == 0 {
		return nil, fmt.Errorf("ERROR: no Up/Down annotations found, so no statements were executed")
	}

	return &result, nil
}
