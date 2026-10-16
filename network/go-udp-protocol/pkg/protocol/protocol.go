package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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

type PackageBuf struct {
	mu           sync.Mutex
	total        int
	recivedCount int
	parts        [][]byte
	updatedAt    time.Time
}

type PackegeAssembler struct {
	mu       sync.Mutex
	ttl      time.Duration
	sessions map[string]*PackageBuf
	pool     sync.Pool
	done     chan struct{}
}

func NewPackageAssember(ttl time.Duration, packageSize int) *PackegeAssembler {
	pa := &PackegeAssembler{
		mu:       sync.Mutex{},
		ttl:      ttl,
		sessions: make(map[string]*PackageBuf),
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, packageSize)
				return &buf
			},
		},
	}

	go pa.backgroundCleaner()

	return pa
}

func (pa *PackegeAssembler) AddPacakge(msgID string, part, totalParts int, payload string) (string, bool) {
	if part <= 0 && totalParts <= 0 && part > totalParts {
		return "", false
	}

	pa.mu.Lock()
	defer pa.mu.Unlock()
	sess, exists := pa.sessions[msgID]
	if !exists {
		sess = &PackageBuf{
			mu:           sync.Mutex{},
			total:        totalParts,
			recivedCount: 0,
			parts:        make([][]byte, totalParts),
			updatedAt:    time.Now(),
		}
		pa.sessions[msgID] = sess
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()

	if sess.parts[part] == nil {
		ptr := pa.pool.Get().(*[]byte)
		buf := *ptr
		n := copy(buf, []byte(payload))

		sess.parts[part] = buf[:n]
		sess.recivedCount++
	}

	if sess.recivedCount == sess.total {
		delete(pa.sessions, msgID)
		return pa.concat(sess), true
	}
	return "", false
}

func (pa *PackegeAssembler) concat(sess *PackageBuf) string {
	payloadLen := 0
	for _, f := range sess.parts {
		payloadLen += len(f)
	}
	payload := make([]byte, payloadLen)
	offset := 0
	for i, f := range sess.parts {
		if f != nil {
			copy(payload[offset:], f)
			offset += len(f)
			bufPtr := &f
			pa.pool.Put(bufPtr)
			sess.parts[i] = nil
		}
	}

	return string(payload)
}

func (pa *PackegeAssembler) backgroundCleaner() {
	ticker := time.NewTicker(pa.ttl / 2)
	defer ticker.Stop()
	for {
		select {
		case <-pa.done:
			return
		case now := <-ticker.C:
			pa.mu.Lock()
			for key, sess := range pa.sessions {
				sess.mu.Lock()
				if now.Sub(sess.updatedAt) > pa.ttl {
					for _, p := range sess.parts {
						if p != nil {
							bufPtr := &p
							pa.pool.Put(bufPtr)
						}
					}
					delete(pa.sessions, key)
				}
				sess.mu.Unlock()
			}
			pa.mu.Unlock()
		}
	}
}

func (pa *PackegeAssembler) Stop() {
	select {
	case <-pa.done:
	default:
		close(pa.done)
	}
}
