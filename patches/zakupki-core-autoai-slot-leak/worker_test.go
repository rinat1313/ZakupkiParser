package autoai

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/rinat1313/zakupki-core/internal/control"
)

// TestPreheldSlotMustReleaseEvenOnEarlyReturn фиксирует контракт auto-AI:
// BeginAnalyze в tick() + Release в AnalyzeOptions обязан освобождаться
// даже если AnalyzeTender выходит до вызова LLM (иначе FreeAnalyzeSlots=0 навсегда).
func TestPreheldSlotMustReleaseEvenOnEarlyReturn(t *testing.T) {
	ctrl := control.New()
	ctrl.SetAnalyzeCapacity(1)

	ctx, release, ok := ctrl.BeginAnalyze(context.Background())
	if !ok {
		t.Fatal("BeginAnalyze failed")
	}
	if ctrl.FreeAnalyzeSlots() != 0 {
		t.Fatalf("expected 0 free slots while held, got %d", ctrl.FreeAnalyzeSlots())
	}

	// Имитация исправленного AnalyzeTender: defer release ДО любой ранней ошибки.
	err := func() error {
		defer release()
		_ = ctx
		return errNoText // early return path (пустой корпус / ListDocuments)
	}()
	if err == nil {
		t.Fatal("expected early error")
	}

	if ctrl.FreeAnalyzeSlots() != 1 {
		t.Fatalf("slot leak: free=%d active=%d (auto-AI would freeze)", ctrl.FreeAnalyzeSlots(), ctrl.AnalyzeActiveCount())
	}
	if ctrl.AnalyzeActive() {
		t.Fatal("analyze still marked active after release")
	}
}

func TestReleaseIsIdempotent(t *testing.T) {
	ctrl := control.New()
	ctrl.SetAnalyzeCapacity(1)
	_, release, ok := ctrl.BeginAnalyze(context.Background())
	if !ok {
		t.Fatal("BeginAnalyze failed")
	}
	release()
	release() // sync.Once — не должно увести active в минус без восстановления
	if ctrl.FreeAnalyzeSlots() != 1 {
		t.Fatalf("free=%d after double release", ctrl.FreeAnalyzeSlots())
	}
}

func TestWorkerDoesNotClaimWhenNoFreeSlots(t *testing.T) {
	ctrl := control.New()
	ctrl.SetAnalyzeCapacity(1)
	_, release, ok := ctrl.BeginAnalyze(context.Background())
	if !ok {
		t.Fatal("BeginAnalyze failed")
	}

	var claimed atomic.Int32
	// Как tick(): при free<=0 новых claim нет.
	free := ctrl.FreeAnalyzeSlots()
	if free > 0 {
		claimed.Add(1)
	}
	if claimed.Load() != 0 {
		t.Fatal("must not claim when slots exhausted")
	}

	release()
	if ctrl.FreeAnalyzeSlots() != 1 {
		t.Fatal("slot not freed")
	}
}
