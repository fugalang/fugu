package errors

import "github.com/fugalang/fugu/internal/token"

// Не менять массив!!
var errors = []Error{}

type Error struct {
	Code        uint16
	CodeName    string   // название ошибки, для удобства
	Message     string   // сообщение об ошибке, кратокое описание
	Arrow       string   // строка с указанием места ошибки и пояснением. если равен "BLOCK=FALSE" то не будет блока с кодом, а только указание на строку и столбец
	Description []string // подробное описание ошибки

	Start uint64
	End   uint64
	Pos   token.Position
}
