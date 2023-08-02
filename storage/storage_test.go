package storage

import (
	"testing"
)

func TestFakeStorageMethod(t *testing.T) {
	nfs := NewFakeStorage()
	test_out := len(nfs.SearchRecipes(""))
	if test_out != 2 {
		t.Errorf("Expected 2 results got %d", test_out)
	}

	test_out = len(nfs.SearchRecipes("chicken"))
	if test_out != 1 {
		t.Errorf("Expected 1 results got %d", test_out)
	}
}
