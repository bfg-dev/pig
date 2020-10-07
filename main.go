package main

//
// REFACTOR ME!!!
//

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/bfg-dev/pig/db"
	"github.com/bfg-dev/pig/file"
	"github.com/bfg-dev/pig/migration"
	"github.com/bfg-dev/pig/output"

	pgx "github.com/jackc/pgx/v4"
)

var (
	flags    = flag.NewFlagSet("pig", flag.ExitOnError)
	dir      = flags.String("dir", ".", "directory with migration files")
	note     = flags.String("note", "", "custom note for migrations")
	gitinfo  = flags.String("gitinfo", "", "custom git information (branch or tag) for migrations")
	onlyPlan = flags.Bool("only-plan", false, "show migration plan")
	listView = flags.Bool("list-view", false, "Use 'list' insteed of 'table' view")
)

var (
	usagePrefix = `Usage: pig [OPTIONS] DBSTRING COMMAND

Example:
	pig "user=postgres dbname=postgres sslmode=disable" status
Options:
`

	usageCommands = `
Commands:
    init                            Init database
    up                              Up all available migrations
    up-migration NAME               Up a specific NAME
    up-gitinfo GITINFO              Up a specific GITINFO (affects only known migrations)
    up-note NOTE                    Up a specific NOTE (affects only known migrations)
    down-migration NAME             Roll back a specific NAME
    down-gitinfo GITINFO            Roll back a specific GITINFO
    down-note NOTE                  Roll back a specific NOTE
    reset                           Roll back all migrations
    status                          Dump the migration status
    history                         Show migration history
    history-migration NAME          Show history for a specific NAME
    history-gitinfo GITINFO         Show history for a specific GITINFO
    history-note NOTE               Show history for a specific NOTE
    graph [pngname]                 Draw png graph. (default pngname is output.png)
		graph-migration NAME [pngname]  Draw png graph for specific NAME (default pngname is NAME.png)
		generate NAME                   Generate sql template
`
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	if args[0] == "generate" {
		if len(args[1]) == 0 {
			output.Fatal("Please provide NAME")
		}
		output.Info1("Generating file", args[1])
		err := file.GenerateNewSQLFile(dir, args[1])
		if err != nil {
			output.Fatal(err)
		}
		output.OK("File generated", args[1])
		return
	}

	dbstring, command := args[0], args[1]

	db, err := pgx.Connect(context.Background(), dbstring)
	if err != nil {
		output.Fatalf("-dbstring=%q: %v", dbstring, err)
	}

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[2:]...)
	}

	manager := migration.NewManager(db, "pig_version_table", "pig_history_table", *dir)

	output.Info1("Waiting for global lock")
	if err := manager.GetDBGlobalLock(); err != nil {
		output.Fatalf("global lock error: %v", err)
	}

	run(command, manager, *note, *gitinfo, *onlyPlan, arguments)
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

func upAndDown(f func(*migration.Meta) error, migrations *migration.Migrations, note, gitinfo string) {
	if len(migrations.Items) == 0 {
		log.Println("There are no migrations to execute")
		return
	}

	for _, m := range migrations.Items {
		m.Note = note
		m.GITinfo = gitinfo
		output.Info2("Executing: ", m.Name)
		if err := f(m); err != nil {
			output.Fatal(err)
		}
	}
	output.OK("All migrations executed")
}

func printMigrations(migs *migration.Migrations) {
	if *listView {
		output.PrintMigrationsStructs(migs)
	} else {
		output.PrintMigrationsTable(migs)
	}
}

func printHistory(history []*db.RecHistory) {
	if *listView {
		output.PrintHistoryStructs(history)
	} else {
		output.PrintHistoryTable(history)
	}
}

