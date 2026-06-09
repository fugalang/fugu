package reporter

type Code uint16

const (
	TestError Code = iota

	// NoClosing
	LexerNoClosing // не закрыт блок и тд
)

// TODO: надо будет очев сделать чтобы возвращал или на англ или на ру
func (c Code) String() string {
	switch c {
	case TestError:
		return "тестовая ошибка"
	case LexerNoClosing:
		return "пропущен закрывающий символ"
	default:
		return "неизвестная ошибка"
	}
}

func (c Code) Notes() []string {
	switch c {
	case TestError:
		return []string{
			"Ой, кажется, я сломался изнутри! Это внутренняя отладочная ошибка.",
			"Я не должен был спотыкаться на этом месте при обычной работе.",
			"Пожалуйста, создайте об этом баг-репорт:",
			"  https://github.com/fugalang/fugu/issues",
			"Если вам не трудно, опишите в репорте сценарий, при котором вы обнаружили эту проблему.",
		}

	case LexerNoClosing:
		return []string{""}
	default:
		return []string{}
	}
}

func (c Code) Code() string {
	switch c {
	case TestError:
		return "TEST"
	case LexerNoClosing:
		return "LNC1"
	default:
		return "I!K"
	}
}

func (c Code) Arrow() string {
	switch c {
	case TestError:
		return ""
	case LexerNoClosing:
		return "закрой за собой!"
	default:
		return ""
	}
}
