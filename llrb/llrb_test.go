// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

var _ = fmt.Sprintf("dummpy print")

func TestCases(t *testing.T) {
	tree := NewLLRB()
	tree.Upsert(KeyInt{1, -1})
	tree.Upsert(KeyInt{1, -1})
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(KeyInt{1, -1}) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(KeyInt{1, -1})
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(KeyInt{1, -1}) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(KeyInt{1, -1})
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(KeyInt{1, -1}) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := NewLLRB()
	n := 100
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(n - i), -1})
	}
	i := 0
	tree.Range(KeyInt{0, -1}, KeyInt{100, -1}, "high", func(key Key) bool {
		i++
		if key.(KeyInt).key != int64(i) {
			t.Errorf("bad order: got %d, expect %d", key.(KeyInt), i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := NewLLRB()
	order := []KeyString{
		{"ab", -1},
		{"aba", -1},
		{"abc", -1},
		{"a", -1},
		{"aa", -1},
		{"aaa", -1},
		{"b", -1},
		{"a-", -1},
		{"a!", -1},
	}
	for _, i := range order {
		tree.Upsert(i)
	}
	k := 0
	low, high := KeyString{"ab", -1}, KeyString{"ac", -1}
	tree.Range(low, high, "low", func(key Key) bool {
		if k > 3 {
			t.Fatalf("returned more keys than expected")
		}
		i1 := order[k].key
		i2 := key.(KeyString).key
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	tree := NewLLRB()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	j := 0
	tree.Range(KeyInt{-1, -1}, KeyInt{1000, -1}, "none", func(key Key) bool {
		if key.(KeyInt).key != int64(j) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := NewLLRB()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		replaced := tree.Upsert(KeyInt{int64(perm[i]), -1})
		if replaced == nil || replaced.(KeyInt).key != int64(perm[i]) {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := NewLLRB()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	for i := 0; i < n; i++ {
		tree.Delete(KeyInt{int64(i), -1})
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := NewLLRB()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	if tree.Delete(KeyInt{200, -1}) != nil {
		t.Errorf("deleted non-existent key")
	}
	if tree.Delete(KeyInt{-2, -1}) != nil {
		t.Errorf("deleted non-existent key")
	}
	for i := 0; i < n; i++ {
		u := tree.Delete(KeyInt{int64(i), -1})
		if u == nil || u.(KeyInt).key != int64(i) {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(KeyInt{200, -1}) != nil {
		t.Errorf("deleted non-existent key")
	}
	if tree.Delete(KeyInt{-2, -1}) != nil {
		t.Errorf("deleted non-existent key")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := NewLLRB()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(KeyInt{int64(i), -1})
	}
	j := 0
	tree.Range(KeyInt{0, -1}, KeyInt{100, -1}, "low", func(key Key) bool {
		switch j {
		case 0:
			if key.(KeyInt).key != int64(0) {
				t.Errorf("expecting 0")
			}
		case 1:
			if key.(KeyInt).key != int64(n-1) {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestRandomInsertStats(t *testing.T) {
	tree := NewLLRB()
	n := 100000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt{int64(perm[i]), -1})
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := NewLLRB()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.Insert(KeyInt{int64(perm[i]), -1})
		}
	}
	j := 0
	tree.Range(KeyInt{-1, -1}, KeyInt{999, -1}, "high", func(key Key) bool {
		if key.(KeyInt).key != int64(j/2) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func BenchmarkInsert(b *testing.B) {
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt{int64(b.N - i), -1})
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt{int64(b.N - i), -1})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(KeyInt{int64(i), -1})
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt{int64(b.N - i), -1})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}
