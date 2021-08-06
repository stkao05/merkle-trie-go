package merkle

import (
	"testing"
	"time"
)

func TestMurmurCompatibility(t *testing.T) {
	trie := &Trie{}

	Insert(trie, "test_id", 1628125200)
	expected := uint32(1519865632)

	if trie.Hash != expected {
		t.Errorf("murmur hash imcompatable. expect: %d actual: %d", trie.Hash, expected)
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
