package lexer

import (
	"errors"
	"fmt"
	"github.com/Dor1ma/Strawberry/token"
	"os"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

var eof = rune(-1)

var (
	// ошибки для идентификатора
	errUnterminated = errors.New("unterminated string")
	errEspace       = errors.New("invalid escape char")
	errInvalidChar  = errors.New("invalid unicode char")

	// ошибки для чисел
	errLessPower = errors.New("power is required")
)

// Lexer - делает лексический анализ текста и разбивает его на токены
type Lexer struct {
	s        *scanner.Scanner
	char     rune
	tokenBuf *strings.Builder
}

func (l *Lexer) consume() {
	if l.isAtEnd() {
		return
	}
	char := l.s.Next()
	if char == scanner.EOF {
		l.char = eof
		return
	}
	l.char = char
}

func (l *Lexer) peek() rune {
	ch := l.s.Peek()
	return ch
}

func (l *Lexer) skip() {
	for unicode.IsSpace(l.char) {
		l.consume()
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.char == eof
}

func (l *Lexer) match(ch rune) bool {
	l.consume()
	if l.isAtEnd() || l.char != ch {
		return false
	}
	l.consume()
	return true
}

func (l *Lexer) error(msg string) {
	// ToDo: исправить баг с позицией
	fmt.Fprintf(os.Stderr, "%s %s\n", l.Pos().String(), msg)
}

func (l *Lexer) readIdentifier() string {
	l.tokenBuf.Reset()
	for isAlphaNumeric(l.char) {
		l.tokenBuf.WriteRune(l.char)
		l.consume()
	}
	return l.tokenBuf.String()
}

func (l *Lexer) readString() (string, error) {
	l.tokenBuf.Reset()
	l.consume()
	if l.char == '"' {
		l.consume()
		return "", nil
	}

	for l.char != '"' {
		if l.isAtEnd() {
			l.error(errUnterminated.Error())
			return "", errUnterminated
		} else if l.char == '\\' {
			peekCh := l.peek()
			if peekCh == eof {
				l.error(errEspace.Error())
				return "", errEspace
			}
			l.consume()
			switch peekCh {
			case '"':
				l.tokenBuf.WriteRune('"')
			case 'u':
				code := make([]rune, 4)
				for i := range code {
					l.consume()
					if !unicode.Is(unicode.Hex_Digit, l.char) {
						l.error(errInvalidChar.Error())
						return "", errInvalidChar
					}
					code[i] = l.char
				}
				l.tokenBuf.WriteRune(charCode2Rune(string(code)))
			}
		} else {
			l.tokenBuf.WriteRune(l.char)
		}
		l.consume()
	}
	// end ".
	l.consume()
	return l.tokenBuf.String(), nil
}

func (l *Lexer) readNumber() (string, error) {

	l.tokenBuf.Reset()
	for unicode.IsNumber(l.char) {
		l.tokenBuf.WriteRune(l.char)
		l.consume()
	}

	if l.char == '.' {
		if !unicode.IsNumber(l.peek()) {
			return l.tokenBuf.String(), nil
		}
		l.tokenBuf.WriteRune(l.char)
		l.consume()
		for unicode.IsNumber(l.char) {
			l.tokenBuf.WriteRune(l.char)
			l.consume()
		}
	}

	if l.char == 'E' || l.char == 'e' {
		seenPower := false
		l.tokenBuf.WriteRune(l.char)
		l.consume()
		if l.char == '+' || l.char == '-' {
			l.tokenBuf.WriteRune(l.char)
			l.consume()
		}
		for unicode.IsNumber(l.char) {
			seenPower = true
			l.tokenBuf.WriteRune(l.char)
			l.consume()
		}
		if !seenPower {
			l.error(errLessPower.Error())
			return "", errLessPower
		}
	}

	return l.tokenBuf.String(), nil
}

// NextToken читает и возвращает токены или литералы
func (l *Lexer) NextToken() (tok token.Token, literal string) {
	l.skip()

	switch l.char {
	case '(':
		tok = token.LeftParen
		literal = "("
	case ')':
		tok = token.RightParen
		literal = ")"
	case '[':
		tok = token.LeftBracket
		literal = "["
	case ']':
		tok = token.RightBracket
		literal = "]"
	case '{':
		tok = token.LeftBrace
		literal = "{"
	case '}':
		tok = token.RightBrace
		literal = "}"
	case ',':
		tok = token.Comma
		literal = ","
	case '.':
		tok = token.Dot
		literal = "."
	case '-':
		tok = token.Minus
		literal = "-"
	case '+':
		tok = token.Plus
		literal = "+"
	case ';':
		tok = token.Semicolon
		literal = ";"
	case '/':
		tok = token.Slash
		literal = "/"
	case '*':
		tok = token.Star
		literal = "*"
	case '!':
		if l.match('=') {
			tok = token.BangEqual
			literal = "!="
		} else {
			tok = token.Bang
			literal = "!"
		}
		return
	case '=':
		if l.match('=') {
			tok = token.EqualEqual
			literal = "=="
		} else {
			tok = token.Equal
			literal = "="
		}
		return
	case '>':
		if l.match('=') {
			tok = token.GreaterThanOrEqual
			literal = ">="
		} else {
			tok = token.Greater
			literal = ">"
		}
		return
	case '<':
		if l.match('=') {
			tok = token.LessThanOrEqual
			literal = "<="
		} else {
			tok = token.Less
			literal = "<"
		}
		return
	case '"':
		liter, err := l.readString()
		if err != nil {
			return token.Illegal, liter
		}
		tok = token.String
		literal = liter
		return
	case eof:
		tok = token.EOF
		return
	default:
		if unicode.IsLetter(l.char) {
			literal = l.readIdentifier()
			tok = token.Lookup(literal)
			return
		} else if unicode.IsNumber(l.char) {
			liter, err := l.readNumber()
			if err != nil {
				return token.Illegal, ""
			}
			tok = token.Number
			literal = liter
			return
		}

		tok = token.Illegal
		literal = ""
	}

	l.consume()
	return
}

func (l *Lexer) Pos() scanner.Position {
	return l.s.Pos()
}

func charCode2Rune(code string) rune {
	v, err := strconv.ParseInt(code, 16, 32)
	if err != nil {
		return unicode.ReplacementChar
	}
	return rune(v)
}

func isAlphaNumeric(ch rune) bool { return unicode.IsLetter(ch) || unicode.IsNumber(ch) || ch == '_' }

// конструктор
func New(input string) *Lexer {
	s := &scanner.Scanner{}
	s.Init(strings.NewReader(input))
	l := &Lexer{
		s:        s,
		tokenBuf: &strings.Builder{},
	}
	l.consume()
	return l
}
