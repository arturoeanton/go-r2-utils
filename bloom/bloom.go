package bloom

import (
	"fmt"

	"github.com/arturoeanton/go-r2-utils/hash"
	"github.com/arturoeanton/go-r2-utils/memset"
)

type Bloom struct {
	memset *memset.MemSet
	repeat map[uint64]uint64
	fxHash func(string) uint64
}

func NewBloom() *Bloom {
	return &Bloom{memset: memset.NewMemSet(), fxHash: hash.HashStringUint64, repeat: make(map[uint64]uint64)}
}

// return count repeated hash and true if exist hash collision
func (b *Bloom) Add(param interface{}) (uint64, bool) {
	value := fmt.Sprint(param)
	hashValue := b.fxHash(value)
	if b.memset.Contains(hashValue) {
		if _, ok := b.repeat[hashValue]; !ok {
			b.repeat[hashValue] = 1
			return 2, true
		}
		b.repeat[hashValue]++
		return (b.repeat[hashValue] + 1), true
	}
	b.memset.Add(hashValue)
	return 1, false
}

func (b *Bloom) Remove(param interface{}) uint64 {
	value := fmt.Sprint(param)
	hashValue := b.fxHash(value)
	if b.memset.Contains(hashValue) {
		if _, ok := b.repeat[hashValue]; ok {
			b.repeat[hashValue]--
			l := b.repeat[hashValue]
			if b.repeat[hashValue] <= 0 {
				delete(b.repeat, hashValue)
			}
			return (l + 1)
		}
	}
	b.memset.Remove(hashValue)
	return 0
}

func (b *Bloom) Contains(param interface{}) bool {
	value := fmt.Sprint(param)
	return b.memset.Contains(b.fxHash(value))
}

func (b *Bloom) Count(param interface{}) uint64 {
	value := fmt.Sprint(param)
	hashValue := b.fxHash(value)
	if b.memset.Contains(hashValue) {
		if _, ok := b.repeat[hashValue]; ok {
			return (b.repeat[hashValue] + 1)
		}
		return 1
	}
	return 0
}
