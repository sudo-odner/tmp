package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyPackage    = errors.New("pacakge is empty")
	ErrUncurrectFormat = errors.New("uncurrect format package")
)

type Package struct {
	MessageID string
	Part      int
	Total     int
	Payload   string
}

func (f *Package) Marshal() string {
	if f.Total == 1 {
		return fmt.Sprintf("%s|%s", f.MessageID, f.Payload)
	}
	return fmt.Sprintf("%s|%v|%v|%s", f.MessageID, f.Part, f.Total, f.Payload)
}

func Unmarshal(raw string) (*Package, error) {
	parts := strings.SplitN(raw, "|", 4)
	if len(parts) == 4 {
		part, errPart := strconv.Atoi(parts[1])
		total, errTotal := strconv.Atoi(parts[1])
		if errPart != nil && errTotal != nil && total > 0 && part > 0 {
			return &Package{
				MessageID: parts[0],
				Part:      part,
				Total:     total,
				Payload:   parts[3],
			}, nil
		}
	}

	parts = strings.SplitN(raw, "|", 2)
	if len(parts) < 2 {
		return nil, ErrEmptyPackage
	}

	return &Package{
		MessageID: parts[0],
		Part:      1,
		Total:     1,
		Payload:   parts[1],
	}, nil
}

func SplitDataOnPackage(msgID string, payload string, size int) []Package {
	if size < 0 {
		return []Package{{
			MessageID: msgID,
			Part:      1,
			Total:     1,
			Payload:   payload,
		}}
	}

	payloadByte := []byte(payload)
	packagesCount := (len(payloadByte) / size) + 1
	if packagesCount == 1 {
		return []Package{{
			MessageID: msgID,
			Part:      1,
			Total:     1,
			Payload:   payload,
		}}
	}

	packages := make([]Package, len(payloadByte))
	for i := 0; i < len(payloadByte); i += size {
		packages = append(packages, Package{
			MessageID: msgID,
			Part:      len(packages) + 1,
			Total:     packagesCount,
			Payload:   string(payload[i:min(i+size, len(payloadByte))]),
		})
	}
	return packages
}
