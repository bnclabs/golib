package llrb

import "testing"

func TestDict(t *testing.T) {
	d := NewDict()
	if d.Len() != 0 {
		t.Fatalf("expected an empty dict")
	}
	if d.Size() != 0 {
		t.Fatalf("expected an empty dict")
	}
	d.InsertBulk(
		&KeyInt{9, 90},
		&KeyInt{3, 30},
		&KeyInt{8, 80},
		&KeyInt{1, 10},
	)
	d.UpsertBulk(
		&KeyInt{5, 50},
		&KeyInt{2, 20},
		&KeyInt{7, 70},
		&KeyInt{1, 100},
		&KeyInt{1, 1000},
	)
	if d.Has(&KeyInt{5, 0}) == false {
		t.Fatalf("expected key")
	}
	if d.Has(&KeyInt{6, 0}) == true {
		t.Fatalf("unexpected key")
	}
	if itm := d.Get(&KeyInt{1, 0}).(*KeyInt); itm.Key != 1 || itm.Value != 1000 {
		t.Fatal("failed Get(): %v", itm)
	}
	if item := d.Min().(*KeyInt); item.Key != 1 || item.Value != 1000 {
		t.Fatal("failed Min(): %v", item)
	}
	if item := d.Max().(*KeyInt); item.Key != 9 || item.Value != 90 {
		t.Fatal("failed Max(): %v", item)
	}
	ref := map[int]int64{9: 90, 3: 30, 8: 80, 5: 50, 2: 20, 7: 70, 1: 1000}
	d.Range(nil, nil, "both", func(key Item) bool {
		k := key.(*KeyInt).Key
		if _, ok := ref[int(k)]; !ok {
			t.Fatalf("%v expected", k)
		}
		delete(ref, int(k))
		return true
	})
	if len(ref) > 0 {
		t.Fatalf("expected to range on full set")
	}

	d.DeleteMin()
	d.DeleteMax()
	ref = map[int]int64{3: 30, 8: 80, 5: 50, 2: 20, 7: 70}
	d.Range(nil, nil, "both", func(key Item) bool {
		k := key.(*KeyInt).Key
		if _, ok := ref[int(k)]; !ok {
			t.Fatalf("%v expected", k)
		}
		delete(ref, int(k))
		return true
	})
	if len(ref) > 0 {
		t.Fatalf("expected to range on full set")
	}
}
