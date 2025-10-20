// Package unpack provides functionality for simple unpacking compressed strings
// Supports stringd like "a4bc2d5e" which become "aaaabccddddde" after unpacking
// Supporst escape sequences with backslashes
package unpack

import (
	"fmt"
	"unicode"
	"strconv"
	"strings"
)



// TokenType represents the type of token
type TokenType string

const (
	// Symbol represents a single character token
	Symbol TokenType = "Symbol"

	// Coefficient represents a numeric coefficient token  
	Coefficient TokenType = "Coefficient"

	// Escape represents an escape sequence token
	Escape TokenType = "Escape"

	// EOF represents end of string token
	EOF TokenType = "EOF"

	// Empty represents an empty token for error state
	Empty TokenType = ""
)



// Token represents a lexical token
type Token struct {
	tokenType TokenType
	value string
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.tokenType, t.value)
}



// LexerError represents an error during lexical analysis
type LexerError struct {
    Position int
    Char     rune
    Message  string
}

func (e *LexerError) Error() string {
	return fmt.Sprintf("Lexer error at position %d (char '%c'): %s", 
		e.Position, e.Char, e.Message)
}



// Lexer provides functionality to tokenize strings
type Lexer struct {
	Runes []rune
	Position int
	CurrentChar rune
}

func newLexer(str string) *Lexer {
	runes := []rune(str)
    if len(runes) == 0 {
		return &Lexer{Runes: runes, Position: -1, CurrentChar: 0}
    }
    return &Lexer{Runes: runes, Position: 0, CurrentChar: runes[0]}
}

func (lex *Lexer) advance() {
	lex.Position++
	if lex.Position > len(lex.Runes) - 1 {
		lex.CurrentChar = 0
	} else {
		lex.CurrentChar = lex.Runes[lex.Position]
	}
}

func (lex *Lexer) integer() string {
	runes := []rune{}
	for lex.CurrentChar != 0 && unicode.IsDigit(lex.CurrentChar) {
		runes = append(runes, lex.CurrentChar)
		lex.advance()
	}
	return string(runes)
}

func (lex *Lexer) getNextToken() (Token, error) {
	for lex.CurrentChar != 0 {
		if unicode.IsDigit(lex.CurrentChar) {
			coef := lex.integer()
			return Token{tokenType: Coefficient, value: coef}, nil
		}

		if unicode.IsLetter(lex.CurrentChar) {
			symb := lex.CurrentChar
			lex.advance()
			return Token{tokenType: Symbol, value: string(symb)}, nil
		}

		if lex.CurrentChar == '\\' {
            lex.advance()
            if lex.CurrentChar == 0 {
                return Token{tokenType: Empty, value: ""}, 
						&LexerError{
							Position: lex.Position,
							Char:     '\\',
							Message:  "Unexpected end of string after escape character",
						}
            }
            symb := lex.CurrentChar
            lex.advance()
            return Token{tokenType: Symbol, value: string(symb)}, nil
        }

		return Token{tokenType: Empty, value: ""}, 
				&LexerError{
					Position: lex.Position,
					Char: lex.CurrentChar,
					Message: "Unexpected character",
				}
	}

	return Token{tokenType: EOF, value: ""}, nil
}

func tokenize(str string) ([]Token, error) {
	lex := newLexer(str)
	tokens := []Token{}
	for {
		token, err := lex.getNextToken();
		
		if err != nil {
			return nil, err
		}

		if token.tokenType == EOF {
			break
		}

		tokens = append(tokens, token)
	}
	return tokens, nil
}


// Unpack unpacks a compressed string
func Unpack(str string) (string, error) {
	result := ""

	tokens, err := tokenize(str)

	if err != nil {
		return "", fmt.Errorf("Tokenize Error: %w", err)
	}

	if len(tokens) == 0 {
        return "", nil
    }

	if tokens[0].tokenType == Coefficient {
        return "", fmt.Errorf("Invalid String: starts with coefficient")
    }

	hasSymbol := false
    for _, token := range tokens {
        if token.tokenType == Symbol {
            hasSymbol = true
            break
        }
    }
    if !hasSymbol {
        return "", fmt.Errorf("Invalid String: only coefficients")
    }

	currentSymbol := ""
	currentCoefficient := 0
	for _, token := range tokens {
		if token.tokenType == Symbol {
			if currentSymbol != "" {
				result += currentSymbol
			}
			currentSymbol = token.value
			continue
		}

		if token.tokenType == Coefficient {
			currentCoefficient, err = strconv.Atoi(token.value)
			if err != nil {
				return "", fmt.Errorf("Invalid Coefficient: %w", err)
            }

			if currentSymbol == "" {
                return "", fmt.Errorf("Invalid String: coefficient without preceding symbol")
            }

			result += strings.Repeat(currentSymbol, currentCoefficient)
			currentSymbol = ""
			currentCoefficient = 0
		}
	}

	if currentSymbol != "" {
        result += currentSymbol
    }

	return result, nil
}