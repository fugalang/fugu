package diagnostics

import "github.com/fugalang/fugu/internal/diagnostics/errors"

type DiagnosticArena struct {
	errors []errors.Error
}
