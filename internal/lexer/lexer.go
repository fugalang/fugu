package lexer

import (
	"unicode/utf8"

	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	. "github.com/fugalang/fugu/internal/token"
)

type Lexer struct {
	Input    []byte
	fileName string
	da       *diagnostics.Arena

	pos    int
	curPos int
	rn     rune // текущая rune
	rnSize int

	line int
	col  int

	savePos  int
	saveLine int
	saveCol  int
	saveRn   rune
	saveSize int

	tokPos  int
	tokLine int
	tokCol  int
}

func New(input []byte, fileName string, da *diagnostics.Arena) *Lexer {
	if input == nil {
		input = input[:0]
	}
	lex := &Lexer{
		Input:    input,
		fileName: fileName,
		da:       da,
		line:     1,
		col:      0,
	}
	lex.advance()
	return lex
}

func (lex *Lexer) advance() {
	if lex.curPos >= len(lex.Input) {
		lex.pos = lex.curPos
		lex.rn = 0
		lex.rnSize = 0
		return
	}

	lex.pos = lex.curPos

	if b := lex.Input[lex.curPos]; b < utf8.RuneSelf {
		lex.rn = rune(b)
		lex.rnSize = 1
	} else {
		r, size := utf8.DecodeRune(lex.Input[lex.curPos:])
		lex.rn = r
		lex.rnSize = size
	}

	lex.curPos += lex.rnSize

	if lex.rn == '\n' {
		lex.line++
		lex.col = 0
	} else {
		lex.col++
	}
}

func (lex *Lexer) peek() rune {
	if lex.curPos >= len(lex.Input) {
		return 0
	}
	if b := lex.Input[lex.curPos]; b < utf8.RuneSelf {
		return rune(b)
	}
	r, _ := utf8.DecodeRune(lex.Input[lex.curPos:])
	return r
}

// снимки состояния

func (lex *Lexer) freeze() {
	lex.savePos = lex.pos
	lex.saveLine = lex.line
	lex.saveCol = lex.col
	lex.saveRn = lex.rn
	lex.saveSize = lex.rnSize
}

func (lex *Lexer) unfreeze() {
	lex.pos = lex.savePos
	lex.curPos = lex.savePos + lex.saveSize
	lex.line = lex.saveLine
	lex.col = lex.saveCol
	lex.rn = lex.saveRn
	lex.rnSize = lex.saveSize
}

func (lex *Lexer) tok(kind Kind) Token {
	return Token{
		Kind: kind,
		Pos: Position{
			FileName: lex.fileName,
			Line:     uint64(lex.tokLine),
			Column:   uint64(lex.tokCol),
			Offset:   uint64(lex.tokPos),
		},
		Start: uint64(lex.tokPos),
		End:   uint64(lex.pos),
	}
}

