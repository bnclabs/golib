package llrb

import (
	"reflect"
	"testing"
)

func TestAscendGreaterOrEqual(t *testing.T) {
	tree := NewLLRB()
	tree.Insert(KeyInt(4))
	tree.Insert(KeyInt(6))
	tree.Insert(KeyInt(1))
	tree.Insert(KeyInt(3))
	var ary []Key
	tree.AscendGreaterOrEqual(KeyInt(-1), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected := []Key{KeyInt(1), KeyInt(3), KeyInt(4), KeyInt(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(KeyInt(3), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{KeyInt(3), KeyInt(4), KeyInt(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(KeyInt(2), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{KeyInt(3), KeyInt(4), KeyInt(6)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

func TestDescendLessOrEqual(t *testing.T) {
	tree := NewLLRB()
	tree.Insert(KeyInt(4))
	tree.Insert(KeyInt(6))
	tree.Insert(KeyInt(1))
	tree.Insert(KeyInt(3))
	var ary []Key
	tree.DescendLessOrEqual(KeyInt(10), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected := []Key{KeyInt(6), KeyInt(4), KeyInt(3), KeyInt(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(KeyInt(4), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{KeyInt(4), KeyInt(3), KeyInt(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(KeyInt(5), func(i Key) bool {
		ary = append(ary, i)
		return true
	})
	expected = []Key{KeyInt(4), KeyInt(3), KeyInt(1)}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}
