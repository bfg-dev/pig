package output

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bfg-dev/pig/db"
	"github.com/bfg-dev/pig/file"
	"github.com/bfg-dev/pig/migration"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	red     = color.New(color.FgRed).SprintFunc()
	redf    = color.New(color.FgRed).SprintfFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	yellowf = color.New(color.FgYellow).SprintfFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	greenf  = color.New(color.FgGreen).SprintfFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	cyanf   = color.New(color.FgCyan).SprintfFunc()
)

func stringPointerToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

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

// PrintMigrations - print small table
func PrintMigrations(records *migration.Migrations, hideFilename bool) {
	table := tablewriter.NewWriter(os.Stdout)
	if hideFilename {
		table.SetHeader([]string{"ID", "Timestamp", "Name", "Requirements", "Is applied", "Pending"})

		table.SetHeaderColor(
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
		)

		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgHiBlackColor},
			tablewriter.Colors{},
			tablewriter.Colors{tablewriter.FgHiBlueColor},
			tablewriter.Colors{tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.FgYellowColor},
			tablewriter.Colors{tablewriter.FgRedColor},
		)
	} else {
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
	}

	for _, record := range records.Items {
		recordID := "?"
		if record.ID != 0 {
			recordID = fmt.Sprintf("%v", record.ID)
		}

		recordRequirements := []string{}
		for _, req := range record.Requirements {
			recordRequirements = append(recordRequirements, req.Name)
		}

		if hideFilename {
			table.Append([]string{
				recordID,
				record.TStamp.Format("2006-01-02 15:04:05"),
				record.Name,
				strings.Join(recordRequirements, ", "),
				fmt.Sprintf("%v", record.Applied),
				fmt.Sprintf("%v", record.Pending),
			})
		} else {
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
	}
	table.Render() // Send output
}

// PrintHistory - print history
func PrintHistory(records []*db.RecHistory) {
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

// Fatal - log fatal
func Fatal(v ...interface{}) {
	log.Fatal(red(v))
}

// Fatalf - log fatal
func Fatalf(format string, v ...interface{}) {
	log.Fatal(redf(format, v))
}

// Info1 - log info
func Info1(v ...interface{}) {
	log.Println(yellow(v))
}

// Infof1 - log info
func Infof1(format string, v ...interface{}) {
	log.Println(yellowf(format, v))
}

// Info2 - log info
func Info2(v ...interface{}) {
	log.Println(cyan(v))
}

// Infof2 - log info
func Infof2(format string, v ...interface{}) {
	log.Println(cyanf(format, v))
}

// OK - log info
func OK(v ...interface{}) {
	log.Println(green(v))
}

// OKf - log info
func OKf(format string, v ...interface{}) {
	log.Println(greenf(format, v))
}
