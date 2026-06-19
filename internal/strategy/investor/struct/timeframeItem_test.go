package _struct

import "testing"

func TestChildShareAmountEvenSplit(t *testing.T) {
	parent := &TimeframeItem{}
	parent.Config.Resolution = "D"
	parent.AllowedChildAmount = 500

	childA := &TimeframeItem{Parent: parent}
	childA.Config.Resolution = "4h"
	childB := &TimeframeItem{Parent: parent}
	childB.Config.Resolution = "30m"
	parent.Children = []*TimeframeItem{childA, childB}

	if !parent.HasChildren() {
		t.Fatalf("parent must report having children")
	}
	if parent.HasParent() {
		t.Fatalf("root parent must not report having a parent")
	}
	if !childA.HasParent() {
		t.Fatalf("child must report having a parent")
	}

	if got := childA.ChildShareAmount(); got != 250 {
		t.Fatalf("childA share expected 250, got %g", got)
	}
	if got := childB.ChildShareAmount(); got != 250 {
		t.Fatalf("childB share expected 250, got %g", got)
	}
}

func TestChildShareAmountNoParent(t *testing.T) {
	root := &TimeframeItem{}
	if got := root.ChildShareAmount(); got != 0 {
		t.Fatalf("root child share must be 0, got %g", got)
	}
}
