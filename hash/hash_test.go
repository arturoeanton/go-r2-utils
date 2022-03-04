package hash_test

import (
	"testing"

	"github.com/arturoeanton/go-r2-utils/hash"
)

func TestHash(t *testing.T) {
	i := hash.HashStringUint64("test")
	if i != 10090666253179731817 {
		t.Error("HashStringUint64() failed")
	}
}
