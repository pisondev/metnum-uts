package rootfinding

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Expression struct {
	raw  string
	root exprNode
}

func CompileExpression(raw string) (Expression, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Expression{}, fmt.Errorf("ekspresi tidak boleh kosong")
	}

	tokens, err := tokenize(raw)
	if err != nil {
		return Expression{}, err
	}

	parser := parser{tokens: tokens}
	root, err := parser.parseExpression()
	if err != nil {
		return Expression{}, err
	}
	if parser.current().typ != tokenEOF {
		return Expression{}, fmt.Errorf("token tidak dikenali di dekat %q", parser.current().text)
	}

	return Expression{raw: raw, root: root}, nil
}

func (e Expression) Eval(x float64) (float64, error) {
	if e.root == nil {
		return 0, fmt.Errorf("ekspresi belum dikompilasi")
	}
	return e.root.Eval(x)
}

type exprNode interface {
	Eval(x float64) (float64, error)
}

type numberNode struct {
	value float64
}

func (n numberNode) Eval(_ float64) (float64, error) {
	return n.value, nil
}

type variableNode struct{}

func (variableNode) Eval(x float64) (float64, error) {
	return x, nil
}

type unaryNode struct {
	op    tokenType
	right exprNode
}

func (n unaryNode) Eval(x float64) (float64, error) {
	value, err := n.right.Eval(x)
	if err != nil {
		return 0, err
	}

	switch n.op {
	case tokenPlus:
		return value, nil
	case tokenMinus:
		return -value, nil
	default:
		return 0, fmt.Errorf("operator unary tidak didukung")
	}
}

type binaryNode struct {
	op          tokenType
	left, right exprNode
}

func (n binaryNode) Eval(x float64) (float64, error) {
	leftValue, err := n.left.Eval(x)
	if err != nil {
		return 0, err
	}

	rightValue, err := n.right.Eval(x)
	if err != nil {
		return 0, err
	}

	var result float64
	switch n.op {
	case tokenPlus:
		result = leftValue + rightValue
	case tokenMinus:
		result = leftValue - rightValue
	case tokenMultiply:
		result = leftValue * rightValue
	case tokenDivide:
		if NearlyEqual(rightValue, 0) {
			return 0, fmt.Errorf("terjadi pembagian dengan nol")
		}
		result = leftValue / rightValue
	case tokenPower:
		result = math.Pow(leftValue, rightValue)
	default:
		return 0, fmt.Errorf("operator biner tidak didukung")
	}

	if !IsFinite(result) {
		return 0, fmt.Errorf("hasil evaluasi tidak hingga")
	}

	return result, nil
}

type functionNode struct {
	name string
	arg  exprNode
}

func (n functionNode) Eval(x float64) (float64, error) {
	value, err := n.arg.Eval(x)
	if err != nil {
		return 0, err
	}

	var result float64
	switch n.name {
	case "sin":
		result = math.Sin(value)
	case "cos":
		result = math.Cos(value)
	case "tan":
		result = math.Tan(value)
	case "exp":
		result = math.Exp(value)
	case "sqrt":
		result = math.Sqrt(value)
	case "abs":
		result = math.Abs(value)
	case "ln":
		result = math.Log(value)
	case "log":
		result = math.Log10(value)
	default:
		return 0, fmt.Errorf("fungsi %q tidak didukung", n.name)
	}

	if !IsFinite(result) {
		return 0, fmt.Errorf("hasil evaluasi fungsi %s tidak hingga", n.name)
	}

	return result, nil
}

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenNumber
	tokenIdentifier
	tokenPlus
	tokenMinus
	tokenMultiply
	tokenDivide
	tokenPower
	tokenLeftParen
	tokenRightParen
)

type token struct {
	typ   tokenType
	text  string
	value float64
}

