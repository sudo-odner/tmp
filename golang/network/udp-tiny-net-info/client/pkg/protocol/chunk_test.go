package protocol

import (
	"bytes"
	"testing"
)

func TestChunkMarshalUnmarshal(t *testing.T) {
	t.Run("single chunk", func(t *testing.T) {
		ch := Chunk{
			MessageID: "123",
			Part:      1,
			Total:     1,
			Data:      []byte("GET|HTTP"),
		}

		raw := ch.Marshal()
		expected := "123|GET|HTTP"
		if string(raw) != expected {
			t.Errorf("got %s, want %s", string(raw), expected)
		}

		parsed, err := UnmarshalChunk(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if parsed.MessageID != "123" || parsed.Part != 1 || parsed.Total != 1 || !bytes.Equal(parsed.Data, ch.Data) {
			t.Errorf("parsed chunk mismatch: %+v vs %+v", parsed, ch)
		}
	})

	t.Run("multi chunk", func(t *testing.T) {
		ch := Chunk{
			MessageID: "msg-99",
			Part:      2,
			Total:     5,
			Data:      []byte("partial payload chunk data"),
		}

		raw := ch.Marshal()
		expected := "msg-99|2|5|partial payload chunk data"
		if string(raw) != expected {
			t.Errorf("got %s, want %s", string(raw), expected)
		}

		parsed, err := UnmarshalChunk(raw)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if parsed.MessageID != "msg-99" || parsed.Part != 2 || parsed.Total != 5 || !bytes.Equal(parsed.Data, ch.Data) {
			t.Errorf("parsed chunk mismatch: %+v vs %+v", parsed, ch)
		}
	})

	t.Run("empty packet error", func(t *testing.T) {
		_, err := UnmarshalChunk([]byte{})
		if err != ErrEmptyPacket {
			t.Errorf("expected ErrEmptyPacket, got %v", err)
		}
	})

	t.Run("malformed packet without separator", func(t *testing.T) {
		_, err := UnmarshalChunk([]byte("nodividers"))
		if err != ErrMalformedHeader {
			t.Errorf("expected ErrMalformedHeader, got %v", err)
		}
	})
}

func TestSplitIntoChunks(t *testing.T) {
	msgID := "test-msg"
	data := []byte("1234567890abcdefghij") // 20 bytes

	t.Run("payload fits in single chunk", func(t *testing.T) {
		chunks := SplitIntoChunks(msgID, data, 50)
		if len(chunks) != 1 {
			t.Fatalf("expected 1 chunk, got %d", len(chunks))
		}
		if chunks[0].Part != 1 || chunks[0].Total != 1 || !bytes.Equal(chunks[0].Data, data) {
			t.Errorf("unexpected chunk content: %+v", chunks[0])
		}
	})

	t.Run("payload split into multiple chunks", func(t *testing.T) {
		chunkSize := 7
		chunks := SplitIntoChunks(msgID, data, chunkSize)
		expectedChunks := (len(data) + chunkSize - 1) / chunkSize // 3 chunks (7 + 7 + 6)
		if len(chunks) != expectedChunks {
			t.Fatalf("expected %d chunks, got %d", expectedChunks, len(chunks))
		}

		var reassembled []byte
		for i, ch := range chunks {
			if ch.MessageID != msgID {
				t.Errorf("chunk %d has wrong msgID: %s", i, ch.MessageID)
			}
			if ch.Part != i+1 {
				t.Errorf("chunk %d has wrong Part: %d", i, ch.Part)
			}
			if ch.Total != expectedChunks {
				t.Errorf("chunk %d has wrong Total: %d", i, ch.Total)
			}
			reassembled = append(reassembled, ch.Data...)
		}

		if !bytes.Equal(reassembled, data) {
			t.Errorf("reassembled data mismatch: got %s, want %s", string(reassembled), string(data))
		}
	})
}
