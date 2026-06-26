package lexer

import (
	"fmt"
	"testing"

	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/token"
	"github.com/k0kubun/pp/v3"
)

func TestComment(t *testing.T) {
	// 1. Описываем структуру тест-кейса
	tests := []struct {
		name            string
		input           []byte
		checkSpaceTk    bool // проверка токена пробел следущим токеном
		expectedKind    token.Kind
		expectedLiteral []byte
	}{
		{
			name:            "Тест обработки однострочный комментарий",
			input:           []byte("// привет slava"),
			expectedKind:    token.COMMENT,
			expectedLiteral: []byte(" привет slava"),
		},
		{
			name: "Тест обработки многострочного коментария",
			input: []byte(`/* 
Привет, это многострочный коммент для проверки корректной работы лексера) 
*/ 
`),
			checkSpaceTk: true,
			expectedKind: token.M_COMMENT,
			expectedLiteral: []byte(` 
Привет, это многострочный коммент для проверки корректной работы лексера) 
`),
		},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg", &diagnostics.Arena{
			Source: string(tt.input),
		})
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue // переход к след тесту
		}

		lit := tk.Literal(&lex.Input)
		if string(lit) != string(tt.expectedLiteral) {
			t.Errorf("[%s] Неверный литерал.\nОжидался:\n%q\n\nПолучен:\n%q",
				tt.name, tt.expectedLiteral, lit)
		}

		if tt.checkSpaceTk {
			tk = lex.NextToken()
			if tk.Kind != token.SPACING {
				t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
					tt.name, token.SPACING.String(), tk.Kind.String())
			}
		}

		tk = lex.NextToken()
		if tk.Kind != token.EOF {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, token.EOF.String(), tk.Kind.String())
		}

	}
}

func TestOperator(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedKind token.Kind
	}{
		// Операторы диапазонов
		{name: "Исключающий диапазон", input: []byte(".."), expectedKind: token.OP_RANGE},
		{name: "Включающий диапазон", input: []byte("..="), expectedKind: token.RANGE_INCL},
		{name: "Полуоткрытый диапазон", input: []byte("..<"), expectedKind: token.RANGE_HALF_OPEN},

		// Операторы присваивания
		{name: "Обычное переопределение", input: []byte("="), expectedKind: token.ASSIGN},
		{name: "Уменьшение с присваиванием", input: []byte("-="), expectedKind: token.SUB_ASSIGN},
		{name: "Увеличение с присваиванием", input: []byte("+="), expectedKind: token.ADD_ASSIGN},
		{name: "Умножение с присваиванием", input: []byte("*="), expectedKind: token.MUL_ASSIGN},
		{name: "Деление с присваиванием", input: []byte("/="), expectedKind: token.DIV_ASSIGN},
		{name: "Остаток с присваиванием", input: []byte("%="), expectedKind: token.MOD_ASSIGN},
		{name: "Возведение в степень с присваиванием", input: []byte("^="), expectedKind: token.POW_ASSIGN},

		// Логические операторы сравнения
		{name: "Равенство", input: []byte("=="), expectedKind: token.EQ},
		{name: "Неравенство", input: []byte("!="), expectedKind: token.NEQ},
		{name: "Меньше или равно", input: []byte("<="), expectedKind: token.LE},
		{name: "Больше или равно", input: []byte(">="), expectedKind: token.GE},
		{name: "Меньше", input: []byte("<"), expectedKind: token.LT},
		{name: "Больше", input: []byte(">"), expectedKind: token.GT},
		{name: "Логическое НЕ", input: []byte("!"), expectedKind: token.BANG},
		{name: "Логическое И", input: []byte("&&"), expectedKind: token.AND},
		{name: "Логическое ИЛИ", input: []byte("||"), expectedKind: token.OR},

		// Операторы арифметики
		{name: "Минус", input: []byte("-"), expectedKind: token.SUB},
		{name: "Плюс", input: []byte("+"), expectedKind: token.ADD},
		{name: "Умножение", input: []byte("*"), expectedKind: token.MUL},
		{name: "Деление", input: []byte("/"), expectedKind: token.DIV},
		{name: "Остаток от деления", input: []byte("%"), expectedKind: token.MOD},
		{name: "Степень", input: []byte("^"), expectedKind: token.POW},

		// Побитовые операторы
		{name: "Битовый сдвиг влево", input: []byte("<<"), expectedKind: token.SHR_LESS},
		{name: "Битовый сдвиг вправо", input: []byte(">>"), expectedKind: token.SHR_GREATER},
		{name: "Побитовое НЕ", input: []byte("~"), expectedKind: token.BITWISE_NOT},

		// Операторы управления данными
		{name: "Лямбда", input: []byte("=>"), expectedKind: token.ARROW},
		{name: "Пайплайн", input: []byte("|>"), expectedKind: token.PIPE},
		{name: "Тернарный оператор", input: []byte("?:"), expectedKind: token.DEFAULT},
		{name: "Безопасный вызов", input: []byte("?."), expectedKind: token.OPTIONAL_DOT},
		{name: "Взятие ссылки", input: []byte("&"), expectedKind: token.REF},

		// Операторы группировки
		{name: "Левая круглая скобка: ( ", input: []byte("("), expectedKind: token.L_PAREN},
		{name: "Правая круглая скобка: ) ", input: []byte(")"), expectedKind: token.R_PAREN},
		{name: "Левая фигурная скобка: { ", input: []byte("{"), expectedKind: token.L_BRACE},
		{name: "Правая фигурная скобка: } ", input: []byte("}"), expectedKind: token.R_BRACE},
		{name: "Левая квадратная скобка: [ ", input: []byte("["), expectedKind: token.L_BRACK},
		{name: "Правая квадратная скобка: ] ", input: []byte("]"), expectedKind: token.R_BRACK},

		// Операторы разделения
		{name: "Двоеточие", input: []byte(":"), expectedKind: token.COLON},
		{name: "Точка с запятой", input: []byte(";"), expectedKind: token.END},
		{name: "Запятая", input: []byte(","), expectedKind: token.COMMA},
		{name: "Точка", input: []byte("."), expectedKind: token.DOT},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg", &diagnostics.Arena{
			Source: string(tt.input),
		})
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue
		}
	}
}

