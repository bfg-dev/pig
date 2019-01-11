package output

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bfg-dev/pig/migration"
)

// CheckGraphviz - check if graphviz exists
func CheckGraphviz() bool {
	cmd := exec.Command("dot", "-V")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GraphPng - draw png file
func GraphPng(migs *migration.Migrations, filename string) error {

	var (
		dotData string
	)

	dotData = "digraph D {"
	for _, item := range migs.Items {
		for _, req := range item.Requirements {
			dotData += fmt.Sprintf("\"%v\" -> \"%v\"\n", item.Name, req.Name)
		}
	}
	dotData += "}"

	pngFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("dot", "-Tpng")
	cmd.Stdin = strings.NewReader(dotData)
	cmd.Stdout = pngFile
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
