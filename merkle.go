package merkle

import (
	"hash"
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
