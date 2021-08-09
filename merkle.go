package merkle

import (
	"hash"
	"math"
	"sort"
	"time"

	"github.com/spaolacci/murmur3"
)

type Trie struct {
	Hash     uint32
	Children map[uint8]*Trie
}

var Hasher hash.Hash32 = murmur3.New32WithSeed(0)
var SlotMilli int64 = 60 * 1000

func Insert(trie *Trie, id string, timestamp time.Time) {
	defer Hasher.Reset()

	Hasher.Write([]byte(id))
	hash := Hasher.Sum32()

	if trie.Hash == 0 {
		trie.Hash = hash
	} else {
		trie.Hash = trie.Hash ^ hash
	}

	milli := timestamp.UnixNano() / int64(time.Millisecond) / SlotMilli
	ds := numToDigits(milli)
	node := trie

	for _, d := range ds {
		digit := uint8(d)

		if node.Children == nil {
			node.Children = map[uint8]*Trie{}
		}

		child := node.Children[digit]

		if child == nil {
			child = &Trie{Hash: hash}
			node.Children[digit] = child
		} else {
			child.Hash = child.Hash ^ hash
		}

		node = child
	}
}

func BranchPoint(trie1 *Trie, trie2 *Trie) time.Time {
	node1 := trie1
	node2 := trie2
	digits := []uint8{}

loop:
	for {
		switch {
		case node1 == nil && node2 == nil:
			break loop
		case node1 != nil && node2 == nil:
			digits = append(digits, minTime(node1)...)
			break loop
		case node1 == nil && node2 != nil:
			digits = append(digits, minTime(node2)...)
			break loop
		}

		sorted := sortedKeySet(node1.Children, node2.Children)

		// find the smallest digit which hash diff
		foundDiff := false

		for _, d := range sorted {
			c1, c2 := node1.Children[d], node2.Children[d]

			if (c1 != nil && c2 == nil) ||
				(c1 == nil && c2 != nil) ||
				(c1 != nil && c2 != nil && c1.Hash != c2.Hash) {
				digits = append(digits, d)
				node1 = c1
				node2 = c2
				foundDiff = true
				break
			}
		}

		if !foundDiff {
			break loop
		}
	}

	if len(digits) == 0 {
		return time.Time{}
	}

	nano := digitsToNum(digits) * SlotMilli * int64(time.Millisecond)
	return time.Unix(0, nano)
}

func numToDigits(num int64) []int64 {
	digitCount := int(math.Log10(float64(num))) + 1
	res := make([]int64, 0, digitCount)

	for i := digitCount - 1; i >= 0; i-- {
		digit := int64(math.Mod(float64(num)/math.Pow(10, float64(i)), 10))
		res = append(res, digit)
	}

	return res
}

func digitsToNum(digits []uint8) int64 {
	var res = int64(0)
	var n = len(digits)

	for i := 0; i < n; i++ {
		res += int64(digits[i]) * int64(math.Pow10(n-i-1))
	}

	return res
}

func sortedKeySet(t1 map[uint8]*Trie, t2 map[uint8]*Trie) []uint8 {
	var set = map[int]bool{}

	for k := range t1 {
		set[int(k)] = true
	}
	for k := range t2 {
		set[int(k)] = true
	}

	keys := make([]int, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	uintKeys := make([]uint8, 0, len(set))
	for _, k := range keys {
		uintKeys = append(uintKeys, uint8(k))
	}

	return uintKeys
}

func sortedKeys(t map[uint8]*Trie) []uint8 {
	keys := make([]int, 0, len(t))
	for k := range t {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	uintKeys := make([]uint8, 0, len(keys))
	for _, k := range keys {
		uintKeys = append(uintKeys, uint8(k))
	}

	return uintKeys
}

func minTime(trie *Trie) []uint8 {
	keys := []uint8{}

	for trie != nil {
		if trie.Children == nil {
			break
		}

		k := sortedKeys(trie.Children)[0]
		keys = append(keys, k)
		trie = trie.Children[k]
	}

	return keys
}
