package payload

import (
	"bytes"
	"errors"
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
		return true
	default:
		return false
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
		return true
	default:
		return false
	}
}

type Request struct {
	Command Commands
	Args    []string
}

func (r *Request) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteString(string(r.Command))
	buf.WriteByte('|')
	for _, arg := range r.Args {
		buf.WriteString(arg)
	}
	return buf.Bytes()
}

func UnmarshalRequset(raw []byte) (*Request, error) {
	parts := bytes.Split(raw, []byte{'|'})
	if len(parts) == 0 || len(parts[0]) == 0 {
		return nil, ErrInvalidPayload
	}

	cmd := Commands(string(parts[0]))
	if !cmd.IsValid() {
		return nil, ErrUnknownCommands
	}
	var args []string
	if len(parts) > 1 {
		args = make([]string, len(parts)-1)
		for i := 1; i < len(parts); i++ {
			args[i-1] = string(parts[i])
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

func (r *Response) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteString(string(r.Status))
	buf.WriteByte('|')
	buf.WriteString(r.Body)
	return buf.Bytes()
}

func UnmarshalResponse(raw []byte) (*Response, error) {
	parts := bytes.Split(raw, []byte{'|'})
	if len(parts) != 2 {
		return nil, ErrInvalidPayload
	}

	status := Status(string(parts[0]))
	if !status.IsValid() {
		return nil, ErrUnknownStatus
	}
	return &Response{
		Status: status,
		Body:   string(parts[1]),
	}, nil
}
