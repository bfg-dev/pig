package output

import "github.com/fatih/color"

var (
	red     = color.New(color.FgRed).SprintFunc()
	redf    = color.New(color.FgRed).SprintfFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	yellowf = color.New(color.FgYellow).SprintfFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	greenf  = color.New(color.FgGreen).SprintfFunc()
	cyan    = color.New(color.FgCyan).SprintFunc()
	cyanf   = color.New(color.FgCyan).SprintfFunc()
	gray    = color.New(color.FgHiBlack).SprintFunc()
	hiblue  = color.New(color.FgHiBlue).SprintFunc()
)

func init() {
  color.NoColor = false
}

func stringPointerToString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
