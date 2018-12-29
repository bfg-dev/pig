package file

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// RecManager - file manager
type RecManager struct {
	dir string
}

// NewRecManager - returns new RecManager
func NewRecManager(dir string) *RecManager {
	return &RecManager{
		dir: dir,
	}
}

// GetAllShort - get all short file records
func (m *RecManager) GetAllShort() ([]*RecShort, error) {
	var result []*RecShort

	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s directory does not exists", m.dir)
	}

	sqlMigrationFiles, err := filepath.Glob(m.dir + "/**.sql")
	if err != nil {
		return nil, err
	}

	for _, file := range sqlMigrationFiles {
		fstat, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		result = append(result, &RecShort{
			Filename: filepath.Base(file),
			TStamp:   fstat.ModTime(),
		})
	}

	return result, nil
}

// parseFile - parse sql file into full record
func (m *RecManager) parseFile(rec *RecFull) error {
	sqlByteData, err := ioutil.ReadFile(filepath.Join(m.dir, rec.Filename))
	if err != nil {
		return err
	}

	sqlData := string(sqlByteData)
	rec.SQLData = &sqlData

	scanner := bufio.NewScanner(bytes.NewReader(sqlByteData))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, CMDPrefix) {
			option := strings.TrimSpace(line[len(CMDPrefix):])

			if strings.HasPrefix(option, optNamePrefix) {
				rec.Name = strings.TrimSpace(option[len(optNamePrefix):])
				if strings.Index(rec.Name, ",") != -1 {
					return fmt.Errorf("Name can not contain ',': %v", rec.Name)
				}
			}

			if strings.HasPrefix(option, optRequiremetsPrefix) {
				for _, req := range strings.Split(strings.TrimSpace(option[len(optRequiremetsPrefix):]), ",") {
					rec.Requirements = append(rec.Requirements, strings.TrimSpace(req))
				}
			}
		}
	}

	return nil
}

// GetFullFromShort - get full file records from short
func (m *RecManager) GetFullFromShort(shortRecords []*RecShort) ([]*RecFull, error) {
	var result []*RecFull

	for _, rec := range shortRecords {
		fullRec := &RecFull{}
		fullRec.RecShort = *rec
		if err := m.parseFile(fullRec); err != nil {
			return nil, err
		}
		if len(fullRec.Name) == 0 {
			fullRec.Name = strings.TrimSuffix(fullRec.Filename, filepath.Ext(fullRec.Filename))
		}
		result = append(result, fullRec)
	}

	return result, nil
}

// GetAllFull - get all full file records
func (m *RecManager) GetAllFull() ([]*RecFull, error) {
	var result []*RecFull

	shortRecords, err := m.GetAllShort()
	if err != nil {
		return nil, err
	}

	result, err = m.GetFullFromShort(shortRecords)
	return result, err
}

// FindByFilenameFull - find by filename in array
func (m *RecManager) FindByFilenameFull(records []*RecFull, filename string) *RecFull {
	for _, rec := range records {
		if rec.Filename == filename {
			return rec
		}
	}
	return nil
}

// RemoveFromListFull - remove item from short list
func (m *RecManager) RemoveFromListFull(records []*RecFull, removeItem *RecFull) []*RecFull {
	var result []*RecFull

	for _, rec := range records {
		if rec.Filename != removeItem.Filename {
			result = append(result, rec)
		}
	}

	return result
}
