package output

import (
	"fmt"
	"os"
	"strings"

	"pig/db"
	"pig/file"
	"pig/migration"

	"github.com/olekukonko/tablewriter"
)

// PrintSmallDBTable - print small table
func PrintSmallDBTable(records []*db.RecShort) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Timestamp", "Name", "Filename", "Requirements", "Is applied"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
	)

	for _, record := range records {
		table.Append([]string{
			fmt.Sprintf("%v", record.ID),
			record.TStamp.Format("2006-01-02 15:04:05"),
			record.Name,
			record.Filename,
			strings.Join(record.Requirements, ","),
			fmt.Sprintf("%v", record.Applied),
		})
	}
	table.Render() // Send output
}

// PrintSmallFileTable - print small table
func PrintSmallFileTable(records []*file.RecShort) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Filename"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{},
	)

	for _, record := range records {
		table.Append([]string{
			record.TStamp.Format("2006-01-02 15:04:05"),
			record.Filename,
		})
	}
	table.Render() // Send output
}

// PrintFullFileTable - print small table
func PrintFullFileTable(records []*file.RecFull) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "Name", "Filename", "Requirements"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgCyanColor},
	)

	for _, record := range records {
		table.Append([]string{
			record.TStamp.Format("2006-01-02 15:04:05"),
			record.Name,
			record.Filename,
			strings.Join(record.Requirements, ","),
		})
	}
	table.Render() // Send output
}

// PrintMigrationsTable - print small table
func PrintMigrationsTable(records *migration.Migrations) {
	table := tablewriter.NewWriter(os.Stdout)
	// table.SetRowLine(true)
	table.SetHeader([]string{"ID", "Timestamp", "Name", "Filename", "Requirements", "Is applied", "Pending"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	)

	for _, record := range records.Items {
		recordID := "?"
		if record.ID != 0 {
			recordID = fmt.Sprintf("%v", record.ID)
		}

		recordRequirements := []string{}
		for _, req := range record.Requirements {
			recordRequirements = append(recordRequirements, req.Name)
		}

		table.Append([]string{
			recordID,
			record.TStamp.Format("2006-01-02 15:04:05"),
			record.Name,
			record.Filename,
			strings.Join(recordRequirements, ", "),
			fmt.Sprintf("%v", record.Applied),
			fmt.Sprintf("%v", record.Pending),
		})
	}
	table.Render() // Send output
}

// PrintHistoryTable - print history
func PrintHistoryTable(records []*db.RecHistory) {
	table := tablewriter.NewWriter(os.Stdout)
	applied := ""

	table.SetHeader([]string{"ID", "When", "Name", "Is applied", "Note", "GIT info", "Filename"})

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{},
		tablewriter.Colors{tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
	)

	for _, record := range records {
		recordID := "?"
		if record.ID != 0 {
			recordID = fmt.Sprintf("%v", record.ID)
		}

		if record.Applied {
			applied = "+"
		} else {
			applied = "-"
		}

		table.Append([]string{
			recordID,
			record.When.Format("2006-01-02 15:04:05"),
			record.Name,
			applied,
			stringPointerToString(record.Note),
			stringPointerToString(record.GITinfo),
			record.Filename,
		})
	}
	table.Render() // Send output
}
