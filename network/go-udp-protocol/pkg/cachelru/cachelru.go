package cachelru

import "sync"

type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

type CacheLRU struct {
	mu    sync.Mutex
	cap   int
	items map[string]*Node
	head  *Node
	tail  *Node
}

func New(capacity int) *CacheLRU {
	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head

	return &CacheLRU{
		mu:    sync.Mutex{},
		cap:   capacity,
		items: make(map[string]*Node, capacity),
		head:  head,
		tail:  tail,
	}
}

func (c *CacheLRU) remove(n *Node) {
	n.next.prev = n.prev
	n.prev.next = n.next
}

func (c *CacheLRU) pushHead(n *Node) {
	n.next = c.head.next
	n.prev = c.head

	c.head.next.prev = n
	c.head.next = n
}

func (c *CacheLRU) Put(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		node.value = value
		c.remove(node)
		c.pushHead(node)
		return
	}

	if len(c.items) >= c.cap {
		oldest := c.tail.prev
		c.remove(oldest)
		delete(c.items, oldest.key)
	}

	newNode := &Node{key: key, value: value}
	c.pushHead(newNode)
	c.items[key] = newNode
}

func (c *CacheLRU) Get(key string, value string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.items[key]
	if !exists {
		return "", false
	}

	c.remove(node)
	c.pushHead(node)

	return node.value, true
}
