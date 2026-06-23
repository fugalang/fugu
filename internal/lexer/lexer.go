package lexer

import (
	"unicode"
	"unicode/utf8"

	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	. "github.com/fugalang/fugu/internal/token"
)

type Lexer struct {
	Input []byte
	rn    rune // текущая rune

	curPos         uint64 // абсолютное смещение c начала файла
	tokStart       uint64 // абсолютное смещение до начала токена который разбираеться прямо сейчас
	tokStartLine   uint64 // номер строки начала токена
	tokStartColumn uint64 // номер колонки начала токена
	pos            Position

	savePoint saveLexer
	da        *diagnostics.Arena
}

// для заморозки состояния
type saveLexer struct {
	rn             rune
	curPos         uint64
	tokStart       uint64
	tokStartLine   uint64
	tokStartColumn uint64
	pos            Position
}

func New(input []byte, fileName string, da *diagnostics.Arena) *Lexer {
	if input == nil {
		input = make([]byte, 0)
	}
	lex := &Lexer{
		Input:  input,
		curPos: 0,
		pos: Position{
			FileName: fileName,
			Line:     1,
			Column:   0,
			Offset:   0,
		},
	}

	lex.da = da

	lex.advance()
	return lex
}

func (lex *Lexer) Reset() {
	lex = New(lex.Input, lex.pos.FileName, lex.da)
}

func (lex *Lexer) NextToken() Token {
	lex.tokStart = lex.pos.Offset
	lex.tokStartLine = lex.pos.Line
	lex.tokStartColumn = lex.pos.Column

	if lex.rn == 0 {
		return lex.NewToken(EOF)
	}

	if unicode.IsSpace(lex.rn) {
		for unicode.IsSpace(lex.rn) {
			lex.advance()
		}
		return lex.NewToken(SPACING)
	}

	switch lex.rn {
	case '/':
		if lex.peekRn() == '/' {
			return lex.readLineComment()
		} else if lex.peekRn() == '*' {
			return lex.readMultiLineComment()
		} else if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(DIV_ASSIGN)
		}

		lex.advance()
		return lex.NewToken(DIV)

	case '.':
		if lex.peekRn() == '.' {
			lex.advance() // едим первую .
			if lex.peekRn() == '=' {
				lex.advance().advance()
				return lex.NewToken(RANGE_INCL)
			} else if lex.peekRn() == '<' {
				lex.advance().advance()
				return lex.NewToken(RANGE_HALF_OPEN)
			} else if lex.peekRn() == '.' {
				lex.advance().advance()
				return lex.NewToken(OP_ARRAY)
			}
			lex.advance()
			return lex.NewToken(OP_RANGE)
		}
		lex.advance()
		return lex.NewToken(DOT)

	case '<':
		if lex.peekRn() == '<' {
			lex.advance().advance()
			return lex.NewToken(SHR_LESS)
		} else if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(LE)
		} else if lex.peekRn() == '-' {
			return lex.NewToken(CHAN_SEND)
		}
		lex.advance()
		return lex.NewToken(LT)

	case '>':
		if lex.peekRn() == '>' {
			lex.advance().advance()
			return lex.NewToken(SHR_GREATER)
		} else if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(GE)
		}
		lex.advance()
		return lex.NewToken(GT)

	case '-':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(SUB_ASSIGN)
		} else if lex.peekRn() == '>' {
			lex.advance().advance()
			return lex.NewToken(RTN_ARROW)
		}

		lex.advance()
		return lex.NewToken(SUB)

	case '+':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(ADD_ASSIGN)
		}
		lex.advance()
		return lex.NewToken(ADD)

	case '*':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(MUL_ASSIGN)
		}
		lex.advance()
		return lex.NewToken(MUL)

	case '%':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(MOD_ASSIGN)
		}
		lex.advance()
		return lex.NewToken(MOD)

	case '^':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(POW_ASSIGN)
		}
		lex.advance()
		return lex.NewToken(POW)

	case '~':
		lex.advance()
		return lex.NewToken(BITWISE_NOT)

	case '&':
		if lex.peekRn() == '&' {
			lex.advance().advance()
			return lex.NewToken(AND)
		}
		lex.advance()
		return lex.NewToken(REF)

	case '!':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(NEQ)
		}
		lex.advance()
		return lex.NewToken(BANG)

	case '?':
		if lex.peekRn() == ':' {
			lex.advance().advance()
			return lex.NewToken(DEFAULT)
		} else if lex.peekRn() == '.' {
			lex.advance().advance()
			return lex.NewToken(OPTIONAL_DOT)
		}

	case '=':
		if lex.peekRn() == '=' {
			lex.advance().advance()
			return lex.NewToken(EQ)
		} else if lex.peekRn() == '>' {
			lex.advance().advance()
			return lex.NewToken(ARROW)
		}

		lex.advance()
		return lex.NewToken(ASSIGN)

	case '|':
		if lex.peekRn() == '|' {
			lex.advance().advance()
			return lex.NewToken(OR)
		} else if lex.peekRn() == '>' {
			lex.advance().advance()
			return lex.NewToken(PIPE)
		}

	case ':':
		lex.advance()
		return lex.NewToken(COLON)

	case '(':
		lex.advance()
		return lex.NewToken(L_PAREN)

	case ')':
		lex.advance()
		return lex.NewToken(R_PAREN)

	case '{':
		lex.advance()
		return lex.NewToken(L_BRACE)

	case '}':
		lex.advance()
		return lex.NewToken(R_BRACE)

	case '[':
		lex.advance()
		return lex.NewToken(L_BRACK)

	case ']':
		lex.advance()
		return lex.NewToken(R_BRACK)

	case ';':
		lex.advance()
		return lex.NewToken(END)

	case ',':
		lex.advance()
		return lex.NewToken(COMMA)

	case '"':
		return lex.readString()
	case '`':
		return lex.readRawString()

	case '\'':
		return lex.readChar()

	default:
		if unicode.IsDigit(lex.rn) {
			return lex.readNumber()
		} else if unicode.IsLetter(lex.rn) || lex.rn == '_' {
			for unicode.IsLetter(lex.rn) || unicode.IsDigit(lex.rn) || lex.rn == '_' {
				lex.advance()
			}

			return lex.NewToken(SearchKeyword(lex.Input[lex.tokStart:lex.pos.Offset]))
		}

		lex.advance()
		return lex.NewToken(ILLEGAL)

	}

	return lex.NewToken(ILLEGAL)
}

