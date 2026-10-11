package assembly

import (
	"sync"
	"time"
)

type messageBuffer struct {
	total     int
	parts     map[int][]byte
	updatedAt time.Time
}

type Assembler struct {
	mu      sync.Mutex
	ttl     time.Duration
	buffers map[string]*messageBuffer
	stopCh  chan struct{}
}

func NewAssembler(ttl time.Duration) *Assembler {
	a := &Assembler{
		ttl:     ttl,
		buffers: make(map[string]*messageBuffer),
		stopCh:  make(chan struct{}),
	}

	// For clean
	go a.startCleaner()

	return a
}

func (a *Assembler) AddChunk(sessionKey string, part, total int, payload []byte) ([]byte, bool) {
	if part <= 0 || total <= 0 || part > total {
		return nil, false
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	buf, exists := a.buffers[sessionKey]
	if !exists {
		buf = &messageBuffer{
			total:     total,
			parts:     make(map[int][]byte, total),
			updatedAt: time.Now(),
		}
		a.buffers[sessionKey] = buf
	}

	chunkCopy := make([]byte, len(payload))
	copy(chunkCopy, payload)

	buf.parts[part] = chunkCopy
	buf.updatedAt = time.Now()

	if len(buf.parts) == buf.total {
		totalLen := 0
		for i := 1; i <= buf.total; i++ {
			partBytes, ok := buf.parts[i]
			if !ok {
				return nil, false
			}
			totalLen += len(partBytes)
		}

		assembled := make([]byte, 0, totalLen)
		for i := 1; i <= buf.total; i++ {
			assembled = append(assembled, buf.parts[i]...)
		}

		delete(a.buffers, sessionKey)
		return assembled, true
	}

	return nil, false
}

func (a *Assembler) startCleaner() {
	ticker := time.NewTicker(a.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			return
		case now := <-ticker.C:
			a.mu.Lock()
			for key, buf := range a.buffers {
				if now.Sub(buf.updatedAt) > a.ttl {
					delete(a.buffers, key)
				}
			}
			a.mu.Unlock()
		}
	}
}

func (a *Assembler) Stop() {
	select {
	case <-a.stopCh:
	default:
		close(a.stopCh)
	}
}
