package uql

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Token represents a lexical token
type Token struct {
	Type  string
	Value string
	Pos   int
}

// Lexer performs lexical analysis
type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		tokens: make([]Token, 0),
	}
}

// Tokenize converts input string into tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		// Skip whitespace
		if l.input[l.pos] == ' ' || l.input[l.pos] == '\t' || l.input[l.pos] == '\r' {
			l.pos++
			continue
		}

		// Skip newlines but track them
		if l.input[l.pos] == '\n' {
			l.tokens = append(l.tokens, Token{Type: "newline", Value: "\n", Pos: l.pos})
			l.pos++
			continue
		}

		// Comments
		if l.input[l.pos] == '#' {
			start := l.pos
			l.pos++ // skip #
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
			l.tokens = append(l.tokens, Token{Type: "comment", Value: l.input[start+1 : l.pos], Pos: start})
			continue
		}

		// Pipe
		if l.input[l.pos] == '|' {
			l.tokens = append(l.tokens, Token{Type: "pipe", Value: "|", Pos: l.pos})
			l.pos++
			continue
		}

		// Comma
		if l.input[l.pos] == ',' {
			l.tokens = append(l.tokens, Token{Type: "comma", Value: ",", Pos: l.pos})
			l.pos++
			continue
		}

		// Parentheses
		if l.input[l.pos] == '(' {
			l.tokens = append(l.tokens, Token{Type: "lparen", Value: "(", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == ')' {
			l.tokens = append(l.tokens, Token{Type: "rparen", Value: ")", Pos: l.pos})
			l.pos++
			continue
		}

		// Assignment
		if l.input[l.pos] == '=' {
			// Check for ==
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
				l.tokens = append(l.tokens, Token{Type: "eq", Value: "==", Pos: l.pos})
				l.pos += 2
				continue
			}
			l.tokens = append(l.tokens, Token{Type: "assignment", Value: "=", Pos: l.pos})
			l.pos++
			continue
		}

		// String with double quotes
		if l.input[l.pos] == '"' {
			start := l.pos
			l.pos++ // skip opening quote
			value := ""
			for l.pos < len(l.input) && l.input[l.pos] != '"' {
				if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
					// Handle escape sequences
					l.pos++
					switch l.input[l.pos] {
					case 'n':
						value += "\n"
					case 't':
						value += "\t"
					case 'r':
						value += "\r"
					case '\\':
						value += "\\"
					case '"':
						value += "\""
					case '\'':
						value += "'"
					default:
						value += string(l.input[l.pos])
					}
					l.pos++
				} else {
					value += string(l.input[l.pos])
					l.pos++
				}
			}
			if l.pos >= len(l.input) {
				return nil, fmt.Errorf("unterminated string at position %d", start)
			}
			l.pos++ // skip closing quote
			l.tokens = append(l.tokens, Token{Type: "string", Value: value, Pos: start})
			continue
		}

		// String with single quotes
		if l.input[l.pos] == '\'' {
			start := l.pos
			l.pos++ // skip opening quote
			value := ""
			for l.pos < len(l.input) && l.input[l.pos] != '\'' {
				if l.input[l.pos] == '\\' && l.pos+1 < len(l.input) {
					l.pos++
					switch l.input[l.pos] {
					case 'n':
						value += "\n"
					case 't':
						value += "\t"
					case 'r':
						value += "\r"
					case '\\':
						value += "\\"
					case '"':
						value += "\""
					case '\'':
						value += "'"
					default:
						value += string(l.input[l.pos])
					}
					l.pos++
				} else {
					value += string(l.input[l.pos])
					l.pos++
				}
			}
			if l.pos >= len(l.input) {
				return nil, fmt.Errorf("unterminated string at position %d", start)
			}
			l.pos++ // skip closing quote
			l.tokens = append(l.tokens, Token{Type: "sq_string", Value: value, Pos: start})
			continue
		}

		// Numbers - check if it's a digit or negative number
		if l.isDigit(l.input[l.pos]) || l.isNegativeNumber() {
			start := l.pos
			if l.input[l.pos] == '-' {
				l.pos++
			}
			// Parse integer and decimal parts
			for l.pos < len(l.input) && l.isNumberChar(l.input[l.pos]) {
				l.pos++
			}
			l.tokens = append(l.tokens, Token{Type: "number", Value: l.input[start:l.pos], Pos: start})
			continue
		}

		// Operators
		if l.pos+1 < len(l.input) {
			twoChar := l.input[l.pos : l.pos+2]
			switch twoChar {
			case ">=":
				l.tokens = append(l.tokens, Token{Type: "gte", Value: ">=", Pos: l.pos})
				l.pos += 2
				continue
			case "<=":
				l.tokens = append(l.tokens, Token{Type: "lte", Value: "<=", Pos: l.pos})
				l.pos += 2
				continue
			case "!=":
				l.tokens = append(l.tokens, Token{Type: "ne", Value: "!=", Pos: l.pos})
				l.pos += 2
				continue
			case "=~":
				l.tokens = append(l.tokens, Token{Type: "regex_match", Value: "=~", Pos: l.pos})
				l.pos += 2
				continue
			case "!~":
				l.tokens = append(l.tokens, Token{Type: "regex_not_match", Value: "!~", Pos: l.pos})
				l.pos += 2
				continue
			}
		}

		// Single char operators
		if l.input[l.pos] == '>' {
			l.tokens = append(l.tokens, Token{Type: "gt", Value: ">", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '<' {
			l.tokens = append(l.tokens, Token{Type: "lt", Value: "<", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '+' {
			l.tokens = append(l.tokens, Token{Type: "plus", Value: "+", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '-' {
			l.tokens = append(l.tokens, Token{Type: "dash", Value: "-", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '*' {
			l.tokens = append(l.tokens, Token{Type: "mul", Value: "*", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '/' {
			l.tokens = append(l.tokens, Token{Type: "divide", Value: "/", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '%' {
			l.tokens = append(l.tokens, Token{Type: "mod", Value: "%", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '!' {
			l.tokens = append(l.tokens, Token{Type: "exclaim", Value: "!", Pos: l.pos})
			l.pos++
			continue
		}
		if l.input[l.pos] == '~' {
			l.tokens = append(l.tokens, Token{Type: "tilde", Value: "~", Pos: l.pos})
			l.pos++
			continue
		}

		// Identifiers and keywords
		if (l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z') || (l.input[l.pos] >= 'A' && l.input[l.pos] <= 'Z') || l.input[l.pos] == '_' {
			start := l.pos
			for l.pos < len(l.input) && ((l.input[l.pos] >= 'a' && l.input[l.pos] <= 'z') || (l.input[l.pos] >= 'A' && l.input[l.pos] <= 'Z') || (l.input[l.pos] >= '0' && l.input[l.pos] <= '9') || l.input[l.pos] == '_' || l.input[l.pos] == '-') {
				l.pos++
			}
			value := l.input[start:l.pos]
			l.tokens = append(l.tokens, Token{Type: "identifier", Value: value, Pos: start})
			continue
		}

		return nil, fmt.Errorf("unexpected character '%c' at position %d", l.input[l.pos], l.pos)
	}

	return l.tokens, nil
}

// Parser parses tokens into commands
type Parser struct {
	tokens []Token
	pos    int
}

// NewParser creates a new parser
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Helper methods for lexer

// isDigit checks if a character is a digit
func (l *Lexer) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isNegativeNumber checks if current position is start of a negative number
func (l *Lexer) isNegativeNumber() bool {
	return l.input[l.pos] == '-' && l.pos+1 < len(l.input) && l.isDigit(l.input[l.pos+1])
}

// isNumberChar checks if a character is valid in a number (digit, dot, or exponent notation)
func (l *Lexer) isNumberChar(c byte) bool {
	return l.isDigit(c) || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-'
}

// Parse parses a UQL query string into commands
func Parse(query string) ([]Command, error) {
	if query == "" {
		return []Command{{Type: CmdHello}}, nil
	}

	lexer := NewLexer(strings.TrimSpace(query))
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	return parser.ParseCommands()
}

// ParseCommands parses multiple commands separated by pipes
func (p *Parser) ParseCommands() ([]Command, error) {
	commands := make([]Command, 0)

	for p.pos < len(p.tokens) {
		// Skip newlines and comments at the start
		p.skipWhitespaceAndComments()

		if p.pos >= len(p.tokens) {
			break
		}

		cmd, err := p.parseCommand()
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)

		// Skip newlines and comments after command
		p.skipWhitespaceAndComments()

		// Check for pipe
		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "pipe" {
			p.pos++ // consume pipe
			p.skipWhitespaceAndComments()
		}
	}

	if len(commands) == 0 {
		return []Command{{Type: CmdHello}}, nil
	}

	return commands, nil
}

// skipWhitespaceAndComments skips newlines and comments
func (p *Parser) skipWhitespaceAndComments() {
	for p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "newline" || p.tokens[p.pos].Type == "comment") {
		p.pos++
	}
}

// parseCommand parses a single command
func (p *Parser) parseCommand() (Command, error) {
	if p.pos >= len(p.tokens) {
		return Command{}, errors.New("unexpected end of input")
	}

	token := p.tokens[p.pos]

	// Handle comment as a command
	if token.Type == "comment" {
		p.pos++
		return Command{Type: CmdComment, Value: token.Value}, nil
	}

	if token.Type != "identifier" {
		return Command{}, fmt.Errorf("expected command name, got %s", token.Type)
	}

	cmdName := token.Value
	p.pos++

	switch cmdName {
	case "hello":
		return Command{Type: CmdHello}, nil
	case "ping":
		return Command{Type: CmdPing, Value: "pong"}, nil
	case "echo":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		value, err := p.parseString()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdEcho, Value: value}, nil
	case "count":
		return Command{Type: CmdCount}, nil
	case "limit":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		num, err := p.parseNumber()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdLimit, Value: num}, nil
	case "order":
		// Check if next token is "by"
		if p.pos >= len(p.tokens) || p.tokens[p.pos].Value != "by" {
			return Command{}, errors.New("expected 'by' after 'order'")
		}
		p.pos++
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		orderArgs, err := p.parseOrderByArgs()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdOrderBy, Value: orderArgs}, nil
	case "project":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		projections, err := p.parseProjections()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdProject, Value: projections}, nil
	case "project-away":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		fields, err := p.parseFieldList()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdProjectAway, Value: fields}, nil
	case "project-reorder":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		fields, err := p.parseFieldList()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdProjectReorder, Value: fields}, nil
	case "scope":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		field, err := p.parseString()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdScope, Value: RefType(field)}, nil
	case "where":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		conditions, err := p.parseWhereConditions()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdWhere, Value: conditions}, nil
	case "distinct":
		// Check if there's a field name
		if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "identifier") {
			if err := p.skipWhitespace(); err != nil {
				return Command{}, err
			}
			field, err := p.parseStringOrIdentifier()
			if err != nil {
				return Command{}, err
			}
			return Command{Type: CmdDistinct, Value: RefType(field)}, nil
		}
		return Command{Type: CmdDistinct, Value: nil}, nil
	case "parse-json":
		args, err := p.parseParseArgs()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdParseJSON, Value: args}, nil
	case "parse-csv":
		args, err := p.parseParseArgs()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdParseCSV, Value: args}, nil
	case "parse-xml":
		args, err := p.parseParseArgs()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdParseXML, Value: args}, nil
	case "parse-yaml":
		args, err := p.parseParseArgs()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdParseYAML, Value: args}, nil
	case "extend":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		extensions, err := p.parseExtensions()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdExtend, Value: extensions}, nil
	case "summarize":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		summarizeItem, err := p.parseSummarize()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdSummarize, Value: summarizeItem}, nil
	case "pivot":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		pivotItem, err := p.parsePivot()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdPivot, Value: pivotItem}, nil
	default:
		return Command{}, fmt.Errorf("unknown command: %s", cmdName)
	}
}

// skipWhitespace skips whitespace tokens
func (p *Parser) skipWhitespace() error {
	// In our token stream, whitespace is already handled
	return nil
}

// parseString parses a string value
func (p *Parser) parseString() (string, error) {
	if p.pos >= len(p.tokens) {
		return "", errors.New("expected string")
	}
	token := p.tokens[p.pos]
	if token.Type != "string" && token.Type != "sq_string" {
		return "", fmt.Errorf("expected string, got %s", token.Type)
	}
	p.pos++
	return token.Value, nil
}

// parseStringOrIdentifier parses a string or identifier
func (p *Parser) parseStringOrIdentifier() (string, error) {
	if p.pos >= len(p.tokens) {
		return "", errors.New("expected string or identifier")
	}
	token := p.tokens[p.pos]
	if token.Type != "string" && token.Type != "sq_string" && token.Type != "identifier" {
		return "", fmt.Errorf("expected string or identifier, got %s", token.Type)
	}
	p.pos++
	return token.Value, nil
}

// parseNumber parses a number value
func (p *Parser) parseNumber() (float64, error) {
	if p.pos >= len(p.tokens) {
		return 0, errors.New("expected number")
	}
	token := p.tokens[p.pos]
	if token.Type != "number" {
		return 0, fmt.Errorf("expected number, got %s", token.Type)
	}
	p.pos++
	num, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", token.Value)
	}
	return num, nil
}

// parseOrderByArgs parses order by arguments
func (p *Parser) parseOrderByArgs() ([]OrderByArg, error) {
	args := make([]OrderByArg, 0)

	for {
		field, err := p.parseString()
		if err != nil {
			return nil, err
		}

		direction := "asc" // default
		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
			dir := p.tokens[p.pos].Value
			if dir == "asc" || dir == "desc" {
				direction = dir
				p.pos++
			}
		}

		args = append(args, OrderByArg{Field: field, Direction: direction})

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "comma" {
			break
		}
		p.pos++ // consume comma
	}

	return args, nil
}