func (lex *Lexer) NextToken() Token {
	lex.tokPos = lex.pos
	if isSpace(lex.rn) {
		for isSpace(lex.rn) {
			lex.advance()
		}
		return lex.tok(SPACING)
	}

	lex.tokPos = lex.pos
	lex.tokLine = lex.line
	lex.tokCol = lex.col

	if lex.rn == 0 {
		return lex.tok(EOF)
	}

	ch := lex.rn
	lex.advance()

	switch ch {
	case '/':
		switch lex.rn {
		case '/':
			return lex.lineComment()
		case '*':
			return lex.multiLineComment()
		case '=':
			lex.advance()
			return lex.tok(DIV_ASSIGN)
		}
		return lex.tok(DIV)

	case '.':
		if lex.rn == '.' {
			lex.advance()
			switch lex.rn {
			case '=':
				lex.advance()
				return lex.tok(RANGE_INCL)
			case '<':
				lex.advance()
				return lex.tok(RANGE_HALF_OPEN)
			case '.':
				lex.advance()
				return lex.tok(OP_ARRAY)
			}
			return lex.tok(OP_RANGE)
		}
		return lex.tok(DOT)

	case '<':
		switch lex.rn {
		case '<':
			lex.advance()
			return lex.tok(SHR_LESS)
		case '=':
			lex.advance()
			return lex.tok(LE)
		case '-':
			return lex.tok(CHAN_SEND)
		}
		return lex.tok(LT)

	case '>':
		switch lex.rn {
		case '>':
			lex.advance()
			return lex.tok(SHR_GREATER)
		case '=':
			lex.advance()
			return lex.tok(GE)
		}
		return lex.tok(GT)

	case '-':
		switch lex.rn {
		case '=':
			lex.advance()
			return lex.tok(SUB_ASSIGN)
		case '>':
			lex.advance()
			return lex.tok(RTN_ARROW)
		}
		return lex.tok(SUB)

	case '+':
		if lex.rn == '=' {
			lex.advance()
			return lex.tok(ADD_ASSIGN)
		}
		return lex.tok(ADD)

	case '*':
		if lex.rn == '=' {
			lex.advance()
			return lex.tok(MUL_ASSIGN)
		}
		return lex.tok(MUL)

	case '%':
		if lex.rn == '=' {
			lex.advance()
			return lex.tok(MOD_ASSIGN)
		}
		return lex.tok(MOD)

	case '^':
		if lex.rn == '=' {
			lex.advance()
			return lex.tok(POW_ASSIGN)
		}
		return lex.tok(POW)

	case '~':
		return lex.tok(BITWISE_NOT)

	case '&':
		if lex.rn == '&' {
			lex.advance()
			return lex.tok(AND)
		}
		return lex.tok(REF)

	case '!':
		if lex.rn == '=' {
			lex.advance()
			return lex.tok(NEQ)
		}
		return lex.tok(BANG)

	case '?':
		switch lex.rn {
		case ':':
			lex.advance()
			return lex.tok(DEFAULT)
		case '.':
			lex.advance()
			return lex.tok(OPTIONAL_DOT)
		}
		return lex.tok(ILLEGAL)

	case '=':
		switch lex.rn {
		case '=':
			lex.advance()
			return lex.tok(EQ)
		case '>':
			lex.advance()
			return lex.tok(ARROW)
		}
		return lex.tok(ASSIGN)

	case '|':
		switch lex.rn {
		case '|':
			lex.advance()
			return lex.tok(OR)
		case '>':
			lex.advance()
			return lex.tok(PIPE)
		}
		return lex.tok(ILLEGAL)

	case ':':
		return lex.tok(COLON)
	case '(':
		return lex.tok(L_PAREN)
	case ')':
		return lex.tok(R_PAREN)
	case '{':
		return lex.tok(L_BRACE)
	case '}':
		return lex.tok(R_BRACE)
	case '[':
		return lex.tok(L_BRACK)
	case ']':
		return lex.tok(R_BRACK)
	case ';':
		return lex.tok(END)
	case ',':
		return lex.tok(COMMA)

	case '"':
		return lex.readString()
	case '`':
		return lex.readRawString()
	case '\'':
		return lex.readChar()

	default:
		if ch >= '0' && ch <= '9' {
			return lex.readNumber(ch)
		}

		if isIdentStart(ch) {
			return lex.readIdent()
		}
		return lex.tok(ILLEGAL)
	}
}

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	}
	return r > 0x7F && isSpaceUnicode(r)
}

func isSpaceUnicode(r rune) bool {
	switch r {
	case 0x00A0,
		0x1680,
		0x2000, 0x2001, 0x2002, 0x2003,
		0x2004, 0x2005, 0x2006, 0x2007,
		0x2008, 0x2009, 0x200A,
		0x2028, 0x2029,
		0x202F, 0x205F,
		0x3000,
		0xFEFF:
		return true
	}
	return false
}

func isIdentStart(r rune) bool {
	return r == '_' ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r > 0x7F && isLetterUnicode(r)
}

func isIdentContinue(r rune) bool {
	return r == '_' ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r > 0x7F && (isLetterUnicode(r) || isDigitUnicode(r))
}

func isLetterUnicode(r rune) bool {
	return r >= 0xAA
}

func isDigitUnicode(r rune) bool {
	return r >= 0x0660
}

func (lex *Lexer) lineComment() Token {
	lex.advance() // '/'
	for lex.rn != '\n' && lex.rn != 0 {
		lex.advance()
	}
	return lex.tok(COMMENT)
}

