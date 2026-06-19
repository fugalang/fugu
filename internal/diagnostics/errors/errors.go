package errors

import "github.com/fugalang/fugu/internal/token"

// Не менять массив!!
var Errors = []Error{
	{
		Code:     0,
		CodeName: "TestError",
		Message:  "тестовая ошибка для проверки механизма диагностики",
		Arrow:    "BLOCK=FALSE",
		Description: []string{
			"тест ошибка не должна быть в релизе!!",
		},
	},
	{
		Code:     1,
		CodeName: "LexerIllegal",
		Message:  "недопустимый символ",
		Arrow:    "BLOCK=FALSE",
		Description: []string{
			"символ не распознан, возможно он не поддерживается или это опечатка",
		},
	},
	{
		Code:     2,
		CodeName: "LexerNoClosing",
		Message:  "не найден закрывающий символ",
		Arrow:    "Закрой за собой!",
		Description: []string{
			"открывающий символ не был закрыт до конца файла, возможно пропущен или это опечатка",
			"проверьте, что все открывающие символы (например, кавычки, скобки) имеют соответствующие закрывающие символы",
			"если ошибка возникает внутри строки, проверьте правильность экранирования символов внутри строки",
			"возможно вы забыли закрыть строку или комментарий, проверьте соответствующие символы в коде",
		},
	},
}

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

func (e Error) Update(tk token.Token) Error {
	e.Start = tk.Start
	e.End = tk.End
	e.Pos = tk.Pos
	return e
}
