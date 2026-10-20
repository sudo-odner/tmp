package protocol

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	ErrEmptyPackage    = errors.New("pacakge is empty")
	ErrUncurrectFormat = errors.New("uncurrect format package")
)

type Packet struct {
	MessageID string
	Part      int
	Total     int
	Payload   []byte
}

func (f *Packet) Marshal() []byte {
	var buf bytes.Buffer

	if f.Total == 1 {
		buf.WriteString(f.MessageID)
		buf.WriteByte('|')
		buf.Write(f.Payload)
		return buf.Bytes()
	}

	buf.WriteString(f.MessageID)
	buf.WriteByte('|')
	buf.WriteString(strconv.Itoa(f.Part))
	buf.WriteByte('|')
	buf.WriteString(strconv.Itoa(f.Total))
	buf.WriteByte('|')
	buf.Write(f.Payload)
	return buf.Bytes()
}

// Slice packet
// | MsgID | Part | Total | Payload |
//
// Single packet
// | MsgID | Payload |

func Unmarshal(raw []byte) (*Packet, error) {
	parts := bytes.SplitN(raw, []byte{'|'}, 4)
	if len(parts) == 4 {
		part, errPart := strconv.Atoi(string(parts[1]))
		total, errTotal := strconv.Atoi(string(parts[2]))
		if errPart == nil && errTotal == nil && total > 0 && part > 0 && part <= total {
			return &Packet{
				MessageID: string(parts[0]),
				Part:      part,
				Total:     total,
				Payload:   parts[3],
			}, nil
		}
	}

	parts = bytes.SplitN(raw, []byte{'|'}, 2)
	if len(parts) < 2 {
		return nil, ErrUncurrectFormat
	}

	return &Packet{
		MessageID: string(parts[0]),
		Part:      1,
		Total:     1,
		Payload:   parts[1],
	}, nil
}

func SplitDataOnPackage(msgID string, payload []byte, size int) []Packet {
	if size < 0 {
		return []Packet{{
			MessageID: msgID,
			Part:      1,
			Total:     1,
			Payload:   payload,
		}}
	}

	packagesCount := (len(payload) + size - 1) / size
	packages := make([]Packet, 0, packagesCount)

	partNum := 1
	for i := 0; i < len(payload); i += size {
		packages = append(packages, Packet{
			MessageID: msgID,
			Part:      partNum,
			Total:     packagesCount,
			Payload:   payload[i:min(i+size, len(payload))],
		})
		partNum++
	}
	return packages
}

type PacketBuf struct {
	mu           sync.Mutex
	total        int
	recivedCount int
	parts        [][]byte
	updatedAt    time.Time
}

type PacketAssembler struct {
	mu       sync.Mutex
	ttl      time.Duration
	sessions map[string]*PacketBuf
	pool     sync.Pool
	done     chan struct{}
}

func NewPackageAssember(ttl time.Duration, packageSize int) *PacketAssembler {
	pa := &PacketAssembler{
		ttl:      ttl,
		sessions: make(map[string]*PacketBuf),
		done:     make(chan struct{}),
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

func (pa *PacketAssembler) AddPacakge(msgID string, part, totalParts int, payload []byte) ([]byte, bool) {
	if part <= 0 || totalParts <= 0 || part > totalParts {
		return nil, false
	}

	pa.mu.Lock()
	sess, exists := pa.sessions[msgID]
	if !exists {
		sess = &PacketBuf{
			mu:           sync.Mutex{},
			total:        totalParts,
			recivedCount: 0,
			parts:        make([][]byte, totalParts+1),
			updatedAt:    time.Now(),
		}
		pa.sessions[msgID] = sess
	}
	pa.mu.Unlock()

	sess.mu.Lock()
	sess.updatedAt = time.Now()
	if sess.parts[part] == nil {
		ptr := pa.pool.Get().(*[]byte)
		buf := *ptr
		n := copy(buf, payload)

		sess.parts[part] = buf[:n]
		sess.recivedCount++
	}
	sess.mu.Unlock()

	if sess.recivedCount == sess.total {
		pa.mu.Lock()
		delete(pa.sessions, msgID)
		pa.mu.Unlock()
		return pa.concat(sess), true
	}
	return nil, false
}

func (pa *PacketAssembler) concat(sess *PacketBuf) []byte {
	lenPyaload := 0
	for _, f := range sess.parts {
		lenPyaload += len(f)
	}
	payload := make([]byte, lenPyaload)
	offset := 0
	for i := 1; i < len(sess.parts); i++ {
		if sess.parts[i] != nil {
			copy(payload[offset:], sess.parts[i])
			offset += len(sess.parts[i])
			bufPtr := &sess.parts[i]
			pa.pool.Put(bufPtr)
			sess.parts[i] = nil
		}
	}

	return payload
}

func (pa *PacketAssembler) backgroundCleaner() {
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

func (pa *PacketAssembler) Stop() {
	select {
	case <-pa.done:
	default:
		close(pa.done)
	}
}