// parseProjections parses project arguments
func (p *Parser) parseProjections() ([]interface{}, error) {
	projections := make([]interface{}, 0)

	for {
		// Check for alias assignment
		var alias string
		startPos := p.pos

		// Try to parse string
		if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
			alias = p.tokens[p.pos].Value
			p.pos++

			// Check for assignment
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "assignment" {
				p.pos++ // consume =
				// Parse the value (could be function or reference)
				if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
					// Check if it's a function call
					if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "lparen" {
						// It's a function
						fn, err := p.parseFunction()
						if err != nil {
							return nil, err
						}
						fn.Alias = alias
						projections = append(projections, fn)
					} else {
						// It's a reference
						field, err := p.parseStringOrIdentifier()
						if err != nil {
							return nil, err
						}
						projections = append(projections, RefType(field, alias))
					}
				} else if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
					field, err := p.parseString()
					if err != nil {
						return nil, err
					}
					projections = append(projections, RefType(field, alias))
				} else {
					return nil, errors.New("expected function or field reference after assignment")
				}
			} else {
				// No assignment, just a field reference
				projections = append(projections, RefType(alias))
			}
		} else if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
			// Could be a function or identifier
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "lparen" {
				fn, err := p.parseFunction()
				if err != nil {
					return nil, err
				}
				projections = append(projections, fn)
			} else {
				field := p.tokens[p.pos].Value
				p.pos++
				projections = append(projections, RefType(field))
			}
		} else {
			p.pos = startPos
			break
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "comma" {
			break
		}
		p.pos++ // consume comma
	}

	return projections, nil
}

