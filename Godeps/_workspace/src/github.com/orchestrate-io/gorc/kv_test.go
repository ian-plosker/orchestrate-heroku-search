// Copyright 2014, Orchestrate.IO, Inc.

package gorc

import (
	"testing"
	"testing/quick"
)

func TestKVHasNext(t *testing.T) {
	f := func(results *KVResults) bool {
		return !(results.Next == "" && results.HasNext())
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestKVTrailingGetUri(t *testing.T) {
	f := func(path *Path) bool {
		if path.Ref == "" {
			return path.trailingGetURI() == path.Collection+"/"+path.Key
		}
		return path.trailingGetURI() == path.Collection+"/"+path.Key+"/refs/"+path.Ref
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestKVTrailingPetUri(t *testing.T) {
	f := func(path *Path) bool {
		return path.trailingPutURI() == path.Collection+"/"+path.Key
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
