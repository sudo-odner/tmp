package assembly

import (
	"bytes"
	"testing"
	"time"
)

func TestAssembler_InOrder(t *testing.T) {
	a := NewAssembler(time.Second)
	defer a.Stop()

	sessionKey := "client1#1"

	res, complete := a.AddChunk(sessionKey, 1, 3, []byte("Part1-"))
	if complete || res != nil {
		t.Fatalf("expected incomplete for chunk 1")
	}

	res, complete = a.AddChunk(sessionKey, 2, 3, []byte("Part2-"))
	if complete || res != nil {
		t.Fatalf("expected incomplete for chunk 2")
	}

	res, complete = a.AddChunk(sessionKey, 3, 3, []byte("Part3"))
	if !complete {
		t.Fatalf("expected complete for chunk 3")
	}

	expected := "Part1-Part2-Part3"
	if string(res) != expected {
		t.Errorf("got %s, want %s", string(res), expected)
	}
}

func TestAssembler_OutOfOrder(t *testing.T) {
	a := NewAssembler(time.Second)
	defer a.Stop()

	sessionKey := "client1#2"

	_, complete := a.AddChunk(sessionKey, 2, 3, []byte("World"))
	if complete {
		t.Fatalf("expected incomplete for chunk 2")
	}

	_, complete = a.AddChunk(sessionKey, 3, 3, []byte("!"))
	if complete {
		t.Fatalf("expected incomplete for chunk 3")
	}

	res, complete := a.AddChunk(sessionKey, 1, 3, []byte("Hello "))
	if !complete {
		t.Fatalf("expected complete when all chunks arrived")
	}

	expected := "Hello World!"
	if string(res) != expected {
		t.Errorf("got %s, want %s", string(res), expected)
	}
}

func TestAssembler_DuplicateAndInvalidChunks(t *testing.T) {
	a := NewAssembler(time.Second)
	defer a.Stop()

	sessionKey := "client1#3"

	// Invalid chunk bounds
	_, complete := a.AddChunk(sessionKey, 0, 2, []byte("invalid"))
	if complete {
		t.Errorf("chunk part <= 0 should be rejected")
	}

	_, complete = a.AddChunk(sessionKey, 3, 2, []byte("invalid"))
	if complete {
		t.Errorf("chunk part > total should be rejected")
	}

	// Valid chunk 1
	a.AddChunk(sessionKey, 1, 2, []byte("Hello"))
	// Duplicate chunk 1
	a.AddChunk(sessionKey, 1, 2, []byte("Hello"))

	res, complete := a.AddChunk(sessionKey, 2, 2, []byte(" World"))
	if !complete {
		t.Fatalf("expected complete")
	}
	if !bytes.Equal(res, []byte("Hello World")) {
		t.Errorf("got %s, want Hello World", string(res))
	}
}

func TestAssembler_TTL_Cleanup(t *testing.T) {
	ttl := 100 * time.Millisecond
	a := NewAssembler(ttl)
	defer a.Stop()

	sessionKey := "client1#4"
	a.AddChunk(sessionKey, 1, 2, []byte("Part1"))

	// Wait for cleaner ticker to fire and clean expired buffer
	time.Sleep(300 * time.Millisecond)

	a.mu.Lock()
	_, exists := a.buffers[sessionKey]
	a.mu.Unlock()

	if exists {
		t.Errorf("expected buffer to be cleaned up after TTL")
	}
}
