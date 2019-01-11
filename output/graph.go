package output

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bfg-dev/pig/migration"
)

var (
	colors = []string{
		"#8B0000",
		"#C71585",
		"#FF4500",
		"#FF8C00",
		"#FFD700",
		"#BDB76B",
		"#800080",
		"#4B0082",
		"#006400",
		"#808000",
		"#008B8B",
		"#0000FF",
	}
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
		dotData    string
		colorIndex int
	)

	dotData = "digraph D {"
	for _, item := range migs.Items {
		for _, req := range item.Requirements {
			dotData += fmt.Sprintf("\"%v\" -> \"%v\" [color=\"%v\"];\n", item.Name, req.Name, colors[colorIndex])
			colorIndex++
			if colorIndex >= len(colors) {
				colorIndex = 0
			}
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
