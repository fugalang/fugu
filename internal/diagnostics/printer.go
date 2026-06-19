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
					for i := 0; i < int(err.Pos.Column)-1; i++ {
						if i < len(line) && line[i] == '\t' {
							sb.WriteString("\t")
						} else {
							sb.WriteString(" ")
						}
					}
					lenArr := err.End - err.Start
					for lenArr > 0 {
						sb.WriteString(color.BoldRed("^"))
						lenArr--
					}
					sb.WriteString(" ")
					sb.WriteString(color.BoldRed(err.Arrow))
					sb.WriteString("\n")
				}
			}
			sb.WriteString("\n")
		}

		// Вычисляем максимальную длину строки в описании для корректного отображения блока
		maxLineLength := 0
		offset := 4 // отступ для рамки блока

		for _, desc := range err.Description {
			if len(desc) > maxLineLength {
				maxLineLength = len(desc)
			}
		}

		maxLineLength += offset

		if maxLineLength > 256 {
			maxLineLength = 256
		}
		sb.WriteString("╭")
		for i := 0; i < maxLineLength; i++ {
			sb.WriteString("─")
		}
		sb.WriteString("╮\n")
		for _, desc := range err.Description {
			sb.WriteString("│ ")
			sb.WriteString(desc)
			offset := maxLineLength - len(desc) - 2
			for i := 0; i < offset; i++ {
				sb.WriteString(" ")
			}
			sb.WriteString(" │\n")
		}
		sb.WriteString("╰")
		for i := 0; i < maxLineLength; i++ {
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