func TestLiteral(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		expectedKind    token.Kind
		expectedLiteral []byte
	}{
		{name: "Идентификатор начинающийся числом", input: []byte("10pixel"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("10pixel")},
		{name: "Идентификатор с цифрами и подчёркиванием", input: []byte("2stack"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("2stack")},
		{name: "Обычный идентификатор с подчёркивания", input: []byte("__init__"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("__init__")},

		{name: "Ключевое слово module", input: []byte("mod"), expectedKind: token.MODULE, expectedLiteral: []byte("mod")},
		{name: "Ключевое слово fn", input: []byte("fn"), expectedKind: token.FN, expectedLiteral: []byte("fn")},

		{name: "Целое число (INTEGER)", input: []byte("12443"), expectedKind: token.INTEGER, expectedLiteral: []byte("12443")},
		{name: "Дробное число (FLOATING)", input: []byte("12.3"), expectedKind: token.FLOATING, expectedLiteral: []byte("12.3")},
		{name: "Мнимое целое число (IMAGINARY)", input: []byte("123i"), expectedKind: token.IMAGINARY, expectedLiteral: []byte("123i")},
		{name: "Мнимое дробное число (IMAGINARY)", input: []byte("12.3i"), expectedKind: token.IMAGINARY, expectedLiteral: []byte("12.3i")},

		{name: "Обычная строка (STRING)", input: []byte(`"hello world"`), expectedKind: token.STRING, expectedLiteral: []byte("hello world")},
		{name: "Сырая строка (RAW_STRING)", input: []byte("`multiline code`"), expectedKind: token.RAW_STRING, expectedLiteral: []byte("multiline code")},
		{name: "Одиночный символ (CHARACTER)", input: []byte("'я'"), expectedKind: token.CHARACTER, expectedLiteral: []byte("я")},
		{name: "Экранированный символ (CHARACTER)", input: []byte("'\\n'"), expectedKind: token.CHARACTER, expectedLiteral: []byte("\\n")},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg", &diagnostics.Arena{
			Source: string(tt.input),
		})
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue
		}

		if string(tk.Literal(&lex.Input)) != string(tt.expectedLiteral) {
			t.Errorf("[%s] Неверный литерал. Ожидался: %q, получен: %q",
				tt.name, tt.expectedLiteral, tk.Literal(&lex.Input))
			continue
		}
	}
}

func TestLexerStabilization(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedKind token.Kind
	}{
		{
			name:         "cтабилизация после незакрытого многострочного комментария",
			input:        []byte("/* незакрытый комментарий \n fn main() {}"),
			expectedKind: token.FN,
		},
		{
			name: "cтабилизация после незакрытой обычной строки",
			input: []byte(`"незакрытая строка
if x == 5 {}`),
			expectedKind: token.IF,
		},
		{
			name:         "cтабилизация после незакрытой сырой строки",
			input:        []byte("`незакрытая сырая строка \n else { return }"),
			expectedKind: token.ELSE,
		},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg", &diagnostics.Arena{
			Source: string(tt.input),
		})

		firstTok := lex.NextToken()
		if firstTok.Kind != token.ILLEGAL {
			t.Fatalf("[%s] Первый токен обязан быть ILLEGAL. Получен: %s", tt.name, firstTok.Kind.String())
		}

		secondTok := lex.NextToken()
		if secondTok.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена после стабилизации. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), secondTok.Kind.String())
		}

	}
}

func TestLexerString(t *testing.T) {
	input := []byte(`
"привет ${mod "${gen_name(pkg)}" } "
`)
	lex := New(input, "main.fg", diagnostics.New(string(input)))

	tks := []token.Token{}
	for {
		tk := lex.NextToken()
		switch tk.Kind {
		case token.COMMENT, token.M_COMMENT, token.SPACING:
			continue
		}
		if tk.Kind == token.EOF {
			break
		}
		tks = append(tks, tk)
	}

	pp.Println(tks)
	for _, t := range tks {
		fmt.Println(
			t.Kind.String()+":",
			string(t.Literal(&input)),
			t.Start,
		)
	}
}