// parseFieldList parses a comma-separated list of fields
func (p *Parser) parseFieldList() ([]TypedValue, error) {
	fields := make([]TypedValue, 0)

	for {
		field, err := p.parseString()
		if err != nil {
			return nil, err
		}
		fields = append(fields, RefType(field))

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "comma" {
			break
		}
		p.pos++ // consume comma
	}

	return fields, nil
}

// parseWhereConditions parses where clause conditions
func (p *Parser) parseWhereConditions() ([]interface{}, error) {
	// Simplified implementation - just parse the expression as a string for now
	conditions := make([]interface{}, 0)
	// This is a placeholder - full implementation would parse the actual condition logic
	return conditions, nil
}

// parseParseArgs parses parse command arguments
func (p *Parser) parseParseArgs() ([][]ParseArg, error) {
	// For now, return empty args - full implementation would parse actual arguments
	return [][]ParseArg{}, nil
}

// parseExtensions parses extend command extensions
func (p *Parser) parseExtensions() ([]interface{}, error) {
	extensions := make([]interface{}, 0)

	for {
		// Similar to projections
		var alias string

		if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
			alias = p.tokens[p.pos].Value
			p.pos++

			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "assignment" {
				p.pos++ // consume =
				if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
					// Check if it's a function call
					if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "lparen" {
						fn, err := p.parseFunction()
						if err != nil {
							return nil, err
						}
						fn.Alias = alias
						extensions = append(extensions, fn)
					} else {
						field, err := p.parseStringOrIdentifier()
						if err != nil {
							return nil, err
						}
						extensions = append(extensions, RefType(field, alias))
					}
				} else {
					return nil, errors.New("expected function or field reference after assignment")
				}
			}
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "comma" {
			break
		}
		p.pos++ // consume comma
	}

	return extensions, nil
}

