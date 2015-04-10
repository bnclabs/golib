package llrb

import "sort"

type Dict struct {
	dict map[int]int64
	size int
}

func NewDict() *Dict {
	return &Dict{dict: make(map[int]int64)}
}

func (d *Dict) Clone() *Dict {
	newd := &Dict{dict: make(map[int]int64), size: d.size}
	for k, v := range d.dict {
		newd.dict[k] = v
	}
	return newd
}

func (d *Dict) Len() int {
	return len(d.dict)
}

func (d *Dict) Size() int {
	return d.size
}

func (d *Dict) Has(key Item) bool {
	kint := key.(*KeyInt).Key
	_, ok := d.dict[int(kint)]
	return ok
}

func (d *Dict) Get(key Item) Item {
	kint := key.(*KeyInt).Key
	if v, ok := d.dict[int(kint)]; ok {
		return &KeyInt{kint, v}
	}
	return nil
}

func (d *Dict) Min() Item {
	if len(d.dict) == 0 {
		return nil
	}
	kint := d.sorted()[0]
	return &KeyInt{int64(kint), d.dict[kint]}
}

func (d *Dict) Max() Item {
	if len(d.dict) == 0 {
		return nil
	}
	keys := d.sorted()
	kint := keys[len(keys)-1]
	return &KeyInt{int64(kint), d.dict[kint]}
}

func (d *Dict) DeleteMin() Item {
	if len(d.dict) == 0 {
		return nil
	}
	min := d.Min().(*KeyInt).Key
	return d.Delete(&KeyInt{min, -1})
}

func (d *Dict) DeleteMax() Item {
	if len(d.dict) == 0 {
		return nil
	}
	max := d.Max().(*KeyInt).Key
	return d.Delete(&KeyInt{max, -1})
}

func (d *Dict) UpsertBulk(keys ...Item) {
	for _, key := range keys {
		d.Upsert(key)
	}
}

func (d *Dict) Upsert(key Item) Item {
	kint := key.(*KeyInt)
	if v, ok := d.dict[int(kint.Key)]; ok {
		d.dict[int(kint.Key)] = kint.Value
		return &KeyInt{kint.Key, v}
	}
	d.dict[int(kint.Key)] = kint.Value
	d.size += 16
	return nil
}

func (d *Dict) InsertBulk(keys ...Item) {
	for _, key := range keys {
		d.Insert(key)
	}
}

func (d *Dict) Insert(key Item) {
	kint := key.(*KeyInt)
	d.dict[int(kint.Key)] = kint.Value
	d.size += 16
}

func (d *Dict) Delete(key Item) Item {
	if len(d.dict) == 0 {
		return nil
	}
	kint := key.(*KeyInt)
	value, ok := d.dict[int(kint.Key)]
	if ok {
		d.size -= 16
		delete(d.dict, int(kint.Key))
		return &KeyInt{kint.Key, value}
	}
	return nil
}

func (d *Dict) Range(low, high Item, incl string, iter KeyIterator) {
	keys := d.sorted()
	lkey, hkey := 0, len(keys)
	i := 0
	if low != nil {
		lint := low.(*KeyInt).Key
		switch incl {
		case "low", "both":
			for ; i < len(d.dict); i++ {
				if int64(keys[i]) < lint {
					continue
				}
				break
			}
			lkey = i
		default:
			for ; i < len(d.dict); i++ {
				if int64(keys[i]) <= lint {
					continue
				}
				break
			}
			lkey = i
		}
	}
	if high != nil {
		hint := high.(*KeyInt).Key
		switch incl {
		case "high", "both":
			for ; i < len(d.dict); i++ {
				if int64(keys[i]) <= hint {
					continue
				}
				break
			}
			hkey = i
		default:
			for ; i < len(d.dict); i++ {
				if int64(keys[i]) < hint {
					continue
				}
				break
			}
			hkey = i
		}
	}
	for i := lkey; i < hkey; i++ {
		key := &KeyInt{int64(keys[i]), d.dict[keys[i]]}
		if !iter(key) {
			return
		}
	}
}

func (d *Dict) sorted() []int {
	keys := make([]int, 0, len(d.dict))
	for key, _ := range d.dict {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}
