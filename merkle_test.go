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

// data is a slice of [id, timestamp] tuple
func createTrie(datas [][]string) (*Trie, error) {
	trie := Trie{}

	for _, d := range datas {
		id, timestamp := d[0], d[1]
		ts, err := time.Parse(time.RFC3339, timestamp)

		if err != nil {
			return nil, err
		}

		Insert(&trie, id, ts)
	}

	return &trie, nil
}

func TestMurmurCompatibility(t *testing.T) {
	trie := &Trie{}
	ts := time.Unix(0, 1628125200*int64(time.Microsecond))

	Insert(trie, "test_id", ts)
	expected := uint32(1519865632)

	if trie.Hash != expected {
		t.Fatalf("murmur hash imcompatable. expect: %d actual: %d", trie.Hash, expected)
	}
}

func TestInsert(t *testing.T) {
	datas := [][]string{
		{"0", "2021-08-05T00:00:00.000Z"},
		{"1", "2021-08-05T01:00:00.000Z"},
		{"2", "2021-08-05T02:00:00.000Z"},
		{"3", "2021-08-05T03:00:00.000Z"},
	}

	trie, err := createTrie(datas)

	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(1216950095)
	if trie.Hash != expected {
		t.Errorf("murmur hash imcompatable. expect: %d actual: %d", expected, trie.Hash)
	}
}
