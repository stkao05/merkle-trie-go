package merkle

import (
	"hash"
	"sort"
	"strconv"
	"time"

	"github.com/spaolacci/murmur3"
)

type Trie struct {
	Hash     uint32
	Children map[uint8]*Trie
}

var Hasher hash.Hash32 = murmur3.New32WithSeed(0)
var Slot int64 = 60 * 1000

func Insert(trie *Trie, id string, timestamp time.Time) {
	defer Hasher.Reset()

	ms := timestamp.UnixNano() / int64(time.Microsecond)
	timestr := strconv.FormatInt(ms/Slot, 10)
	Hasher.Write([]byte(id))
	hash := Hasher.Sum32()

	if trie.Hash == 0 {
		trie.Hash = hash
	} else {
		trie.Hash = trie.Hash ^ hash
	}

	node := trie
	for _, c := range timestr {
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

func BranchPoint(trie1 *Trie, trie2 *Trie) (time.Time, error) {
	node1 := trie1
	node2 := trie2
	timestr := ""

	for {
		sorted := sortedKeySet(node1.Children, node2.Children)

		diffKey := -1
		for _, k := range sorted {
			c1, c2 := node1.Children[k], node2.Children[k]

			if (c1 != nil && c2 != nil) && (c1.Hash != c2.Hash) {
				diffKey = int(k)
				break
			}

			if c1 == nil && c2 == nil {
				break
			}

			// either one of c1 and c2 has value
			diffKey = int(k)
		}

		if diffKey == -1 {
			break
		} else {
			timestr += strconv.Itoa(diffKey)
		}
	}

	if timestr == "" {
		return time.Time{}, nil
	}

	ms, err := strconv.Atoi(timestr)
	return time.Unix(0, int64(ms) * int64(time.Microsecond)), err
}
