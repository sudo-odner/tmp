package cache

import "testing"

func TestLRUCache(t *testing.T) {
	c := New(2)

	// Test basic put and get
	c.Put("k1", "v1")
	c.Put("k2", "v2")

	if val, ok := c.Get("k1"); !ok || val != "v1" {
		t.Fatalf("expected k1=v1, got val=%s, ok=%v", val, ok)
	}

	if val, ok := c.Get("k2"); !ok || val != "v2" {
		t.Fatalf("expected k2=v2, got val=%s, ok=%v", val, ok)
	}

	// Update existing key
	c.Put("k1", "v1_updated")
	if val, ok := c.Get("k1"); !ok || val != "v1_updated" {
		t.Fatalf("expected k1=v1_updated, got val=%s, ok=%v", val, ok)
	}

	// Eviction test: accessing k1 makes k2 the oldest.
	// Adding k3 should evict k2.
	c.Put("k3", "v3")

	if _, ok := c.Get("k2"); ok {
		t.Fatalf("expected k2 to be evicted")
	}

	if val, ok := c.Get("k1"); !ok || val != "v1_updated" {
		t.Fatalf("expected k1 to be present")
	}

	if val, ok := c.Get("k3"); !ok || val != "v3" {
		t.Fatalf("expected k3 to be present")
	}
}
