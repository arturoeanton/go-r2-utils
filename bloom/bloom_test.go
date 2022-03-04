package bloom_test

import (
	"testing"

	"github.com/arturoeanton/go-r2-utils/bloom"
)

func TestBloom(t *testing.T) {
	bl := bloom.NewBloom()
	bl.Add("HelloWorld.")
	bl.Add("HelloWorld")
	bl.Add("HelloWorld")
	bl.Add("HelloWorld")

	if bl.Count("HelloWorld.") != 1 {
		t.Error("Count failed")
	}
	if bl.Count("HelloWorld") != 3 {
		t.Error("Count failed")
	}
	if !bl.Contains("HelloWorld") {
		t.Error("bloom does not contain HelloWorld ")
	}
	bl.Remove("HelloWorld")
	if bl.Count("HelloWorld") != 2 {
		t.Error("Count failed")
	}
	if !bl.Contains("HelloWorld") {
		t.Error("bloom does not contain HelloWorld ")
	}
	bl.Remove("HelloWorld")
	if bl.Count("HelloWorld") != 1 {
		t.Error("Count failed")
	}
	if !bl.Contains("HelloWorld") {
		t.Error("bloom does not contain HelloWorld ")
	}
	bl.Remove("HelloWorld")
	if bl.Count("HelloWorld") != 0 {
		t.Error("Count failed")
	}
	if bl.Contains("HelloWorld") {
		t.Error("bloom contain HelloWorld ")
	}
	bl.Remove("HelloWorld")
	if bl.Contains("HelloWorld") {
		t.Error("bloom contain HelloWorld ")
	}

	bl.Add("HelloWorld")
	if bl.Count("HelloWorld") != 1 {
		t.Error("Count failed")
	}

	if bl.Count("HelloWorld.") != 1 {
		t.Error("Count failed")
	}

}