func (lex *Lexer) multiLineComment() Token {
	lex.advance() // '*'
	lex.freeze()
	for {
		if lex.rn == 0 {
			tk := lex.tok(ILLEGAL)
			lex.da.Add(errors.Errors[2].Update(tk))
			lex.unfreeze()
			lex.stabilize()
			return tk
		}
		if lex.rn == '*' && lex.peek() == '/' {
			lex.advance()
			lex.advance()
			break
		}
		lex.advance()
	}
	return lex.tok(M_COMMENT)
}

func (lex *Lexer) readString() Token {
	isTemplate := false
	lex.freeze()

	for lex.rn != '"' && lex.rn != 0 {
		if lex.rn == '\\' {
			lex.advance()
			lex.advance()
			continue
		}
		if lex.rn == '$' && lex.peek() == '{' {
			isTemplate = true
		}
		lex.advance()
	}

	if lex.rn == 0 {
		tk := lex.tok(ILLEGAL)
		lex.da.Add(errors.Errors[2].Update(tk))
		lex.unfreeze()
		lex.stabilize()
		return tk
	}

	lex.advance() // '"'

	if isTemplate {
		return lex.tok(T_STRING)
	}
	return lex.tok(STRING)
}

func (lex *Lexer) readRawString() Token {
	lex.freeze()
	for lex.rn != '`' && lex.rn != 0 {
		lex.advance()
	}
	if lex.rn == 0 {
		tk := lex.tok(ILLEGAL)
		lex.da.Add(errors.Errors[2].Update(tk))
		lex.unfreeze()
		lex.stabilize()
		return tk
	}
	lex.advance()
	return lex.tok(RAW_STRING)
}

func (lex *Lexer) readChar() Token {
	if lex.rn == '\\' {
		lex.advance()
		lex.advance()
	} else if lex.rn != '\'' && lex.rn != 0 {
		lex.advance()
	}
	if lex.rn != '\'' {
		return lex.tok(ILLEGAL)
	}
	lex.advance()
	return lex.tok(CHARACTER)
}

func (lex *Lexer) readIdent() Token {
	for isIdentContinue(lex.rn) {
		lex.advance()
	}
	lit := lex.Input[lex.tokPos:lex.pos]
	return lex.tok(SearchKeyword(lit))
}

func (lex *Lexer) readNumber(first rune) Token {
	_ = first
	isFloat := false
	isIdent := false

	for {
		ch := lex.rn
		if ch >= '0' && ch <= '9' {
			lex.advance()
			continue
		}
		if ch == '.' {
			if isIdent || isFloat {
				break
			}
			next := lex.peek()
			if next < '0' || next > '9' {
				break
			}
			isFloat = true
			lex.advance()
			continue
		}
		if isIdentContinue(ch) {
			isIdent = true
			lex.advance()
			continue
		}
		break
	}

	lit := lex.Input[lex.tokPos:lex.pos]

	if isIdent {
		n := len(lit)
		if n >= 2 && lit[n-1] == 'i' {
			onlyDigits := true
			for _, b := range lit[:n-1] {
				if (b < '0' || b > '9') && b != '.' {
					onlyDigits = false
					break
				}
			}
			if onlyDigits {
				return lex.tok(IMAGINARY)
			}
		}
		return lex.tok(IDENTIFIER)
	}

	if isFloat {
		return lex.tok(FLOATING)
	}
	return lex.tok(INTEGER)
}

func (lex *Lexer) stabilize() {
	k := map[Kind]bool{
		FN: true, IF: true, ELSE: true, SWITCH: true,
		CASE: true, RETURN: true, ENUM: true, SELECT: true,
	}
	for {
		lex.freeze()
		tk := lex.NextToken()
		switch {
		case tk.Kind == EOF:
			return
		case tk.Kind == SPACING || tk.Kind == COMMENT || tk.Kind == M_COMMENT:
			continue
		case k[tk.Kind], tk.Kind == R_BRACE, tk.Kind == END:
			lex.unfreeze()
			return
		}
	}
}
