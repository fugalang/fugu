// Package diagnostics
//
// Здесь допускается использование более тяжёлых по ресурсам решений,
// поскольку этот пакет вызывается только при возникновении ошибок.
package diagnostics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/pkg/color"
)

type DiagnosticArena struct {
	Source string
	Errors []errors.Error
}

var Da = DiagnosticArena{}

func (a *DiagnosticArena) AddError(err errors.Error) {
	a.Errors = append(a.Errors, err)
}

func (a *DiagnosticArena) HasErrors() bool {
	return len(a.Errors) > 0
}

func (a *DiagnosticArena) Clear() {
	a.Errors = nil
}

func (a *DiagnosticArena) Print() {
	for i, err := range a.Errors {
		var sb strings.Builder

		sb.WriteString(color.BoldRed("× "))
		sb.WriteString(color.BoldRed(err.CodeName))
		sb.WriteString(color.Gray(": "))
		sb.WriteString(err.Message)
		sb.WriteString("\n")

		if err.Pos.Line != 0 {
			sb.WriteString(color.Blue("──> "))
		} else {
			sb.WriteString(color.Blue("Module: "))
		}
		sb.WriteString(color.BoldBlue(err.Pos.FileName))
		sb.WriteString(" ")
		if err.Pos.Line != 0 {
			sb.WriteString(color.BoldBlue(fmt.Sprint(err.Pos.Line)))
			sb.WriteString(color.BoldBlue(":"))
			sb.WriteString(color.BoldBlue(fmt.Sprint(err.Pos.Column)))
		}
		if err.Pos.Line != 0 {
			sb.WriteString("\n\n")
		} else {
			sb.WriteString("\n")
		}

		if err.Arrow != "BLOCK=FALSE" {
			lines := GetLine(a.Source, int(err.Pos.Line), int(err.Pos.Line)-5)

			maxLine := len(strings.Split(a.Source, "\n"))
			width := len(strconv.Itoa(maxLine))

			for i, line := range lines {
				lineNum := int(err.Pos.Line) - len(lines) + i + 1

				sb.WriteString(color.Black(
					fmt.Sprintf("%*d │ ", width, lineNum),
				))
				sb.WriteString(line)
				sb.WriteString("\n")

				if lineNum == int(err.Pos.Line) {
					sb.WriteString(color.Black(
						fmt.Sprintf("%*s │ ", width, ""),
					))

					col := int(err.Pos.Column) - 1
					if col < 0 {
						col = 0
					}

					runes := []rune(line)

					if col > len(runes) {
						col = len(runes)
					}

					for i := 0; i < col; i++ {
						sb.WriteString(" ")
					}

					length := int(err.End - err.Start)
					if length <= 0 {
						length = 1
					}

					for i := 0; i < length; i++ {
						sb.WriteString(color.BoldRed("^"))
					}

					sb.WriteString(" ")
					sb.WriteString(color.PastelYellow(err.Arrow))
					sb.WriteString("\n")
				}
			}

			sb.WriteString("\n")
		}

		for _, desc := range err.Description {
			sb.WriteString(color.BoldGreen("• "))
			sb.WriteString(desc)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")

		fmt.Println(sb.String())
		a.Errors = append(a.Errors[:i], a.Errors[i+1:]...)
	}
}

func GetLine(source string, lineNumber int, start int) []string {
	lines := strings.Split(source, "\n")

	idx := lineNumber - 1
	if idx < 0 || idx >= len(lines) {
		return nil
	}

	if len(lines) == 1 {
		return []string{lines[0]}
	}

	from := start
	if from < 0 {
		from = 0
	}

	if from >= len(lines) || from > idx {
		return []string{lines[idx]}
	}

	return lines[from : idx+1]
}
