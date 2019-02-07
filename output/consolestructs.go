package output

import (
	"fmt"
	"strings"

	"github.com/bfg-dev/pig/db"
	"github.com/bfg-dev/pig/migration"
)

// PrintHistoryStructs - print history
func PrintHistoryStructs(records []*db.RecHistory) {
	var applied string

	fmt.Println(gray("-----------------"))

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

		fmt.Println(gray("ID:        "), gray(recordID))
		fmt.Println(gray("When:      "), record.When.Format("2006-01-02 15:04:05"))
		fmt.Println(gray("Name:      "), hiblue(record.Name))
		fmt.Println(gray("Is applied:"), yellow(applied))
		fmt.Println(gray("Note:      "), stringPointerToString(record.Note))
		fmt.Println(gray("GIT info:  "), stringPointerToString(record.GITinfo))
		fmt.Println(gray("Filename:  "), record.Filename)
		fmt.Println(gray("-----------------"))

	}
}

// PrintMigrationsStructs - print small table
func PrintMigrationsStructs(records *migration.Migrations) {
	fmt.Println(gray("-----------------"))
	for _, record := range records.Items {
		recordID := "?"
		if record.ID != 0 {
			recordID = fmt.Sprintf("%v", record.ID)
		}

		recordRequirements := []string{}
		for _, req := range record.Requirements {
			recordRequirements = append(recordRequirements, req.Name)
		}

		fmt.Println(gray("ID:          "), gray(recordID))
		fmt.Println(gray("Timestamp:   "), record.TStamp.Format("2006-01-02 15:04:05"))
		fmt.Println(gray("Name:        "), hiblue(record.Name))
		fmt.Println(gray("Filename:    "), record.Filename)
		fmt.Println(gray("Requirements:"), cyan(strings.Join(recordRequirements, ", ")))
		fmt.Println(gray("Is applied:  "), yellowf("%v", record.Applied))
		fmt.Println(gray("Pending:     "), redf("%v", record.Pending))
		fmt.Println(gray("-----------------"))
	}
}
