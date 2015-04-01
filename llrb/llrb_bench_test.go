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
