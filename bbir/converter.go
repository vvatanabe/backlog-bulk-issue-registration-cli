package bbir

import (
	"github.com/pkg/errors"
)

func NewCommandConverter(builder CommandBuilder, message Messages) *CommandConverter {
	return &CommandConverter{builder, message}
}

type CommandConverter struct {
	builder  CommandBuilder
	messages Messages
}

func (c *CommandConverter) Convert(lines []*Line, callbacks ...Callback) ([]*Command, error) {

	opts := NewDefaultCallbackOptions()
	for _, cb := range callbacks {
		cb(opts)
	}

	var (
		commands     []*Command
		multiErrsMap = make(map[string]error)
	)

	opts.Before()
	for i, line := range lines {

		lineNum := i + 1

		command, err := c.builder.Build(lineNum, line)
		if err != nil {
			multiErrsMap[c.messages.Line(lineNum)] = err
			continue
		}

		if command.HasUnregisteredParentIssue {

			prevCmdIndex := len(commands) - 1
			if prevCmdIndex < 0 {
				err := errors.New(c.messages.ParentIssueIsNotRegistered("*"))
				multiErrsMap[c.messages.Line(lineNum)] = NewMultipleErrors([]error{err})
				continue
			}

			prevCmd := commands[prevCmdIndex]
			if prevCmd.ParentIssueID != nil {
				err := errors.New(c.messages.ParentIssueAlreadyRegisteredAsChildIssue(*prevCmd.ParentIssueID))
				multiErrsMap[c.messages.Line(lineNum)] = NewMultipleErrors([]error{err})
				continue
			}

			prevCmd.AddChild(command)

		} else {
			commands = append(commands, command)
		}
		opts.Each()
	}
	opts.After()

	if len(multiErrsMap) > 0 {
		return commands, NewMultipleErrorsMap(multiErrsMap)
	}
	return commands, nil
}
