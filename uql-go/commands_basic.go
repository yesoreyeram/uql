package uql

// evalComment evaluates a comment command
func evalComment(prev CommandResult, cmd Command) (CommandResult, error) {
	// Comments don't modify the output
	return prev, nil
}

// evalHello evaluates a hello command
func evalHello(prev CommandResult, cmd Command) (CommandResult, error) {
	return CommandResult{
		Output:  "hello",
		Context: prev.Context,
	}, nil
}

// evalPing evaluates a ping command
func evalPing(prev CommandResult, cmd Command) (CommandResult, error) {
	return CommandResult{
		Output:  "pong",
		Context: prev.Context,
	}, nil
}

// evalEcho evaluates an echo command
func evalEcho(prev CommandResult, cmd Command) (CommandResult, error) {
	value, ok := cmd.Value.(string)
	if !ok {
		return prev, nil
	}
	return CommandResult{
		Output:  value,
		Context: prev.Context,
	}, nil
}
