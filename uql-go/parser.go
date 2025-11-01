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
			// Check for -- (double dash for parse options)
			if l.pos+1 < len(l.input) && l.input[l.pos+1] == '-' {
				l.tokens = append(l.tokens, Token{Type: "double_dash", Value: "--", Pos: l.pos})
				l.pos += 2 // Skip both dashes
				continue
			}
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
	case "jsonata":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		expression, err := p.parseString()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdJSONata, Value: expression}, nil
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
	case "range":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		rangeItem, err := p.parseRange()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdRange, Value: rangeItem}, nil
	case "mv-expand":
		if err := p.skipWhitespace(); err != nil {
			return Command{}, err
		}
		mvExpandItem, err := p.parseMvExpand()
		if err != nil {
			return Command{}, err
		}
		return Command{Type: CmdMvExpand, Value: mvExpandItem}, nil
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
func (p *Parser) parseWhereConditions() ([]TypedValue, error) {
	// Parse where expression: field operator value
	// Example: "a" == 10 or "a" in (10, 20)
	conditions := make([]TypedValue, 0, 3)

	// Parse left-hand side (field or value)
	lhs, err := p.parseWhereArgument()
	if err != nil {
		return nil, err
	}
	conditions = append(conditions, lhs)

	// Parse operator
	if p.pos >= len(p.tokens) {
		return nil, errors.New("expected operator in where clause")
	}

	operator, err := p.parseWhereOperator()
	if err != nil {
		return nil, err
	}
	conditions = append(conditions, operator)

	// Parse right-hand side (value, field, or array)
	rhs, err := p.parseWhereArgument()
	if err != nil {
		return nil, err
	}
	conditions = append(conditions, rhs)

	return conditions, nil
}

// parseWhereOperator parses comparison operators
func (p *Parser) parseWhereOperator() (TypedValue, error) {
	if p.pos >= len(p.tokens) {
		return TypedValue{}, errors.New("expected operator")
	}

	token := p.tokens[p.pos]
	var op string

	// Handle multi-token operators
	switch token.Type {
	case "eq": // ==
		op = "=="
		p.pos++
	case "ne": // !=
		op = "!="
		p.pos++
	case "gte": // >=
		op = ">="
		p.pos++
	case "lte": // <=
		op = "<="
		p.pos++
	case "gt": // >
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "assignment" {
			op = ">="
			p.pos += 2
		} else {
			op = ">"
			p.pos++
		}
	case "lt": // <
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "assignment" {
			op = "<="
			p.pos += 2
		} else {
			op = "<"
			p.pos++
		}
	case "regex_not_match": // !~
		op = "!~"
		p.pos++
	case "regex_match": // =~
		op = "=~"
		p.pos++
	case "assignment": // = (could be part of =~, ==, etc)
		if p.pos+1 < len(p.tokens) {
			next := p.tokens[p.pos+1]
			if next.Type == "tilde" {
				op = "=~"
				p.pos += 2
			} else if next.Type == "eq" {
				op = "=="
				p.pos += 2
			} else {
				return TypedValue{}, fmt.Errorf("unexpected operator: %s", token.Type)
			}
		} else {
			return TypedValue{}, fmt.Errorf("unexpected operator: %s", token.Type)
		}
	case "exclaim": // ! (could be !=, !~, !contains, !in, etc)
		if p.pos+1 < len(p.tokens) {
			next := p.tokens[p.pos+1]
			if next.Type == "eq" {
				op = "!="
				p.pos += 2
			} else if next.Type == "tilde" {
				op = "!~"
				p.pos += 2
			} else if next.Type == "identifier" {
				switch next.Value {
				case "contains":
					op = "!contains"
					p.pos += 2
				case "contains_cs":
					op = "!contains_cs"
					p.pos += 2
				case "startswith":
					op = "!startswith"
					p.pos += 2
				case "startswith_cs":
					op = "!startswith_cs"
					p.pos += 2
				case "endswith":
					op = "!endswith"
					p.pos += 2
				case "endswith_cs":
					op = "!endswith_cs"
					p.pos += 2
				case "in":
					op = "!in"
					p.pos += 2
				case "in~":
					op = "!in~"
					p.pos += 2
				default:
					return TypedValue{}, fmt.Errorf("unexpected operator: !%s", next.Value)
				}
			} else {
				return TypedValue{}, fmt.Errorf("unexpected operator after !")
			}
		} else {
			return TypedValue{}, errors.New("unexpected end after !")
		}
	case "identifier":
		// Handle word operators like "in", "between", "contains", etc.
		switch token.Value {
		case "in", "between", "inside", "outside", "in~",
			"contains", "contains_cs", "startswith", "startswith_cs",
			"endswith", "endswith_cs", "matches", "not":
			if token.Value == "not" && p.pos+1 < len(p.tokens) {
				next := p.tokens[p.pos+1]
				if next.Type == "identifier" && next.Value == "contains" {
					op = "!contains"
					p.pos += 2
				} else if next.Type == "identifier" && next.Value == "contains_cs" {
					op = "!contains_cs"
					p.pos += 2
				} else {
					return TypedValue{}, fmt.Errorf("unexpected 'not' operator")
				}
			} else if token.Value == "matches" && p.pos+1 < len(p.tokens) {
				next := p.tokens[p.pos+1]
				if next.Type == "identifier" && next.Value == "regex" {
					op = "matches regex"
					p.pos += 2
				} else {
					return TypedValue{}, errors.New("expected 'regex' after 'matches'")
				}
			} else {
				op = token.Value
				p.pos++
			}
		default:
			return TypedValue{}, fmt.Errorf("unexpected identifier in operator position: %s", token.Value)
		}
	default:
		return TypedValue{}, fmt.Errorf("unexpected token type for operator: %s", token.Type)
	}

	return TypedValue{Type: "operation", Value: op}, nil
}

