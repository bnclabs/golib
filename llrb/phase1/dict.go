package main

import "github.com/prataprc/golib/llrb"

type Dict struct {
	dict map[int]int64
	min  int
	max  int
}

func NewDict() *Dict {
	return &Dict{dict: make(map[int]int64)}
}

func (d Dict) Get(key llrb.Item) llrb.Item {
	kint := key.(*llrb.KeyInt)
	if v, ok := d.dict[int(kint.Key)]; ok {
		return &llrb.KeyInt{kint.Key, v}
	}
	return nil
}

func (d Dict) Min() llrb.Item {
	if len(d.dict) == 0 {
		return nil
	}
	return &llrb.KeyInt{int64(d.min), d.dict[d.min]}
}

func (d Dict) Max() llrb.Item {
	if len(d.dict) == 0 {
		return nil
	}
	return &llrb.KeyInt{int64(d.max), d.dict[d.max]}
}

func (d Dict) DelMin() llrb.Item {
	if len(d.dict) == 0 {
		return nil
	}
	return d.Delete(&llrb.KeyInt{int64(d.min), -1})
}

func (d Dict) DelMax() llrb.Item {
	if len(d.dict) == 0 {
		return nil
	}
	return d.Delete(&llrb.KeyInt{int64(d.min), -1})
}

func (d Dict) Upsert(key llrb.Item) llrb.Item {
	kint := key.(*llrb.KeyInt)
	if v, ok := d.dict[int(kint.Key)]; ok {
		d.dict[int(kint.Key)] = kint.Value
		return &llrb.KeyInt{kint.Key, v}
	}
	d.dict[int(kint.Key)] = kint.Value
	return nil
}

func (d Dict) Insert(key llrb.Item) llrb.Item {
	kint := key.(*llrb.KeyInt)
	d.dict[int(kint.Key)] = kint.Value
	return kint
}

func (d Dict) Delete(key llrb.Item) llrb.Item {
	if len(d.dict) == 0 {
		return nil
	}
	kint := key.(*llrb.KeyInt)
	value, ok := d.dict[int(kint.Key)]
	if ok {
		delete(d.dict, int(kint.Key))
		return &llrb.KeyInt{kint.Key, value}
	}
	return nil
}
