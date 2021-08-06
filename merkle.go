package merkle

import (
	"hash"
	"sort"
	"strconv"

	"github.com/spaolacci/murmur3"
)

type Trie struct {
	Hash     uint32
	Children map[uint8]*Trie
}

var Hasher hash.Hash32 = murmur3.New32WithSeed(0)
var Slot int64 = 60 * 1000

func Insert(trie *Trie, id string, timestamp int64) {
	defer Hasher.Reset()

	timeslot := strconv.FormatInt(int64(timestamp/Slot), 10)
	Hasher.Write([]byte(id))
	hash := Hasher.Sum32()

	if trie.Hash == 0 {
		trie.Hash = hash
	} else {
		trie.Hash = trie.Hash ^ hash
	}

	node := trie
	for _, c := range timeslot {
		child := node.Children[uint8(c)]

		if child == nil {
			child = &Trie{Hash: hash}
		} else {
			child.Hash = child.Hash ^ hash
		}

		node = child
	}
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

func BranchPoint(trie1 *Trie, trie2 *Trie) (timestamp int64) {
// 	node1 := trie1
// 	node2 := trie2
// 	time := ""
// 
	// 	for {
	//     sorted :=  sortedKeySet(node1.Children, node2.Children)
	//
	//     var diffKey uint8
	//     for k := range(sorted) {
	//       var h1, h2 uint32
	//
	//
	//
	//
	//     }
	//
	//
	//
	// 	}

	return 0
}
