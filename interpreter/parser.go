package interpreter

import (
	"strconv"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[TokenType]int{
	PLUS:     SUM,
	MINUS:    SUM,
	SLASH:    PRODUCT,
	ASTERISK: PRODUCT,
	LPAREN:   CALL,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l      *Lexer
	errors []*ParserError

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: make([]*ParserError, 0),
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(INT, p.parseIntegerLiteral)
	p.registerPrefix(FLOAT, p.parseFloatLiteral)
	p.registerPrefix(TRUE, p.parseBoolean)
	p.registerPrefix(FALSE, p.parseBoolean)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(LPAREN, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		if p.curToken.Type == ILLEGAL {
			p.errors = append(p.errors, NewIllegalTokenError(&p.curToken))
		}
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	var errs []string
	for _, e := range p.errors {
		errs = append(errs, e.Msg)
	}
	return errs
}

func (p *Parser) peekError(t TokenType) {
	p.errors = append(p.errors, NewExpectedTokenError(&p.peekToken, t))
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case TYPE_INT, TYPE_UINT, TYPE_FLOAT, TYPE_UNOFLOAT, TYPE_BOOL, TYPE_STRING:
		return p.parseVarStatement()
	case IDENT:
		// Could be an assignment or an expression statement
		// If peek is ASSIGN, it's an assignment
		if p.peekToken.Type == ASSIGN {
			return p.parseAssignStatement()
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() Statement {
	stmt := &VarDecl{Token: p.curToken, Type: p.curToken.Type}

	// Check for optional range: type(min, max) name
	if p.peekToken.Type == LPAREN {
		p.nextToken() // consume type
		p.nextToken() // consume '('

		stmt.RangeMin = p.parseExpression(LOWEST)

		if !p.expectPeek(COMMA) {
			return nil
		}

		p.nextToken() // consume comma
		stmt.RangeMax = p.parseExpression(LOWEST)

		if !p.expectPeek(RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional assignment: = value
	if p.peekToken.Type == ASSIGN {
		p.nextToken() // consume name
		p.nextToken() // consume '='
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// TODO: check it out again
func (p *Parser) parseAssignStatement() Statement {
	stmt := &AssignStmt{Name: &Identifier{Token: p.curToken, Value: p.curToken.Literal}}

	p.nextToken() // consume identifier
	// The `parseStatement` function dispatches to `parseAssignStatement` only when `p.peekToken.Type` is `ASSIGN`.
	// After `p.nextToken()` is called, `p.curToken` should be `ASSIGN`.
	// This check acts as a safeguard, as the dispatch logic in `parseStatement` should ensure `p.curToken.Type` is `ASSIGN` here.
	if p.curToken.Type != ASSIGN {
		return nil // This state indicates an internal parser error or an unexpected token sequence.
	}
	stmt.Token = p.curToken // The '=' token

	p.nextToken() // consume '='
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExprStmt {
	stmt := &ExprStmt{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, NewIntegerParseError(&p.curToken))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, NewFloatParseError(&p.curToken))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curToken.Type == TRUE}
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &UnaryExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &BinaryExpr{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpr{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == COMMA {
		p.nextToken() // consume comma
		p.nextToken() // consume expression
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	p.errors = append(p.errors, NewNoPrefixParseFnError(&p.curToken, t))
}
