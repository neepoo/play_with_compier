package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type DfaState int

const (
	DfaInitial DfaState = iota
	DfaIf
	DfaId_if1
	DfaId_if2
	DfaElse
	DfaId_else1
	DfaId_else2
	DfaId_else3
	DfaId_else4
	DfaInt
	DfaId_int1
	DfaId_int2
	DfaId_int3
	DfaId
	DfaGT
	DfaGE
	DfaAssignment
	DfaPlus
	DfaMinus
	DfaStar
	DfaSlash
	DfaSemiColon
	DfaLeftParen
	DfaRightParen
	DfaIntLiteral
)

type TokenType int

const (
	Initial TokenType = iota
	Plus              // +
	Minus             // -
	Star              // *
	Slash             // /

	GE // >=
	GT // >
	EQ // ==
	LE // <=
	LT // <

	SemiColon  // ;
	LeftParen  // (
	RightParen // )

	Assignment // =

	If
	Else

	Int

	Identifier //标识符

	IntLiteral    //整型字面量
	StringLiteral //字符串字面量
)

func (tk TokenType) String() string {
	if tk == Initial {
		return "Initial"
	} else if tk == Plus {
		return "Plus"
	} else if tk == Minus {
		return "Minus"
	} else if tk == Star {
		return "Star"
	} else if tk == Slash {
		return "Slash"
	} else if tk == GE {
		return "GE"
	} else if tk == GT {
		return "GT"
	} else if tk == EQ {
		return "EQ"
	} else if tk == LE {
		return "LE"
	} else if tk == LT {
		return "LT"
	} else if tk == SemiColon {
		return "SemiColon"
	} else if tk == LeftParen {
		return "LeftParen"
	} else if tk == RightParen {
		return "RightParen"
	} else if tk == Assignment {
		return "Assignment"
	} else if tk == If {
		return "If"
	} else if tk == Else {
		return "Else"
	} else if tk == Int {
		return "Int"
	} else if tk == Identifier {
		return "Identifier"
	} else if tk == IntLiteral {
		return "IntLiteral"
	} else if tk == StringLiteral {
		return "StringLiteral"
	} else {
		return "UNKNOWN TYPE"
	}

}

type Token struct {
	tokenType TokenType
	tokenText string
}

type Tokens []Token

func (ts Tokens) print() {
	fmt.Printf("%s\t%s\n", "text", "type")
	fmt.Println()
	for _, token := range ts {
		fmt.Printf("%s\t%s\n", token.tokenText, token.tokenType.String())
	}
}

type SimpleLexer struct {
	s            string
	tokenText    strings.Builder
	tokens       Tokens
	currentToken Token
}

func isAlpha(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsBlank(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func (s *SimpleLexer) initToken(ch rune) DfaState {
	if s.tokenText.Len() > 0 {
		s.currentToken.tokenText = s.tokenText.String()
		s.tokens = append(s.tokens, s.currentToken)

		s.tokenText = strings.Builder{}
		s.currentToken = Token{}
	}
	var newState DfaState = DfaInitial
	if isAlpha(ch) {
		if ch == 'i' {
			newState = DfaId_int1
		} else {
			newState = DfaId
		}
		s.currentToken.tokenType = Identifier
		s.tokenText.WriteRune(ch)
	} else if isDigit(ch) {
		newState = DfaIntLiteral
		s.currentToken.tokenType = IntLiteral
		s.tokenText.WriteRune(ch)
	} else if ch == '>' {
		newState = DfaGT
		s.currentToken.tokenType = GT
		s.tokenText.WriteRune(ch)
	} else if ch == '+' {
		newState = DfaPlus
		s.currentToken.tokenType = Plus
		s.tokenText.WriteRune(ch)
	} else if ch == '-' {
		newState = DfaMinus
		s.currentToken.tokenType = Minus
		s.tokenText.WriteRune(ch)
	} else if ch == '*' {
		newState = DfaStar
		s.currentToken.tokenType = Star
		s.tokenText.WriteRune(ch)
	} else if ch == '/' {
		newState = DfaSlash
		s.currentToken.tokenType = Slash
		s.tokenText.WriteRune(ch)
	} else if ch == ';' {
		newState = DfaSemiColon
		s.currentToken.tokenType = SemiColon
		s.tokenText.WriteRune(ch)
	} else if ch == '(' {
		newState = DfaLeftParen
		s.currentToken.tokenType = LeftParen
		s.tokenText.WriteRune(ch)
	} else if ch == ')' {
		newState = DfaRightParen
		s.currentToken.tokenType = RightParen
		s.tokenText.WriteRune(ch)
	} else if ch == '=' {
		newState = DfaAssignment
		s.currentToken.tokenType = Assignment
		s.tokenText.WriteRune(ch)
	} else {
		newState = DfaInitial
	}

	return newState
}

func (s *SimpleLexer) tokenize() Tokens {
	reader := bufio.NewReader(strings.NewReader(s.s))
	state := DfaInitial
	s.tokens = Tokens{}
	s.currentToken = Token{}
	s.tokenText = strings.Builder{}
	var ch rune
	var err error
	for {
		ch, _, err = reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("reader rune error : %v", err)
		}

		switch state {
		case DfaInitial:
			state = s.initToken(ch)
		case DfaId:
			if isAlpha(ch) || isDigit(ch) {

				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaGT:
			if ch == '=' {
				s.currentToken.tokenType = GE
				state = DfaGE
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaGE:
			if IsBlank(ch) {
				state = s.initToken(ch)
			}
		case DfaAssignment, DfaPlus, DfaMinus, DfaStar, DfaSlash, DfaSemiColon, DfaLeftParen:
			{
				continue
			}
		case DfaRightParen:
			state = s.initToken(ch)
		case DfaIntLiteral:
			if isDigit(ch) {
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaId_int1:
			if ch == 'n' {
				state = DfaId_int2
				s.tokenText.WriteRune(ch)
			} else if isDigit(ch) || isAlpha(ch) {
				state = DfaId
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaId_int2:
			if ch == 't' {
				state = DfaId_int3
				s.tokenText.WriteRune(ch)
			} else if isAlpha(ch) || isDigit(ch) {
				state = DfaId
				s.tokenText.WriteRune(ch)
			} else {
				state = s.initToken(ch)
			}
		case DfaId_int3:
			if IsBlank(ch) {
				s.currentToken.tokenType = Int
				state = s.initToken(ch)
			} else {
				state = DfaId
				s.tokenText.WriteRune(ch)
			}
		default:
			state = s.initToken(ch)
		}

	}
	if s.tokenText.Len() > 0 {
		s.initToken(ch)
	}
	return s.tokens
}

func NewSimpleLexer(text string) SimpleLexer {
	return SimpleLexer{
		s:            text,
		tokenText:    strings.Builder{},
		tokens:       []Token{},
		currentToken: Token{},
	}
}

func main() {

	ss := []string{"int age = 45;", "inta age = 45;", "in age = 45;", "age >= 45;", "age > 45;"}
	for _, s := range ss {
		fmt.Println(s)
		lexer := NewSimpleLexer(s)
		res := lexer.tokenize()
		res.print()
		fmt.Println()
	}
}
