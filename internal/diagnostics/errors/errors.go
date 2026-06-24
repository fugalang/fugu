package errors

import "github.com/fugalang/fugu/internal/token"

type Severity uint8

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "error"
	}
}

// Не менять массив!!
var Errors = []Error{
	{
		Code:     0,
		CodeName: "TestError",
		Severity: SeverityWarning,
		Message:  "тестовая ошибка для проверки механизма диагностики",
		Description: []string{
			"тест-ошибка не должна попадать в релиз!!",
		},
	},
	{
		Code:     1,
		CodeName: "LexerIllegal",
		Severity: SeverityError,
		Message:  "недопустимый символ",
		Description: []string{
			"символ не распознан, возможно он не поддерживается или это опечатка",
		},
	},
	{
		Code:          2,
		CodeName:      "LexerNoClosing",
		Severity:      SeverityError,
		Message:       "не найден закрывающий символ",
		IsShowSnippet: true,
		Arrow:         "Закрой за собой!",
		Description: []string{
			"открывающий символ не был закрыт до конца файла, возможно пропущен или это опечатка",
			"проверьте, что все открывающие символы (например, кавычки, скобки) имеют соответствующие закрывающие символы",
			"если ошибка возникает внутри строки, проверьте правильность экранирования символов внутри строки",
			"возможно вы забыли закрыть строку или комментарий, проверьте соответствующие символы в коде",
		},
	},
	{
		Code:     3,
		CodeName: "ErrorLoadLibrary",
		Severity: SeverityError,
		Message:  "не удалось загрузить библиотеку",
		Description: []string{
			"не удалось загрузить библиотеку. Причина ошибки:",
		},
	},
	{
		Code:     4,
		CodeName: "ExecutingCommands",
		Severity: SeverityError,
		Message:  "не удалось выполнить команду",
		Description: []string{
			"ошибка выполнения команды: ",
		},
	},
	{
		Code:     5,
		CodeName: "ParsingError",
		Severity: SeverityError,
		Message:  "ошибка разбора",
		Description: []string{
			"ожидалось: было получено:",
		},
	},
	{
		Code:     6,
		CodeName: "ErrorActionMap",
		Severity: SeverityWarning,
		Message:  "ошибка работы с таблицей",
		Description: []string{
			"создайте https://github.com/fugalang/fugu/issues",
		},
	},
}

type Error struct {
	Code          uint16
	Severity      Severity
	CodeName      string // название ошибки, для удобства
	Message       string // сообщение об ошибке, кратокое описание
	Arrow         string // строка с указанием места ошибки и пояснением.
	IsShowSnippet bool
	Description   []string

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

func (e Error) IU(fileModule string, description []string) Error {
	e.Description = description
	e.Pos.FileName = fileModule
	e.Pos.Line = 0
	return e
}

func (e *Error) Error() string {
	return e.Message
}