func (lex *Lexer) readLineComment() Token {
	lex.advance().advance() // '//'

	// останавливаемся перед '\n'
	for lex.rn != '\n' && lex.rn != 0 {
		lex.advance()
	}

	return lex.NewToken(COMMENT)
}

func (lex *Lexer) readMultiLineComment() Token {
	lex.advance().advance() // '/*'

	lex.freezing()

	for {
		if lex.rn == 0 {
			tk := lex.NewToken(ILLEGAL)
			lex.da.Add(errors.Errors[2].Update(tk))
			lex.unfreeze()
			lex.stabilization()
			return tk
		}

		if lex.rn == '*' && lex.peekRn() == '/' {
			lex.advance().advance() // '*', '/'
			break
		}

		lex.advance()
	}

	return lex.NewToken(M_COMMENT)
}

func (lex *Lexer) readString() Token {
	lex.advance() // '"'
	isTemplate := false

	lex.freezing()

	for lex.rn != '"' && lex.rn != 0 {
		if lex.rn == '\\' {
			lex.advance().advance()
			continue
		}

		if lex.rn == '$' && lex.peekRn() == '{' {
			isTemplate = true
		}

		lex.advance()
	}

	if lex.rn == 0 {
		tk := lex.NewToken(ILLEGAL)
		lex.da.Add(errors.Errors[2].Update(tk))
		lex.unfreeze()
		lex.stabilization()
		return tk
	}

	lex.advance() // '"'

	if isTemplate {
		return lex.NewToken(T_STRING)
	}
	return lex.NewToken(STRING)
}

func (lex *Lexer) readRawString() Token {
	lex.advance() // '`'

	lex.freezing()

	for lex.rn != '`' && lex.rn != 0 {
		lex.advance()
	}

	if lex.rn == 0 {
		tk := lex.NewToken(ILLEGAL)
		lex.da.Add(errors.Errors[2].Update(tk))
		lex.unfreeze()
		lex.stabilization()
		return tk
	}

	lex.advance()
	return lex.NewToken(RAW_STRING)
}

