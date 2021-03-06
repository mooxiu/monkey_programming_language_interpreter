package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int  // current position
	readPosition int  // next position
	ch           byte // current char
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// already on the last char of the input
		l.ch = 0 // 0 is the end of the file
	} else {
		l.ch = l.input[l.readPosition]
	}
	// move to the next position
	l.position = l.readPosition
	l.readPosition += 1
	return
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			thisChar := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.EQ,
				Literal: string(thisChar) + string(l.ch),
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			thisChar := l.ch
			l.readChar()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(thisChar) + string(l.ch),
			}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	// here situation can be difficult
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tokType, ok := token.IdentsToken[tok.Literal]
			if !ok {
				tok.Type = token.IDENT
			} else {
				tok.Type = tokType
			}
		} else if isDigit(l.ch) {
			tok.Literal = l.readNum()
			// TODO: add more number type in the future
			tok.Type = token.INT
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	originalPosition := l.position
	// move until next position is not letter
	// TODO: actually identifier can end up with number: like `s1`
	for isLetter(l.peekChar()) {
		l.readChar()
	}
	return l.input[originalPosition:l.readPosition]
}

// move to next until l.ch is not " "
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
	return
}

/*
	isWhitespace decide whether this char can be ignored
	ref: https://ja.wikipedia.org/wiki/ASCII

	TODO: optimize the sort
*/
func isWhitespace(ch byte) bool {
	return ch == '\t' || ch == '\n' || ch == '\r' || ch == ' '
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNum() string {
	originalPosition := l.position
	// move until next position is not digit
	for isDigit(l.peekChar()) {
		l.readChar()
	}
	return l.input[originalPosition:l.readPosition]
}

// see next position char
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
