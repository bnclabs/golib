// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math"
	"math/rand"
	"testing"
)

func TestCases(t *testing.T) {
	tree := NewLLRB()
	tree.Upsert(KeyInt(1))
	tree.Upsert(KeyInt(1))
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(KeyInt(1)) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(KeyInt(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(KeyInt(1)) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(KeyInt(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(KeyInt(1)) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := NewLLRB()
	n := 100
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt(n - i))
	}
	i := 0
	tree.AscendGreaterOrEqual(KeyInt(0), func(key Key) bool {
		i++
		if key.(KeyInt) != KeyInt(i) {
			t.Errorf("bad order: got %d, expect %d", key.(KeyInt), i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := NewLLRB()
	order := []KeyString{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.Upsert(i)
	}
	k := 0
	tree.AscendRange(KeyString("ab"), KeyString("ac"), func(key Key) bool {
		if k > 3 {
			t.Fatalf("returned more keys than expected")
		}
		i1 := order[k]
		i2 := key.(KeyString)
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
		tree.Upsert(KeyInt(perm[i]))
	}
	j := 0
	tree.AscendGreaterOrEqual(KeyInt(0), func(key Key) bool {
		if key.(KeyInt) != KeyInt(j) {
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
		tree.Upsert(KeyInt(perm[i]))
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced := tree.Upsert(KeyInt(perm[i])); replaced == nil || replaced.(KeyInt) != KeyInt(perm[i]) {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := NewLLRB()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt(perm[i]))
	}
	for i := 0; i < n; i++ {
		tree.Delete(KeyInt(i))
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := NewLLRB()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt(perm[i]))
	}
	if tree.Delete(KeyInt(200)) != nil {
		t.Errorf("deleted non-existent key")
	}
	if tree.Delete(KeyInt(-2)) != nil {
		t.Errorf("deleted non-existent key")
	}
	for i := 0; i < n; i++ {
		if u := tree.Delete(KeyInt(i)); u == nil || u.(KeyInt) != KeyInt(i) {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(KeyInt(200)) != nil {
		t.Errorf("deleted non-existent key")
	}
	if tree.Delete(KeyInt(-2)) != nil {
		t.Errorf("deleted non-existent key")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := NewLLRB()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.Upsert(KeyInt(perm[i]))
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(KeyInt(i))
	}
	j := 0
	tree.AscendGreaterOrEqual(KeyInt(0), func(key Key) bool {
		switch j {
		case 0:
			if key.(KeyInt) != KeyInt(0) {
				t.Errorf("expecting 0")
			}
		case 1:
			if key.(KeyInt) != KeyInt(n-1) {
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
		tree.Upsert(KeyInt(perm[i]))
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt(b.N - i))
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(KeyInt(i))
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := NewLLRB()
	for i := 0; i < b.N; i++ {
		tree.Upsert(KeyInt(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := NewLLRB()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.Insert(KeyInt(perm[i]))
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(KeyInt(0), func(key Key) bool {
		if key.(KeyInt) != KeyInt(j/2) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}
