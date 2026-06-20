package diagnostics

import (
	"testing"

	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/token"
)

func TestPrintDiagnostics(t *testing.T) {
	arena := &DiagnosticArena{
		Source: `fn main() {
print("Hello, World!")
}`}

	arena.AddError(errors.Errors[3].Update(token.Token{
		Kind: token.FN,
		Pos: token.Position{
			FileName: "main.fg",
			Line:     1,
			Column:   1,
			Offset:   0,
		},
		Start: 0,
		End:   2,
	}))

	if !arena.HasErrors() {
		t.Errorf("Expected arena to have errors, but it does not.")
	}

	arena.Print()

	arena.Clear()

	if arena.HasErrors() {
		t.Errorf("Expected arena to be cleared of errors, but it still has errors.")
	}
}
