package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrEmptyPacket     = errors.New("packet is empty")
	ErrMalformedHeader = errors.New("malformed envelope header")
)

type Chunk struct {
	MessageID string
	Part      int
	Total     int
	Data      []byte
}

func (e *Chunk) Marshal() []byte {
	if e.Total > 1 {
		hdr := fmt.Sprintf("%s|%d|%d|", e.MessageID, e.Part, e.Total)
		res := make([]byte, 0, len(hdr)+len(e.Data))
		res = append(res, hdr...)
		res = append(res, e.Data...)
		return res
	}

	hdr := fmt.Sprintf("%s|", e.MessageID)
	res := make([]byte, 0, len(hdr)+len(e.Data))
	res = append(res, hdr...)
	res = append(res, e.Data...)
	return res
}

func UnmarshalChunk(raw []byte) (*Chunk, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyPacket
	}
	parts := bytes.SplitN(raw, []byte("|"), 4)

	if len(parts) == 4 {
		part, errPart := strconv.Atoi(string(parts[1]))
		total, errTotal := strconv.Atoi(string(parts[2]))
		if errPart == nil && errTotal == nil && total > 0 && part > 0 {
			return &Chunk{
				MessageID: string(parts[0]),
				Part:      part,
				Total:     total,
				Data:      parts[3],
			}, nil
		}
	}

	idx := bytes.IndexByte(raw, '|')
	if idx == -1 {
		return nil, ErrMalformedHeader
	}

	return &Chunk{
		MessageID: string(raw[:idx]),
		Part:      1,
		Total:     1,
		Data:      raw[idx+1:],
	}, nil
}

func SplitIntoChunks(msgID string, fullPayload []byte, chunkSize int) []Chunk {
	if chunkSize <= 0 {
		chunkSize = 512
	}

	totalLen := len(fullPayload)

	if totalLen <= chunkSize {
		return []Chunk{
			{
				MessageID: msgID,
				Part:      1,
				Total:     1,
				Data:      fullPayload,
			},
		}
	}

	totalChunks := (totalLen + chunkSize - 1) / chunkSize
	chunks := make([]Chunk, 0, totalChunks)

	for i := 0; i < totalLen; i += chunkSize {
		end := i + chunkSize
		if end > totalLen {
			end = totalLen
		}

		partNum := len(chunks) + 1
		chunks = append(chunks, Chunk{
			MessageID: msgID,
			Part:      partNum,
			Total:     totalChunks,
			Data:      fullPayload[i:end],
		})
	}

	return chunks
}