// parseWhereArgument parses a single argument in where clause (field, value, or array)
func (p *Parser) parseWhereArgument() (TypedValue, error) {
	if p.pos >= len(p.tokens) {
		return TypedValue{}, errors.New("unexpected end of expression")
	}

	token := p.tokens[p.pos]

	// Check for array (for 'in', 'between', etc.)
	if token.Type == "lparen" {
		p.pos++ // consume (
		values := make([]interface{}, 0)

		for {
			if p.pos >= len(p.tokens) {
				return TypedValue{}, errors.New("unexpected end in array")
			}

			if p.tokens[p.pos].Type == "rparen" {
				p.pos++ // consume )
				break
			}

			// Parse value
			val, err := p.parseWhereValue()
			if err != nil {
				return TypedValue{}, err
			}
			values = append(values, val)

			// Check for comma
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "comma" {
				p.pos++ // consume ,
			} else if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "rparen" {
				// Will be consumed in next iteration
			} else {
				break
			}
		}

		return TypedValue{Type: "value_array", Value: values}, nil
	}

	// Parse single value or reference
	return p.parseWhereValue()
}

// parseWhereValue parses a single value (number, string, or field reference)
func (p *Parser) parseWhereValue() (TypedValue, error) {
	if p.pos >= len(p.tokens) {
		return TypedValue{}, errors.New("unexpected end of expression")
	}

	token := p.tokens[p.pos]

	switch token.Type {
	case "string":
		// Double-quoted string is a field reference
		p.pos++
		return TypedValue{Type: "ref", Value: token.Value}, nil
	case "sq_string":
		// Single-quoted string is a literal value
		p.pos++
		return TypedValue{Type: "string", Value: token.Value}, nil
	case "number":
		p.pos++
		num, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return TypedValue{}, err
		}
		return TypedValue{Type: "number", Value: num}, nil
	case "identifier":
		// Field reference without quotes
		p.pos++
		return TypedValue{Type: "ref", Value: token.Value}, nil
	default:
		return TypedValue{}, fmt.Errorf("unexpected token type in where argument: %s", token.Type)
	}
}

// parseParseArgs parses parse command arguments
func (p *Parser) parseParseArgs() ([][]ParseArg, error) {
	args := make([]ParseArg, 0)

	// Parse multiple --key value pairs
	for p.pos < len(p.tokens) {
		// Skip whitespace
		if p.tokens[p.pos].Type == "newline" {
			p.pos++
			continue
		}

		// Stop if we hit a pipe or end
		if p.tokens[p.pos].Type == "pipe" {
			break
		}

		// Look for -- (double dash)
		if p.tokens[p.pos].Type != "double_dash" {
			break
		}
		p.pos++ // consume --

		// Get the identifier/key
		if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" {
			return nil, errors.New("expected identifier after --")
		}
		identifier := p.tokens[p.pos].Value
		p.pos++

		// Skip whitespace
		for p.pos < len(p.tokens) && p.tokens[p.pos].Type == "newline" {
			p.pos++
		}

		// Get the value (string or identifier)
		if p.pos >= len(p.tokens) {
			return nil, fmt.Errorf("expected value for --%s", identifier)
		}

		var value string
		switch p.tokens[p.pos].Type {
		case "string", "sq_string":
			value = p.tokens[p.pos].Value
			p.pos++
		case "identifier":
			value = p.tokens[p.pos].Value
			p.pos++
		case "number":
			value = p.tokens[p.pos].Value
			p.pos++
		default:
			return nil, fmt.Errorf("expected value for --%s, got %s", identifier, p.tokens[p.pos].Type)
		}

		args = append(args, ParseArg{
			Identifier: identifier,
			Value:      value,
		})

		// Skip whitespace after value
		for p.pos < len(p.tokens) && p.tokens[p.pos].Type == "newline" {
			p.pos++
		}
	}

	// Return args wrapped in a slice to match the expected type
	if len(args) == 0 {
		return [][]ParseArg{}, nil
	}
	return [][]ParseArg{args}, nil
}

