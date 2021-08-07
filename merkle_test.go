package merkle

import (
	"testing"
	"time"
)

func sliceEquals(s1 []int64, s2 []int64) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func equalDigits(s1 []uint8, s2 []uint8) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestNumToDigits(t *testing.T) {
	res := numToDigits(182719)
	expected := []int64{1, 8, 2, 7, 1, 9}

	if !sliceEquals(res, expected) {
		t.Fatalf("expected: %v actual: %v", expected, res)
	}

	res = numToDigits(2)
	expected = []int64{2}

	if !sliceEquals(res, expected) {
		t.Fatalf("expected: %v actual: %v", expected, res)
	}
}

func TestDigitsToNum(t *testing.T) {
	digits := []uint8{1, 8, 2, 7, 1, 9}
	num := digitsToNum(digits)

	if num != 182719 {
		t.Fatalf("expected: %v actual: %v", num, 182719)
	}

	digits = []uint8{1}
	num = digitsToNum(digits)

	if num != 1 {
		t.Fatalf("expected: %v actual: %v", num, 182719)
	}
}

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

func TestMinTime(t *testing.T) {
	trie := Trie{Children: map[uint8]*Trie{
		1: {
			Children: map[uint8]*Trie{
				6: {},
				7: {},
			},
		},
		2: {
			Children: map[uint8]*Trie{
				3: {},
				4: {},
			},
		},
	}}
	time := minTime(&trie)

	if len(time) != 2 || !(time[0] == 1 && time[1] == 6) {
		t.Fatalf("expected: %v actual: %v", []uint{1, 6}, time)
	}

	trie = Trie{}
	time = minTime(&trie)

	if len(time) != 0 {
		t.Fatalf("expected: %v actual: %v", []uint{}, time)
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

func TestInsert1(t *testing.T) {
	trie, err := createTrie([][]string{
		{"0", "2021-08-05T00:00:00Z"},
	})

	if err != nil {
		t.Fatalf("%v", err)
	}

	digits := minTime(trie)
	expected := []uint8{2, 7, 1, 3, 5, 3, 6, 0}

	if !equalDigits(expected, digits) {
		t.Errorf("expect: %d actual: %d", expected, digits)
	}
}

func TestInsert2(t *testing.T) {
	trie, err := createTrie([][]string{
		{"0", "2021-08-05T00:00:00Z"},
		{"1", "2021-08-05T01:00:00Z"},
		{"2", "2021-08-05T02:00:00Z"},
		{"3", "2021-08-05T03:00:00Z"},
	})

	if err != nil {
		t.Fatal(err)
	}

	expected := uint32(1216950095)
	if trie.Hash != expected {
		t.Errorf("murmur hash imcompatable. expect: %d actual: %d", expected, trie.Hash)
	}
}

func TestBranchPoint(t *testing.T) {
	trie1, err := createTrie([][]string{
		{"0", "2021-08-05T00:00:00Z"},
		{"1", "2021-08-05T01:00:00Z"},
		{"2", "2021-08-05T02:00:00Z"},
		{"3", "2021-08-05T03:00:00Z"},
	})

	if err != nil {
		t.Fatal(err)
	}

	trie2, err := createTrie([][]string{
		{"0", "2021-08-05T00:00:00Z"},
		{"4", "2021-08-05T04:00:00Z"},
	})

	ts := BranchPoint(trie1, trie2)

	actual := ts.UTC().Format(time.RFC3339)
	expect := "2021-08-05T01:00:00Z"

	if expect != actual {
		t.Errorf("expect: %s actual: %s", expect, actual)
	}
}
