package investor

import (
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"testing"
)

func makeTimeframe(resolution, parent string) _struct.TimeframeItem {
	tf := _struct.TimeframeItem{}
	tf.Config.Resolution = resolution
	tf.Config.Parent = parent
	return tf
}

func TestLinkTimeframeTree(t *testing.T) {
	i := &Investor{}
	i.Timeframes = []_struct.TimeframeItem{
		makeTimeframe("30m", "4h"),
		makeTimeframe("4h", "D"),
		makeTimeframe("D", ""),
	}

	if err := i.linkTimeframeTree(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tf30 := &i.Timeframes[0]
	tf4h := &i.Timeframes[1]
	tfD := &i.Timeframes[2]

	if tf30.Parent != tf4h {
		t.Fatalf("30m parent must be 4h")
	}
	if tf4h.Parent != tfD {
		t.Fatalf("4h parent must be D")
	}
	if tfD.HasParent() {
		t.Fatalf("D must be a root")
	}
	if len(tf4h.Children) != 1 || tf4h.Children[0] != tf30 {
		t.Fatalf("4h must have 30m as its only child")
	}
	if len(tfD.Children) != 1 || tfD.Children[0] != tf4h {
		t.Fatalf("D must have 4h as its only child")
	}
}

func TestLinkTimeframeTreeUnknownParent(t *testing.T) {
	i := &Investor{}
	i.Timeframes = []_struct.TimeframeItem{
		makeTimeframe("30m", "1h"),
		makeTimeframe("D", ""),
	}
	if err := i.linkTimeframeTree(); err == nil {
		t.Fatalf("expected error for unknown parent")
	}
}

func TestLinkTimeframeTreeSelfReference(t *testing.T) {
	i := &Investor{}
	i.Timeframes = []_struct.TimeframeItem{
		makeTimeframe("30m", "30m"),
	}
	if err := i.linkTimeframeTree(); err == nil {
		t.Fatalf("expected error for self reference")
	}
}

func TestLinkTimeframeTreeDuplicateResolution(t *testing.T) {
	i := &Investor{}
	i.Timeframes = []_struct.TimeframeItem{
		makeTimeframe("30m", ""),
		makeTimeframe("30m", ""),
	}
	if err := i.linkTimeframeTree(); err == nil {
		t.Fatalf("expected error for duplicate resolution")
	}
}