// parseFunction parses a function call
func (p *Parser) parseFunction() (FunctionCall, error) {
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" {
		return FunctionCall{}, errors.New("expected function name")
	}

	fnName := p.tokens[p.pos].Value
	p.pos++

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "lparen" {
		return FunctionCall{}, errors.New("expected '(' after function name")
	}
	p.pos++ // consume (

	args := make([]TypedValue, 0)

	// Parse arguments
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != "rparen" {
		var arg TypedValue

		token := p.tokens[p.pos]
		switch token.Type {
		case "string", "sq_string":
			arg = StringType(token.Value)
			p.pos++
		case "number":
			num, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return FunctionCall{}, err
			}
			arg = NumberType(num)
			p.pos++
		case "identifier":
			arg = RefType(token.Value)
			p.pos++
		default:
			return FunctionCall{}, fmt.Errorf("unexpected token type in function arguments: %s", token.Type)
		}

		args = append(args, arg)

		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
			p.pos++ // consume comma
		}
	}

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "rparen" {
		return FunctionCall{}, errors.New("expected ')' after function arguments")
	}
	p.pos++ // consume )

	return FunctionCall{
		Type:     "function",
		Operator: FunctionName(fnName),
		Args:     args,
	}, nil
}

// parseSummarize parses a summarize command
// Format: summarize <assignments> [by <fields>]
// Example: summarize "total"=sum("qty") by "category"
func (p *Parser) parseSummarize() (SummarizeItem, error) {
	// Parse summarize assignments
	assignments, err := p.parseSummarizeAssignments()
	if err != nil {
		return SummarizeItem{}, err
	}

	result := SummarizeItem{
		Metrics: assignments,
		By:      make([]TypedValue, 0),
	}

	// Check for "by" keyword
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" && p.tokens[p.pos].Value == "by" {
		p.pos++ // consume "by"

		// Parse the "by" fields
		byFields, err := p.parseFieldList()
		if err != nil {
			return SummarizeItem{}, err
		}
		result.By = byFields
	}

	return result, nil
}