// parseExtensions parses extend command extensions
func (p *Parser) parseExtensions() ([]interface{}, error) {
	extensions := make([]interface{}, 0)

	for {
		// Similar to projections
		var alias string
		startPos := p.pos

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
			} else {
				// No assignment, just a field reference
				extensions = append(extensions, RefType(alias))
			}
		} else if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
			// Could be a function or identifier without alias
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "lparen" {
				fn, err := p.parseFunction()
				if err != nil {
					return nil, err
				}
				extensions = append(extensions, fn)
			} else {
				field := p.tokens[p.pos].Value
				p.pos++
				extensions = append(extensions, RefType(field))
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
		case "string":
			// Double-quoted string is a field reference
			arg = RefType(token.Value)
			p.pos++
		case "sq_string":
			// Single-quoted string is a literal value
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
// Example with alias: pivot "total"=sum("qty"), "fruit", "size"
func (p *Parser) parsePivot() (PivotItem, error) {
	// Check for alias assignment
	var alias string
	if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
		// Might be an alias
		if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == "assignment" {
			alias = p.tokens[p.pos].Value
			p.pos += 2 // consume alias and =
		}
	}

	// Parse the metric function (required)
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
		Alias:    alias,
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

// parseRange parses a range command
// Format: range from <start> to <end> [step <step>]
// Example: range from 1 to 10 step 2
func (p *Parser) parseRange() (RangeValue, error) {
	// Expect "from" keyword
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" || p.tokens[p.pos].Value != "from" {
		return RangeValue{}, errors.New("expected 'from' keyword in range command")
	}
	p.pos++ // consume "from"

	// Parse start value (number or string)
	if p.pos >= len(p.tokens) {
		return RangeValue{}, errors.New("expected start value in range command")
	}

	var start interface{}
	if p.tokens[p.pos].Type == "number" {
		num, err := p.parseNumber()
		if err != nil {
			return RangeValue{}, err
		}
		start = num
	} else if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
		str, err := p.parseString()
		if err != nil {
			return RangeValue{}, err
		}
		start = str
	} else {
		return RangeValue{}, fmt.Errorf("expected number or string for start value, got %s", p.tokens[p.pos].Type)
	}

	// Expect "to" keyword
	if p.pos >= len(p.tokens) || p.tokens[p.pos].Type != "identifier" || p.tokens[p.pos].Value != "to" {
		return RangeValue{}, errors.New("expected 'to' keyword in range command")
	}
	p.pos++ // consume "to"

	// Parse end value (number or string)
	if p.pos >= len(p.tokens) {
		return RangeValue{}, errors.New("expected end value in range command")
	}

	var end interface{}
	if p.tokens[p.pos].Type == "number" {
		num, err := p.parseNumber()
		if err != nil {
			return RangeValue{}, err
		}
		end = num
	} else if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
		str, err := p.parseString()
		if err != nil {
			return RangeValue{}, err
		}
		end = str
	} else {
		return RangeValue{}, fmt.Errorf("expected number or string for end value, got %s", p.tokens[p.pos].Type)
	}

	// Default step values
	var step interface{}
	if _, ok := start.(float64); ok {
		step = 1.0
	} else {
		step = ""
	}

	// Check for optional "step" keyword
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" && p.tokens[p.pos].Value == "step" {
		p.pos++ // consume "step"

		if p.pos >= len(p.tokens) {
			return RangeValue{}, errors.New("expected step value in range command")
		}

		if p.tokens[p.pos].Type == "number" {
			num, err := p.parseNumber()
			if err != nil {
				return RangeValue{}, err
			}
			step = num
		} else if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
			str, err := p.parseString()
			if err != nil {
				return RangeValue{}, err
			}
			step = str
		} else {
			return RangeValue{}, fmt.Errorf("expected number or string for step value, got %s", p.tokens[p.pos].Type)
		}
	}

	return RangeValue{Start: start, End: end, Step: step}, nil
}

// parseMvExpand parses an mv-expand command
// Format: mv-expand "<field>" or mv-expand "<alias>"="<field>"
// Example: mv-expand "users" or mv-expand "user"="users"
func (p *Parser) parseMvExpand() (MvExpandValue, error) {
	var alias string
	var field string

	// Check for alias assignment
	if p.pos < len(p.tokens) && (p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string") {
		firstStr := p.tokens[p.pos].Value
		p.pos++

		// Check if there's an assignment
		if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "assignment" {
			p.pos++ // consume =
			alias = firstStr

			// Parse the actual field name
			if p.pos >= len(p.tokens) {
				return MvExpandValue{}, errors.New("expected field name after '=' in mv-expand")
			}

			if p.tokens[p.pos].Type == "string" || p.tokens[p.pos].Type == "sq_string" {
				field = p.tokens[p.pos].Value
				p.pos++
			} else if p.tokens[p.pos].Type == "identifier" {
				field = p.tokens[p.pos].Value
				p.pos++
			} else {
				return MvExpandValue{}, fmt.Errorf("expected field name in mv-expand, got %s", p.tokens[p.pos].Type)
			}
		} else {
			// No assignment, just the field name
			field = firstStr
		}
	} else if p.pos < len(p.tokens) && p.tokens[p.pos].Type == "identifier" {
		field = p.tokens[p.pos].Value
		p.pos++
	} else {
		return MvExpandValue{}, errors.New("expected field name in mv-expand command")
	}

	return MvExpandValue{Field: field, Alias: alias}, nil
}

// Helper function to match regex patterns
func matchRegex(pattern, text string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(text)
}