func (lex *Lexer) readChar() Token {
	lex.advance()

	if lex.rn == '\\' {
		lex.advance().advance()
	} else if lex.rn != '\'' && lex.rn != 0 {
		lex.advance()
	}

	if lex.rn != '\'' {
		return lex.NewToken(ILLEGAL)
	}

	lex.advance()
	return lex.NewToken(CHARACTER)
}

func (lex *Lexer) readNumber() Token {
	isFloat := false
	isIdent := false

	for unicode.IsDigit(lex.rn) || unicode.IsLetter(lex.rn) || lex.rn == '_' || lex.rn == '.' {

		if lex.rn == '.' {
			if isIdent || isFloat {
				break
			}
			if !unicode.IsDigit(lex.peekRn()) {
				break
			}
			isFloat = true
		}

		if unicode.IsLetter(lex.rn) || lex.rn == '_' {
			isIdent = true
		}

		lex.advance()
	}

	literal := lex.Input[lex.tokStart:lex.pos.Offset]

	if isIdent {
		if literal[len(literal)-1] == 'i' && (!isFloat || len(literal) > 2) {
			onlyDigits := true
			for i := 0; i < len(literal)-1; i++ {
				if (literal[i] < '0' || literal[i] > '9') && literal[i] != '.' {
					onlyDigits = false
					break
				}
			}
			if onlyDigits {
				return lex.NewToken(IMAGINARY)
			}
		}

		return lex.NewToken(IDENTIFIER)
	}

	if isFloat {
		return lex.NewToken(FLOATING)
	}
	return lex.NewToken(INTEGER)
}

func (lex *Lexer) stabilization() {
	tkws := map[Kind]bool{
		FN:     true,
		IF:     true,
		SWITCH: true,
		CASE:   true,
		RETURN: true,
		ENUM:   true,
		SELECT: true,
	}

	for {
		lex.freezing()
		tk := lex.NextToken()

		if tk.Kind == EOF {
			return
		} else if tk.Kind == SPACING || tk.Kind == COMMENT || tk.Kind == M_COMMENT {
			continue
		} else if tkws[tk.Kind] {
			lex.unfreeze()
			return
		} else if tk.Kind == R_BRACE || tk.Kind == END {
			lex.unfreeze()
			return
		}
	}
}

func (lex *Lexer) advance() *Lexer {
	if lex.curPos >= uint64(len(lex.Input)) {
		lex.rn = 0 // \x00
		lex.pos.Offset = lex.curPos
		return lex
	}

	r, size := utf8.DecodeRune(lex.Input[lex.curPos:])

	lex.rn = r
	lex.pos.Offset = lex.curPos
	lex.curPos += uint64(size)

	if lex.rn == '\n' {
		lex.pos.Line++
		lex.pos.Column = 1
	} else {
		lex.pos.Column++
	}

	return lex
}

// возвращает следущий симвл после Lexer.curPos
func (lex *Lexer) peekRn() rune {
	if lex.curPos >= uint64(len(lex.Input)) {
		return 0
	}

	r, _ := utf8.DecodeRune(lex.Input[lex.curPos:])

	return r
}

func (lex *Lexer) NewToken(kind Kind) Token {
	return Token{
		Kind: kind,
		Pos: Position{
			FileName: lex.pos.FileName,
			Line:     lex.tokStartLine,
			Column:   lex.tokStartColumn,
			Offset:   lex.tokStart,
		},
		Start: lex.tokStart,
		End:   lex.pos.Offset,
	}
}

// заморозить состояния
func (lex *Lexer) freezing() {
	lex.savePoint = saveLexer{
		rn:             lex.rn,
		curPos:         lex.curPos,
		tokStart:       lex.tokStart,
		tokStartLine:   lex.tokStartLine,
		tokStartColumn: lex.tokStartColumn,
		pos:            lex.pos,
	}
}

// разморозить :)
func (lex *Lexer) unfreeze() {
	lex.rn = lex.savePoint.rn
	lex.curPos = lex.savePoint.curPos
	lex.tokStart = lex.savePoint.tokStart
	lex.tokStartLine = lex.savePoint.tokStartLine
	lex.tokStartColumn = lex.savePoint.tokStartColumn
	lex.pos = lex.savePoint.pos
}
