package localcache

// reference from https://studygolang.com/articles/23183?fr=sidebar

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
	"time"
)

type LRUCache struct {
	maxSlots   int
	maxMemory  int
	usedMemory int

	mu *sync.RWMutex

	// index
	linkedQueue *list.List
	bucket      map[string]*list.Element

	clearInterval int
}

func NewLRUCache(maxSlots, maxMemory int, interval int) Cache {
	cache := &LRUCache{
		maxSlots:  maxSlots,
		maxMemory: maxMemory,

		clearInterval: interval,

		mu:          new(sync.RWMutex),
		linkedQueue: list.New(),
		bucket:      make(map[string]*list.Element),
	}

	// clean in goroutine
	go cache.backgroundCleanForLRUCache()
	return cache
}

func (lru *LRUCache) backgroundCleanForLRUCache() {
	expiredKeys := make([]string, 0, lru.maxSlots)
	for {
		// time to sleep
		time.Sleep(time.Second * time.Duration(lru.clearInterval))

		lru.mu.Lock()
		for k, v := range lru.bucket {
			if v.Value.(*localValue).isExpired() {
				expiredKeys = append(expiredKeys, k)
			}
		}
		lru.mu.Unlock()
		fmt.Println("got expired keys: ", expiredKeys, ", lru.cleaninterval=", lru.clearInterval)

		// do clean work, there should consider the state of too many keys
		for len(expiredKeys) > 0 {
			lru.mu.Lock()
			for idx, _ := range expiredKeys {
				lru.delete(expiredKeys[idx])
			}
			lru.mu.Unlock()
		}
	}
}

func (lru *LRUCache) statistic() {
	template := fmt.Sprintf(`
------------LOCAL_CACHE------------
---USED_SLOTS     = %d
---MAX_SLOTS      = %d
---USED_MEMOTY    = %d
---MAX_MEMORY     = %d
---CLEAN_INTERVAL = %d
-----------------------------------
`, len(lru.bucket), lru.maxSlots, lru.usedMemory, lru.maxMemory, lru.clearInterval)
	fmt.Println(template)
}

func (lru *LRUCache) checkMemory(oldMem, newMem int) bool {
	if lru.usedMemory+newMem-oldMem > lru.maxMemory {
		return false
	}
	return true
}

func (lru *LRUCache) printLinkedQueue() {
	for idx := 0; idx < lru.linkedQueue.Len(); idx++ {
		if ele := lru.linkedQueue.Front(); ele != nil {
			element := ele.Value.(*localValue)
			fmt.Printf("linkedQueue idx=[%d, name=%s, val=%s]\t", idx, element.keyname, element.value.String())
		}
	}
	fmt.Println()
}

func (lru *LRUCache) set(key string, val []byte, ttl time.Duration) (err error) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	fmt.Println("key=", key, " ,val=", string(val), " , usedslots=", len(lru.bucket), " , maxslots=", lru.maxSlots)

	if ele, ok := lru.bucket[key]; ok {
		buf := ele.Value.(*localValue).value
		// reached the maximum memory, need shrink
		for {
			if lru.checkMemory(buf.Len(), len(val)) {
				break
			}
			lru.shrink()
		}
		buf.Reset()
		buf.Write(val)
		lru.linkedQueue.MoveToFront(ele)
		lru.usedMemory += len(val) - ele.Value.(*localValue).value.Len()
		return
	}
	// reached the maximum memory, need shrink
	for {
		if lru.checkMemory(0, len(val)) {
			break
		}
		lru.shrink()
	}

	// update
	buf := cachePool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Write(val)
	element := lru.linkedQueue.PushFront(&localValue{
		keyname: key,
		value:   buf,
		ttl:     time.Now().Add(ttl),
	})
	lru.usedMemory += len(val)
	lru.bucket[key] = element
	// check slots
	if len(lru.bucket) > lru.maxSlots {
		if tail := lru.linkedQueue.Back(); tail != nil {
			lv, _ := tail.Value.(*localValue)
			lru.delete(lv.keyname)
		}
	}

	return nil
}

func (lru *LRUCache) get(key string) (value []byte) {
	lru.mu.RLock()
	ele, ok := lru.bucket[key]
	lru.mu.RUnlock()

	if !ok {
		return []byte{}
	}

	// check if expired
	lv := ele.Value.(*localValue)
	if lv.isExpired() {
		lru.delete(key)
		return []byte{}
	}

	lru.linkedQueue.MoveToFront(ele)
	return lv.value.Bytes()
}

func (lru *LRUCache) delete(key string) {
	ele, ok := lru.bucket[key]
	if !ok {
		return
	}

	// update
	lv := ele.Value.(*localValue)
	lru.usedMemory += lv.value.Len()
	buf := lv.value
	if buf != nil {
		buf.Reset()
		cachePool.Put(buf)
	}
	lru.linkedQueue.Remove(ele)
	delete(lru.bucket, lv.keyname)
}

// shrink for more memory
func (lru *LRUCache) shrink() {
	if tail := lru.linkedQueue.Back(); tail != nil {
		lv, _ := tail.Value.(*localValue)
		fmt.Println("triggered the shrink for ", lv.keyname)
		lru.delete(lv.keyname)
	}
}
