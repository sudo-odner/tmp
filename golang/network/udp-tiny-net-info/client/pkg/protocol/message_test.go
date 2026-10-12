package protocol

import (
	"reflect"
	"testing"
)

func TestRequestMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		req     *Request
		wantErr bool
	}{
		{
			name: "single command without args",
			req: &Request{
				Command: "LIST",
				Args:    nil,
			},
			wantErr: false,
		},
		{
			name: "command with one arg",
			req: &Request{
				Command: "GET",
				Args:    []string{"HTTP"},
			},
			wantErr: false,
		},
		{
			name: "command with multiple args",
			req: &Request{
				Command: "SEARCH",
				Args:    []string{"web", "traffic"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.req.Marshal()
			got, err := UnmarshalRequest(data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.Command != tt.req.Command {
					t.Errorf("got Command = %s, want %s", got.Command, tt.req.Command)
				}
				if len(tt.req.Args) == 0 && len(got.Args) != 0 {
					t.Errorf("got Args = %v, want empty", got.Args)
				} else if len(tt.req.Args) > 0 && !reflect.DeepEqual(got.Args, tt.req.Args) {
					t.Errorf("got Args = %v, want %v", got.Args, tt.req.Args)
				}
			}
		})
	}

	t.Run("invalid empty request", func(t *testing.T) {
		_, err := UnmarshalRequest([]byte(""))
		if err == nil {
			t.Errorf("expected error for empty request, got nil")
		}
	})
}

func TestResponseMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		resp    *Response
		wantErr bool
	}{
		{
			name: "status ok",
			resp: &Response{
				Status: StatusOK,
				Body:   "80/TCP Hypertext Transfer Protocol",
			},
			wantErr: false,
		},
		{
			name: "status error",
			resp: &Response{
				Status: StatusError,
				Body:   "NOT_FOUND",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.resp.Marshal()
			got, err := UnmarshalResponse(data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.Status != tt.resp.Status || got.Body != tt.resp.Body {
					t.Errorf("got %+v, want %+v", got, tt.resp)
				}
			}
		})
	}

	t.Run("unknown status", func(t *testing.T) {
		_, err := UnmarshalResponse([]byte("INVALID|test"))
		if err != ErrUnknownStatus {
			t.Errorf("expected ErrUnknownStatus, got %v", err)
		}
	})

	t.Run("malformed response without pipe", func(t *testing.T) {
		_, err := UnmarshalResponse([]byte("OK"))
		if err != ErrInvalidMessage {
			t.Errorf("expected ErrInvalidMessage, got %v", err)
		}
	})
}
