package interpreter

import (
	"bytes"
	"strings"
)

// Node interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement interface for statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression interface for expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents an identifier
type Identifier struct {
	Token Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// VarDecl represents a variable declaration statement
type VarDecl struct {
	Token    Token // the token.TYPE_* token (int, float, etc)
	Type     TokenType
	Name     *Identifier
	Value    Expression
	RangeMin Expression // Optional: e.g. int(0, 100)
	RangeMax Expression // Optional
}

func (vd *VarDecl) statementNode()       {}
func (vd *VarDecl) TokenLiteral() string { return vd.Token.Literal }
func (vd *VarDecl) String() string {
	var out bytes.Buffer

	out.WriteString(vd.Token.Literal)

	// Add range info if present
	if vd.RangeMin != nil && vd.RangeMax != nil {
		out.WriteString("(")
		out.WriteString(vd.RangeMin.String())
		out.WriteString(", ")
		out.WriteString(vd.RangeMax.String())
		out.WriteString(")")
	}

	out.WriteString(" ")
	out.WriteString(vd.Name.String())

	if vd.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vd.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// AssignStmt represents an assignment statement
type AssignStmt struct {
	Token Token // the token.ASSIGN token
	Name  *Identifier
	Value Expression
}

func (as *AssignStmt) statementNode()       {}
func (as *AssignStmt) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStmt) String() string {
	var out bytes.Buffer

	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExprStmt represents an expression statement
type ExprStmt struct {
	Token      Token // the first token of the expression
	Expression Expression
}

func (es *ExprStmt) statementNode()       {}
func (es *ExprStmt) TokenLiteral() string { return es.Token.Literal }
func (es *ExprStmt) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BinaryExpr represents a binary expression
type BinaryExpr struct {
	Token    Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (bin *BinaryExpr) expressionNode()      {}
func (bin *BinaryExpr) TokenLiteral() string { return bin.Token.Literal }
func (bin *BinaryExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(bin.Left.String())
	out.WriteString(" " + bin.Operator + " ")
	out.WriteString(bin.Right.String())
	out.WriteString(")")

	return out.String()
}

// UnaryExpr represents a unary expression e.g. -5, !true
type UnaryExpr struct {
	Token    Token // The prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (ue *UnaryExpr) expressionNode()      {}
func (ue *UnaryExpr) TokenLiteral() string { return ue.Token.Literal }
func (ue *UnaryExpr) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ue.Operator)
	out.WriteString(ue.Right.String())
	out.WriteString(")")

	return out.String()
}

// CallExpr represents a function call
type CallExpr struct {
	Token     Token      // The '(' token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpr) expressionNode()      {}
func (ce *CallExpr) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpr) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// IntegerLiteral represents an integer
type IntegerLiteral struct {
	Token Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a float
type FloatLiteral struct {
	Token Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// BooleanLiteral represents a boolean
type BooleanLiteral struct {
	Token Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// BlockStmt represents a block of statements
type BlockStmt struct {
	Token      Token // the { token
	Statements []Statement
}

func (bs *BlockStmt) statementNode()       {}
func (bs *BlockStmt) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// IfStmt represents an if statement
type IfStmt struct {
	Token       Token // the 'if' or 'ifrand' token
	Condition   Expression
	Consequence *BlockStmt
	Alternative Statement // could be another IfStmt (for else if) or BlockStmt (for else)
}

func (is *IfStmt) statementNode()       {}
func (is *IfStmt) TokenLiteral() string { return is.Token.Literal }
func (is *IfStmt) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	if is.Condition != nil {
		out.WriteString(" ")
		out.WriteString(is.Condition.String())
	}
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}
