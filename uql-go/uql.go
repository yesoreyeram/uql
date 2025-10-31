package uql

import (
	"errors"
	"time"
)

// Options represents options for UQL execution
type Options struct {
	Data interface{}
}

// UQL executes a UQL query with the given options
func UQL(query string, options *Options) (interface{}, error) {
	if query == "" {
		return "hello there! provide a valid query", nil
	}

	// Parse the query into AST
	commands, err := Parse(query)
	if err != nil {
		return nil, err
	}

	// Evaluate the commands
	return Evaluate(commands, options)
}

// Evaluate executes a list of commands
func Evaluate(commands []Command, options *Options) (interface{}, error) {
	var output interface{}
	if options != nil {
		output = options.Data
	}

	currentTime := time.Now()
	result := CommandResult{
		Output: output,
		Context: map[string]interface{}{
			"currentTime": currentTime,
		},
	}

	for _, cmd := range commands {
		var err error
		result, err = evaluateCommand(result, cmd)
		if err != nil {
			return nil, err
		}
	}

	return result.Output, nil
}

// evaluateCommand evaluates a single command
func evaluateCommand(prev CommandResult, cmd Command) (CommandResult, error) {
	switch cmd.Type {
	case CmdComment:
		return evalComment(prev, cmd)
	case CmdHello:
		return evalHello(prev, cmd)
	case CmdPing:
		return evalPing(prev, cmd)
	case CmdEcho:
		return evalEcho(prev, cmd)
	case CmdCount:
		return evalCount(prev, cmd)
	case CmdLimit:
		return evalLimit(prev, cmd)
	case CmdOrderBy:
		return evalOrderBy(prev, cmd)
	case CmdProject:
		return evalProject(prev, cmd)
	case CmdProjectAway:
		return evalProjectAway(prev, cmd)
	case CmdProjectReorder:
		return evalProjectReorder(prev, cmd)
	case CmdScope:
		return evalScope(prev, cmd)
	case CmdWhere:
		return evalWhere(prev, cmd)
	case CmdDistinct:
		return evalDistinct(prev, cmd)
	case CmdMvExpand:
		return evalMvExpand(prev, cmd)
	case CmdExtend:
		return evalExtend(prev, cmd)
	case CmdSummarize:
		return evalSummarize(prev, cmd)
	case CmdPivot:
		return evalPivot(prev, cmd)
	case CmdParseJSON:
		return evalParseJSON(prev, cmd)
	case CmdParseCSV:
		return evalParseCSV(prev, cmd)
	case CmdParseXML:
		return evalParseXML(prev, cmd)
	case CmdParseYAML:
		return evalParseYAML(prev, cmd)
	case CmdCommand:
		return evalCommandFunc(prev, cmd)
	case CmdJSONata:
		return evalJSONata(prev, cmd)
	case CmdRange:
		return evalRange(prev, cmd)
	default:
		return prev, errors.New("unknown command type: " + string(cmd.Type))
	}
}
