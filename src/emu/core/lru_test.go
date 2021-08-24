// Copyright (c) 2021 Cisco Systems and/or its affiliates.
// Licensed under the Apache License, Version 2.0 (the "License");
// that can be found in the LICENSE file in the root of the source
// tree.

package core

import (
    "fmt"
    "log"
    "testing"
)

func checkLru(lru *Lru, res LruVal, expLen int, expNilRes bool, evicted bool, expEvicted bool) {
    if (res == nil) != expNilRes {
        log.Panicf("got res %v (not as expected)", res)
    }
    if len(lru.elemsMap) != expLen {
        log.Panicf("got len %v, expected %v", len(lru.elemsMap), expLen)
    }
    if evicted != expEvicted {
        log.Panicf("got evicted %v (not as expected)", evicted)
    }
}

func TestLru(t *testing.T) {
    fmt.Println("start TestLru")
    lru := NewLru(2)

    res, evicted := lru.Add(1, "v1")
    checkLru(lru, res, 1, true, evicted, false)
    res, evicted = lru.Add(1, "v1")
    checkLru(lru, res, 1, false, evicted, false)
    res, evicted = lru.Add(2, "v2")
    checkLru(lru, res, 2, true, evicted, false)
    res, evicted = lru.Add(3, "v3")
    checkLru(lru, res, 2, true, evicted, true)
    res, evicted = lru.Add(3, "v3")
    checkLru(lru, res, 2, false, evicted, false)

    fmt.Println("end TestLru")
}

type Ipfix struct {
    clientBytes     uint64
    clientIpV4      uint32
    clientPkts      uint64
    clientPort      uint16
    proto           uint8
    serverBytes     uint64
    serverIpV4      uint32
    serverPkts      uint64
    serverPort      uint16
}

type IpfixKey struct {
    clientIpV4      uint32
    clientPort      uint16
    proto           uint8
    serverIpV4      uint32
    serverPort      uint16
}

func (i *Ipfix) String() string {
    return i.getKey().(string)
}

func (i *Ipfix) getKey() LruKey {
    return IpfixKey{i.clientIpV4,
        i.clientPort,
        i.proto,
        i.serverIpV4,
        i.serverPort,
    }
}

func BenchmarkLru(b *testing.B) {
    fmt.Println("start BenchmarkLru")
    const maxEntries = 2
    lru := NewLru(maxEntries)
    rec := Ipfix{}

    for i:=0; i<b.N; i++ {
        rec.serverIpV4 = uint32(i) // changes the key
        oldRec, evicted := lru.Add(rec.getKey(), &rec)
        if oldRec != nil {
            b.Fatal("oldRec is not nil")
        }
        if i < maxEntries && evicted || i >= maxEntries && !evicted {
            b.Fatalf("i is %v and evicted is %v", i, evicted)
        }
        oldRec, evicted = lru.Add(rec.getKey(), &rec)
        if oldRec == nil || evicted == true {
            b.Fatalf("oldRec is %v and evicted is %v", oldRec, evicted)
        }
        oldIpfix := *oldRec.(*Ipfix)
        if oldIpfix != rec {
            b.Fatalf("old rec is %v and new rec is %v", oldIpfix, rec)
        }
    }
    fmt.Printf("end BenchmarkLru, elems map: %v\n", lru.elemsMap)
}