package storage

import (
	"testing"
)

func TestStorage(t *testing.T) {
	s := New()

	t.Run("ByKey exact and case-insensitive", func(t *testing.T) {
		val, ok := s.ByKey("HTTP")
		if !ok || val == "" {
			t.Fatalf("expected HTTP to exist, got ok=%v", ok)
		}

		valLower, okLower := s.ByKey("http")
		if !okLower || valLower != val {
			t.Fatalf("expected case-insensitive lookup for http, got %s", valLower)
		}

		_, okNotFound := s.ByKey("NONEXISTENT")
		if okNotFound {
			t.Fatalf("expected false for nonexistent key")
		}
	})

	t.Run("Search substring", func(t *testing.T) {
		// Searching for "traffic" should match HTTP and HTTPS
		res := s.Search("traffic")
		if len(res) == 0 {
			t.Fatalf("expected matches for 'traffic'")
		}

		// Searching for protocol name lowercase
		resHTTP := s.Search("http")
		if len(resHTTP) < 2 {
			t.Fatalf("expected at least 2 matches for 'http' (HTTP and HTTPS), got %d", len(resHTTP))
		}

		// Searching for empty or unmatched
		resNone := s.Search("supercalifragilistic")
		if len(resNone) != 0 {
			t.Fatalf("expected 0 matches, got %d", len(resNone))
		}
	})

	t.Run("List and Count", func(t *testing.T) {
		list := s.List()
		count := s.Count()

		if len(list) != count {
			t.Fatalf("len(list) %d != count %d", len(list), count)
		}
		if count == 0 {
			t.Fatalf("storage should not be empty")
		}
	})
}
