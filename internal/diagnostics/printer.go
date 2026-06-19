package diagnostics

import (
	"fmt"
	"strings"

	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/pkg/color"
)

type DiagnosticArena struct {
	Source string
	Errors []errors.Error
}

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
	for _, err := range a.Errors {
		var sb strings.Builder

		sb.WriteString(color.BoldRed("Error["))
		sb.WriteString(color.BoldYellow(err.CodeName))
		sb.WriteString(color.BoldRed("]: "))
		sb.WriteString(color.PastelYellow(err.Message))
		sb.WriteString("\n")

		sb.WriteString(color.NoteLabel(err.Pos.FileName))
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprint(err.Pos.Line))
		sb.WriteString(":")
		sb.WriteString(fmt.Sprint(err.Pos.Column))
		sb.WriteString("\n\n")

		if err.Arrow != "BLOCK=FALSE" {
			lines := GetLine(a.Source, int(err.Pos.Line), int(err.Pos.Line)-5)

			for i, line := range lines {
				lineNum := int(err.Pos.Line) - len(lines) + i + 1

				sb.WriteString(color.Black(fmt.Sprintf("%4d | ", lineNum)))
				sb.WriteString(line)
				sb.WriteString("\n")

				if lineNum == int(err.Pos.Line) {
					sb.WriteString(color.Black("     | "))

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

					// highlight range
					length := int(err.End - err.Start)
					if length <= 0 {
						length = 1
					}

					for i := 0; i < length; i++ {
						sb.WriteString(color.BoldRed("^"))
					}

					sb.WriteString(" ")
					sb.WriteString(color.BoldRed(err.Arrow))
					sb.WriteString("\n")
				}
			}

			sb.WriteString("\n")
		}

		maxLineLength := 0

		for _, desc := range err.Description {
			l := len([]rune(desc))
			if l > maxLineLength {
				maxLineLength = l
			}
		}

		if maxLineLength > 120 {
			maxLineLength = 120
		}

		sb.WriteString("╭")
		for i := 0; i < maxLineLength+2; i++ {
			sb.WriteString("─")
		}
		sb.WriteString("╮\n")

		for _, desc := range err.Description {
			lineLen := len([]rune(desc))
			pad := maxLineLength - lineLen

			sb.WriteString("│ ")
			sb.WriteString(desc)

			for i := 0; i < pad; i++ {
				sb.WriteString(" ")
			}

			sb.WriteString(" │\n")
		}

		sb.WriteString("╰")
		for i := 0; i < maxLineLength+2; i++ {
			sb.WriteString("─")
		}
		sb.WriteString("╯\n")

		fmt.Println(sb.String())
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