func tokenize(input string) ([]token, error) {
	tokens := make([]token, 0, len(input)+1)

	for i := 0; i < len(input); {
		ch := input[i]

		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			i++
			continue
		}

		switch ch {
		case '+':
			tokens = append(tokens, token{typ: tokenPlus, text: "+"})
			i++
			continue
		case '-':
			tokens = append(tokens, token{typ: tokenMinus, text: "-"})
			i++
			continue
		case '*':
			tokens = append(tokens, token{typ: tokenMultiply, text: "*"})
			i++
			continue
		case '/':
			tokens = append(tokens, token{typ: tokenDivide, text: "/"})
			i++
			continue
		case '^':
			tokens = append(tokens, token{typ: tokenPower, text: "^"})
			i++
			continue
		case '(':
			tokens = append(tokens, token{typ: tokenLeftParen, text: "("})
			i++
			continue
		case ')':
			tokens = append(tokens, token{typ: tokenRightParen, text: ")"})
			i++
			continue
		}

		if isDigit(ch) || ch == '.' {
			start := i
			dotSeen := false
			expSeen := false

			for i < len(input) {
				curr := input[i]
				switch {
				case isDigit(curr):
					i++
				case curr == '.' && !dotSeen && !expSeen:
					dotSeen = true
					i++
				case (curr == 'e' || curr == 'E') && !expSeen:
					expSeen = true
					i++
					if i < len(input) && (input[i] == '+' || input[i] == '-') {
						i++
					}
					if i >= len(input) || !isDigit(input[i]) {
						return nil, fmt.Errorf("eksponen tidak valid di dekat %q", input[start:i])
					}
				default:
					goto doneNumber
				}
			}

		doneNumber:
			text := input[start:i]
			value, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, fmt.Errorf("angka tidak valid: %q", text)
			}
			tokens = append(tokens, token{typ: tokenNumber, text: text, value: value})
			continue
		}

		if isLetter(ch) {
			start := i
			for i < len(input) && isIdentifierChar(input[i]) {
				i++
			}
			text := strings.ToLower(input[start:i])
			tokens = append(tokens, token{typ: tokenIdentifier, text: text})
			continue
		}

		return nil, fmt.Errorf("karakter tidak dikenali: %q", string(ch))
	}

	tokens = append(tokens, token{typ: tokenEOF})
	return tokens, nil
}

type parser struct {
	tokens []token
	pos    int
}

func (p *parser) parseExpression() (exprNode, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		curr := p.current()
		if curr.typ != tokenPlus && curr.typ != tokenMinus {
			return left, nil
		}

		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		left = binaryNode{op: curr.typ, left: left, right: right}
	}
}

func (p *parser) parseTerm() (exprNode, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for {
		curr := p.current()
		if curr.typ != tokenMultiply && curr.typ != tokenDivide {
			return left, nil
		}

		p.pos++
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		left = binaryNode{op: curr.typ, left: left, right: right}
	}
}

func (p *parser) parseUnary() (exprNode, error) {
	curr := p.current()
	if curr.typ == tokenPlus || curr.typ == tokenMinus {
		p.pos++
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return unaryNode{op: curr.typ, right: right}, nil
	}

	return p.parsePower()
}

func (p *parser) parsePower() (exprNode, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	if p.current().typ != tokenPower {
		return left, nil
	}

	p.pos++
	right, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	return binaryNode{op: tokenPower, left: left, right: right}, nil
}

func (p *parser) parsePrimary() (exprNode, error) {
	curr := p.current()
	switch curr.typ {
	case tokenNumber:
		p.pos++
		return numberNode{value: curr.value}, nil
	case tokenIdentifier:
		p.pos++
		if p.current().typ == tokenLeftParen {
			functionName := curr.text
			p.pos++
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			if p.current().typ != tokenRightParen {
				return nil, fmt.Errorf("kurung tutup tidak ditemukan setelah fungsi %s", functionName)
			}
			p.pos++
			return functionNode{name: functionName, arg: arg}, nil
		}

		switch curr.text {
		case "x":
			return variableNode{}, nil
		case "pi":
			return numberNode{value: math.Pi}, nil
		case "e":
			return numberNode{value: math.E}, nil
		default:
			return nil, fmt.Errorf("identifier %q tidak dikenali", curr.text)
		}
	case tokenLeftParen:
		p.pos++
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current().typ != tokenRightParen {
			return nil, fmt.Errorf("kurung tutup tidak ditemukan")
		}
		p.pos++
		return node, nil
	default:
		return nil, fmt.Errorf("token %q tidak valid pada posisi ini", curr.text)
	}
}

func (p *parser) current() token {
	if p.pos >= len(p.tokens) {
		return token{typ: tokenEOF}
	}
	return p.tokens[p.pos]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isIdentifierChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}
