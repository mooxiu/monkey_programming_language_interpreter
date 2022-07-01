package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}

	p.nextToken()
	p.nextToken()
	return p
}

// move to next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	// iterate through the whole input, until current position is token.EOF
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// decide the kind of statement
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
		Name:  &ast.Identifier{},
		Value: nil,
	}

	// parse the name
	// the ident to assign ... must ... of course ... is ident
	nextIdent := p.peekToken
	if nextIdent.Type != token.IDENT {
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: nextIdent,
		Value: nextIdent.Literal,
	}
	// go to next
	p.nextToken()

	// parse the assign token -> "="
	if p.peekToken.Type != token.ASSIGN {
		return nil
	}
	p.nextToken()

	// parse the value: expression
	// TODO: jump this part

	for p.peekToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}
