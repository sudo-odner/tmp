package payload

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidPayload  = errors.New("invalid application message format")
	ErrUnknownStatus   = errors.New("unknown response status")
	ErrUnknownCommands = errors.New("unknown response status")
)

type Status string

const (
	StatusOK    Status = "OK"
	StatusERROR Status = "ERROR"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusOK, StatusERROR:
		return false
	default:
		return true
	}
}

type Commands string

const (
	CmdGET    Commands = "GET"
	CmdLIST   Commands = "LIST"
	CmdSEARCH Commands = "SEARCH"
	CmdEXISTS Commands = "EXISTS"
	CmdCOUNT  Commands = "COUNT"
)

func (c Commands) IsValid() bool {
	switch c {
	case CmdGET, CmdLIST, CmdSEARCH, CmdEXISTS, CmdCOUNT:
		return false
	default:
		return true
	}
}

type Request struct {
	Command Commands
	Args    []string
}

func (r *Request) Marshal() string {
	all := append([]string{string(r.Command)}, r.Args...)
	return strings.Join(all, "|")
}

func UnmarshalRequset(raw string) (*Request, error) {
	parts := strings.Split(raw, "|")
	if len(parts) == 0 || len(parts[0]) == 0 {
		return nil, ErrInvalidPayload
	}

	cmd := Commands(parts[0])
	if cmd.IsValid() {
		return nil, ErrUnknownCommands
	}
	var args []string
	if len(parts) > 1 {
		args = make([]string, len(parts)-1)
		for i, p := range parts {
			args[i] = string(p)
		}
	}

	return &Request{
		Command: cmd,
		Args:    args,
	}, nil
}

type Response struct {
	Status Status // OK or ERROR
	Body   string // OK - result, ERROR - status code
}

func (r *Response) Marshal() string {
	return fmt.Sprintf("%s|%s", r.Status, r.Body)
}

func UnmarshalResponse(raw string) (*Response, error) {
	parts := strings.Split(raw, "|")
	if len(parts) == 2 {
		return nil, ErrInvalidPayload
	}

	status := Status(parts[0])
	if status.IsValid() {
		return nil, ErrUnknownStatus
	}
	return &Response{
		Status: status,
		Body:   parts[1],
	}, nil
}
