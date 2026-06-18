package diagnostics

type Code uint16

const (
	NoError Code = iota
	TestError
	LexerNoClosing
	ParserCantStartWork
	StateDoesNotToken
)

type Msg interface {
	Code() string
	Msg() string
	Notes() []string
	Arrow() string
	IsUseBlock() bool
}

type meta struct {
	code     string
	msg      string
	notes    []string
	arrow    string
	useBlock bool
}

var defaultNotes = []string{
	"Внутренняя ошибка отладки: сбой компонента.",
	"Некорректное состояние, не ожидаемое при штатной работе.",
	"Пожалуйста, сообщите об ошибке по адресу:",
	"  https://github.com/fugalang/fugu/issues",
}

var codes = []meta{
	NoError: {
		code: "NoError",
		msg:  "не найденна",
	},

	TestError: {
		code:  "TestError",
		msg:   "тестовая ошибка",
		notes: defaultNotes,
	},

	LexerNoClosing: {
		code:     "LexerNoClosing",
		msg:      "пропущен закрывающий символ",
		useBlock: true,
		arrow:    "закрой за собой!",
		notes:    defaultNotes,
	},

	ParserCantStartWork: {
		code:     "ParserCantStartWork",
		msg:      "нету возможности запучтить разбор",
		useBlock: true,
		notes: []string{
			"Исправьте прошлые ошибки, чтобы парсер отработал корректно.",
		},
	},

	StateDoesNotToken: {
		code: "StateDoesNotToken",
		msg:  "ошибка при работе с таблицой состояний",
	},
}

func (c Code) meta() meta {
	i := int(c)
	if i < 0 || i >= len(codes) {
		return meta{
			code: "Unknown",
			msg:  "неизвестная ошибка",
		}
	}
	return codes[i]
}

func (c Code) Code() string {
	return c.meta().code
}

func (c Code) Msg() string {
	return c.meta().msg
}

func (c Code) Notes() []string {
	return c.meta().notes
}

func (c Code) Arrow() string {
	return c.meta().arrow
}

func (c Code) IsUseBlock() bool {
	return c.meta().useBlock
}
