package cache

import "sync"

type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

type LRUCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*Node
	head     *Node
	tail     *Node
}

func New(cap int) *LRUCache {
	head, tail := &Node{}, &Node{}
	head.next, tail.prev = tail, head

	return &LRUCache{
		mu:       sync.Mutex{},
		capacity: cap,
		items:    make(map[string]*Node, cap),
		head:     head,
		tail:     tail,
	}
}

func (c *LRUCache) remove(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (c *LRUCache) pushHead(n *Node) {
	n.next = c.head.next
	n.prev = c.head
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRUCache) Get(key string) (string, bool) {
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

func (c *LRUCache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.items[key]; exists {
		node.value = value
		c.remove(node)
		c.pushHead(node)
		return
	}

	if len(c.items) >= c.capacity {
		oldest := c.tail.prev
		c.remove(oldest)
		delete(c.items, oldest.key)
	}

	newNode := &Node{key: key, value: value}
	c.pushHead(newNode)
	c.items[key] = newNode
}
