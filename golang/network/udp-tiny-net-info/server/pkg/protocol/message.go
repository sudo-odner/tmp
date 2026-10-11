package protocol

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrInvalidMessage = errors.New("invalied application message format")
	ErrUnknownStatus  = errors.New("unknown response status")
)

// Status status response server
type Status string

const (
	StatusOK    Status = "OK"
	StatusError Status = "ERROR"
)

// Request - application request message
// Format: CMD|ARG1|ARG2...
type Request struct {
	Command string
	Args    []string
}

func (r *Request) Marshal() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.Command)

	for _, arg := range r.Args {
		buf.WriteByte('|')
		buf.WriteString(arg)
	}
	return buf.Bytes()
}

func UnmarshalRequest(data []byte) (*Request, error) {
	parts := bytes.Split(data, []byte("|"))
	if len(parts) == 0 || len(parts[0]) == 0 {
		return nil, ErrInvalidMessage
	}

	cmd := string(parts[0])
	var args []string
	if len(parts) > 1 {
		args = make([]string, len(parts)-1)
		for i, p := range parts[1:] {
			args[i] = string(p)
		}
	}

	return &Request{
		Command: cmd,
		Args:    args,
	}, nil
}

// Response - application response message
// Format: STATUS|BODY
type Response struct {
	Status Status
	Body   string // OK - result, ERROR - status code
}

func (r *Response) Marshal() []byte {
	result := fmt.Sprintf("%s|%s", r.Status, r.Body)
	return []byte(result)
}

func UnmarshalResponse(data []byte) (*Response, error) {
	parts := bytes.SplitN(data, []byte("|"), 2)
	if len(parts) < 2 {
		return nil, ErrInvalidMessage
	}

	st := Status(parts[0])
	if st != StatusOK && st != StatusError {
		return nil, ErrUnknownStatus
	}
	return &Response{
		Status: st,
		Body:   string(parts[1]),
	}, nil
}
