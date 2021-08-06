package merkle

import (
	"testing"
	"time"
)

func TestSortedKeySet(t *testing.T) {
	c1 := map[uint8]*Trie{
		1: nil,
		2: nil,
		3: nil,
	}

	c2 := map[uint8]*Trie{
		2: nil,
		4: nil,
	}

	sorted := sortedKeySet(c1, c2)
	expected := []uint8{1, 2, 3, 4}

	for i, k := range sorted {
		if k != expected[i] {
			t.Fatalf("expected: %v actual: %v", expected, sorted)
		}
	}
}

func TestMurmurCompatibility(t *testing.T) {
	trie := &Trie{}

	Insert(trie, "test_id", 1628125200)
	expected := uint32(1519865632)

	if trie.Hash != expected {
		t.Fatalf("murmur hash imcompatable. expect: %d actual: %d", trie.Hash, expected)
	}
}

func TestInsert(t *testing.T) {
	trie := &Trie{}

	datas := [][]string{
		{"2021-08-05T00:00:00.000Z", "0"},
		{"2021-08-05T01:00:00.000Z", "1"},
		{"2021-08-05T02:00:00.000Z", "2"},
		{"2021-08-05T03:00:00.000Z", "3"},
	}

	for _, d := range datas {
		timestamp, id := d[0], d[1]
		ts, err := time.Parse(time.RFC3339, timestamp)
		ms := ts.UnixNano() / int64(time.Millisecond)

		if err != nil {
			t.Fatalf("Test code error: unable to parse timestamp: %s", timestamp)
		}

		Insert(trie, id, ms)
	}

	expected := uint32(1216950095)
	if trie.Hash != expected {
		t.Errorf("murmur hash imcompatable. expect: %d actual: %d", expected, trie.Hash)
	}
}
