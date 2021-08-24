// Copyright (c) 2021 Cisco Systems and/or its affiliates.
// Licensed under the Apache License, Version 2.0 (the "License");
// that can be found in the LICENSE file in the root of the source
// tree.

package core

import (
    "unsafe"
)

/* LRU
Uses Dlist + map
*/

type LruKey interface {}
type LruVal interface {}
type lruElem struct {
    dlist DList
    LruKey
    LruVal
}

type Lru struct {
    elemsList  DList
    elemsMap map[LruKey]*lruElem
    maxEntries int
}

func toElem(dlist *DList) *lruElem {
    return (*lruElem)(unsafe.Pointer(dlist))
}

func NewLru(maxEntries int) *Lru {
    lru := &Lru{
        elemsMap:   make(map[LruKey]*lruElem),
        maxEntries: maxEntries,
    }
    lru.elemsList.SetSelf()
    return lru
}

func (l *Lru) Add(key LruKey, val LruVal) (_ LruVal, evicted bool) {
    if l.elemsMap == nil {
        panic("init by NewLru")
    }
    if oldElem, present := l.elemsMap[key]; present { // elem present
        l.elemsList.RemoveNode(&oldElem.dlist)
        l.elemsList.AddFirst(&oldElem.dlist)
        return oldElem.LruVal, false
    }
    if len(l.elemsMap) >= l.maxEntries {
        oldestElem := toElem(l.elemsList.RemoveLast())
        delete(l.elemsMap, oldestElem.LruKey)
        evicted = true
    }
    newElem := &lruElem{dlist: DList{}, LruKey: key, LruVal: val}
    l.elemsList.AddFirst(&newElem.dlist)
    l.elemsMap[key] = newElem
    return
}
