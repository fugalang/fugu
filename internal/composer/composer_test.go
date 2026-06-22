package composer

import (
	"testing"

	"github.com/fugalang/fugu/internal/composer/project"
	"github.com/k0kubun/pp/v3"
)

func TestCacherLoadingLibrary(t *testing.T) {
	prj := project.InitProject("app")
	pp.Println(prj)
}