// parseSummarizeAssignments parses comma-separated summarize assignments
func (p *Parser) parseSummarizeAssignments() ([]SummarizeAssignment, error) {
	assignments := make([]SummarizeAssignment, 0)

	for {
		var alias string

		// Check for alias (optional)
		if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
			alias = p.tokens[p.pos].Value
			p.pos++

			// Check for assignment operator
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "assignment" {
				p.pos++ // consume =
			} else {
				// No assignment, this might be a field reference, backtrack
				p.pos--
				alias = ""
			}
		}

		// Parse the function (e.g., sum("qty"), count(), etc.)
		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" {
			return nil, errors.New("expected function name in summarize")
		}

		fnName := p.tokens[p.pos].Value
		p.pos++

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "lparen" {
			return nil, errors.New("expected '(' after function name in summarize")
		}
		p.pos++ // consume (

		// Parse arguments
		args := make([]TypedValue, 0)
		for p.pos < len(p.tokens) && p.tokens[p.pos].Type != "rparen" {
			if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
				args = append(args, StringType(p.tokens[p.pos].Value))
				p.pos++
			} else if p.tokens[p.pos].Type == "identifier" {
				args = append(args, RefType(p.tokens[p.pos].Value))
				p.pos++
			} else {
				return nil, fmt.Errorf("unexpected token in summarize function args: %s", p.tokens[p.pos].Type)
			}

			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
				p.pos++ // consume comma
			}
		}

		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "rparen" {
			return nil, errors.New("expected ')' after function arguments in summarize")
		}
		p.pos++ // consume )

		assignment := SummarizeAssignment{
			Alias:    alias,
			Operator: FunctionName(fnName),
			Args:     args,
		}
		assignments = append(assignments, assignment)

		// Check for comma (more assignments)
		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
			p.pos++ // consume comma
			continue
		}

		break
	}

	return assignments, nil
}

// parsePivot parses a pivot command
// Format: pivot <function>, [<row_field>], [<col_field>]
// Example: pivot sum("qty"), "fruit", "size"
func (p *Parser) parsePivot() (PivotItem, error) {
	// Parse the metric assignment (required)
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" {
		return PivotItem{}, errors.New("expected function name in pivot")
	}

	fnName := p.tokens[p.pos].Value
	p.pos++

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "lparen" {
		return PivotItem{}, errors.New("expected '(' after function name in pivot")
	}
	p.pos++ // consume (

	// Parse function arguments
	args := make([]TypedValue, 0)
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type != "rparen" {
		if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
			args = append(args, StringType(p.tokens[p.pos].Value))
			p.pos++
		} else if p.tokens[p.pos].Type == "identifier" {
			args = append(args, RefType(p.tokens[p.pos].Value))
			p.pos++
		} else {
			return PivotItem{}, fmt.Errorf("unexpected token in pivot function args: %s", p.tokens[p.pos].Type)
		}

		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
			p.pos++ // consume comma
		}
	}

	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "rparen" {
		return PivotItem{}, errors.New("expected ')' after function arguments in pivot")
	}
	p.pos++ // consume )

	metric := SummarizeAssignment{
		Operator: FunctionName(fnName),
		Args:     args,
	}

	result := PivotItem{
		Metric: metric,
		Fields: make([]TypedValue, 0),
	}

	// Check for comma and optional fields
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
		p.pos++ // consume comma

		// Parse row and column fields
		fields, err := p.parseFieldList()
		if err != nil {
			return PivotItem{}, err
		}
		result.Fields = fields
	}

	return result, nil
}

// Helper function to match regex patterns
func matchRegex(pattern, text string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(text)
}
