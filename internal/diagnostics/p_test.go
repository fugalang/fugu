package diagnostics

import (
	"testing"

	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/token"
)

func TestPrintDiagnostics(t *testing.T) {
	arena := &DiagnosticArena{
		Source: `func main() {
print("Hello, World!")
}`}

	err := errors.Error{
		Code:     1,
		CodeName: "SyntaxError",
		Message:  "Unexpected token",
		Arrow:    "error",
		Description: []string{"This error occurs when the parser encounters an unexpected token in the source code.",
			"Check the syntax of your code and ensure that all tokens are correctly placed.",
		},
		Start: 0,
		End:   4,
		Pos: token.Position{
			FileName: "main.fg", Line: 1, Column: 0,
		},
	}

	arena.AddError(err)

	if !arena.HasErrors() {
		t.Errorf("Expected arena to have errors, but it does not.")
	}

	arena.Print()

	arena.Clear()

	if arena.HasErrors() {
		t.Errorf("Expected arena to be cleared of errors, but it still has errors.")
	}
}