func run(command string, manager *migration.Manager, note, gitinfo string, onlyPlan bool, args []string) {
	switch command {
	case "up":
		output.Info1("UP plan")
		upMigrations, err := manager.GetUpPlan()
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(upMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteUp, upMigrations, note, gitinfo)
		}
	case "up-migration":
		if len(args) == 0 {
			output.Fatal("Please provide NAME")
		}
		output.Info1("UP plan for ", args[0])
		upMigrations, err := manager.GetUpPlanForName(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(upMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteUp, upMigrations, note, gitinfo)
		}
	case "up-gitinfo":
		if len(args) == 0 {
			output.Fatal("Please provide GITINFO")
		}
		output.Info1("UP plan for ", args[0])
		upMigrations, err := manager.GetUpPlanForGITinfo(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(upMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteUp, upMigrations, note, gitinfo)
		}
	case "up-note":
		if len(args) == 0 {
			output.Fatal("Please provide NOTE")
		}
		output.Info1("UP plan for ", args[0])
		upMigrations, err := manager.GetUpPlanForNote(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(upMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteUp, upMigrations, note, gitinfo)
		}
	case "down-migration":
		if len(args) == 0 {
			output.Fatal("Please provide NAME")
		}
		output.Info1("DOWN plan for", args[0])
		downMigrations, err := manager.GetDownPlanForName(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(downMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteDown, downMigrations, note, gitinfo)
		}
	case "down-gitinfo":
		if len(args) == 0 {
			output.Fatal("Please provide GITINFO")
		}
		output.Info1("DOWN plan for", args[0])
		downMigrations, err := manager.GetDownPlanForGITinfo(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(downMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteDown, downMigrations, note, gitinfo)
		}
	case "down-note":
		if len(args) == 0 {
			output.Fatal("Please provide NOTE")
		}
		output.Info1("DOWN plan for", args[0])
		downMigrations, err := manager.GetDownPlanForNote(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(downMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteDown, downMigrations, note, gitinfo)
		}
	case "reset":
		output.Info1("DOWN plan")
		downMigrations, err := manager.GetDownPlan()
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(downMigrations)
		if !onlyPlan {
			upAndDown(manager.ExecuteDown, downMigrations, note, gitinfo)
		}
	case "status":
		output.Info1("Status")
		allMigrations, err := manager.GetAllMigrations()
		if err != nil {
			output.Fatal(err)
		}
		printMigrations(allMigrations)
	case "init":
		output.Info2("Initiating database")
		if err := manager.InitDB(); err != nil {
			output.Fatal(err)
		}
		output.OK("DB initiated")
	case "history":
		output.Info1("History")
		history, err := manager.GetHistory()
		if err != nil {
			output.Fatal(err)
		}
		printHistory(history)
	case "history-migration":
		if len(args) == 0 {
			output.Fatal("Please provide NAME")
		}
		output.Info1("History for ", args[0])
		history, err := manager.GetHistoryForName(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printHistory(history)
	case "history-gitinfo":
		if len(args) == 0 {
			output.Fatal("Please provide GITINFO")
		}
		output.Info1("History for ", args[0])
		history, err := manager.GetHistoryForGITinfo(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printHistory(history)
	case "history-note":
		if len(args) == 0 {
			output.Fatal("Please provide NOTE")
		}
		output.Info1("History for ", args[0])
		history, err := manager.GetHistoryForNote(args[0])
		if err != nil {
			output.Fatal(err)
		}
		printHistory(history)
	case "graph":
		var filename = "output.png"

		output.Info1("Graph")
		allMigrations, err := manager.GetAppliedMigrations()
		if err != nil {
			output.Fatal(err)
		}
		if !output.CheckGraphviz() {
			output.Fatal("Graphviz not found. Please install it. [https://www.graphviz.org/download/]")
		}
		if len(args) != 0 {
			filename = args[0]
		}
		if err := output.GraphPng(allMigrations, filename); err != nil {
			output.Fatal(err)
		}
		output.OK("Saved to", filename)
	case "graph-migration":
		if len(args) == 0 {
			output.Fatal("Please provide NAME")
		}
		output.Info1("Graph plan for", args[0])

		var filename = args[0] + ".png"

		upMigrations, err := manager.GetUpPlanForName(args[0])
		if err != nil {
			output.Fatal(err)
		}
		if !output.CheckGraphviz() {
			output.Fatal("Graphviz not found. Please install it. [https://www.graphviz.org/download/]")
		}
		if len(args) >= 2 {
			filename = args[1]
		}
		if err := output.GraphPng(upMigrations, filename); err != nil {
			output.Fatal(err)
		}
		output.OK("Saved to", filename)
	default:
		output.Fatalf("%q: no such command", command)
	}
}
