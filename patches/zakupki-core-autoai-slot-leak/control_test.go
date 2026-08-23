package control

import (
	"context"
	"testing"
)

func TestFreeAnalyzeSlotsAfterRelease(t *testing.T) {
	c := New()
	c.SetAnalyzeCapacity(2)

	_, r1, ok := c.BeginAnalyze(context.Background())
	if !ok {
		t.Fatal("first BeginAnalyze")
	}
	_, r2, ok := c.BeginAnalyze(context.Background())
	if !ok {
		t.Fatal("second BeginAnalyze")
	}
	if c.FreeAnalyzeSlots() != 0 {
		t.Fatalf("free=%d", c.FreeAnalyzeSlots())
	}
	_, _, ok = c.BeginAnalyze(context.Background())
	if ok {
		t.Fatal("third BeginAnalyze must fail at capacity")
	}

	r1()
	if c.FreeAnalyzeSlots() != 1 {
		t.Fatalf("after r1 free=%d", c.FreeAnalyzeSlots())
	}
	r2()
	if c.FreeAnalyzeSlots() != 2 {
		t.Fatalf("after r2 free=%d", c.FreeAnalyzeSlots())
	}
}
