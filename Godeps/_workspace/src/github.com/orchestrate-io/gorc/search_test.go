// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"testing"
	"testing/quick"
)

func TestSearchHasNext(t *testing.T) {
	f := func(results *SearchResults) bool {
		return !(results.Next == "" && results.HasNext())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestSearchHasPrev(t *testing.T) {
	f := func(results *SearchResults) bool {
		return !(results.Prev == "" && results.HasPrev())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
