package localcache

import (
	"bytes"
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type LFUCache struct {
	maxSlots   int
	maxMemory  int
	usedMemory int

	mu *sync.RWMutex

	// custom min heap
	minHeap *MinHeap
	bucket  map[string]*localValue

	clearInterval int
}

func NewLFUCache(maxSlots, maxMemory int, interval int) Cache {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	cache := &LFUCache{
		maxSlots:  maxSlots,
		maxMemory: maxMemory,

		clearInterval: interval,

		mu:      new(sync.RWMutex),
		minHeap: minHeap,
		bucket:  make(map[string]*localValue),
	}

	go cache.backgroundCleanForLFUCache()
	return cache
}

func (lfu *LFUCache) backgroundCleanForLFUCache() {
	expiredKeys := make([]string, 0, lfu.maxSlots)
	for {
		// time to sleep
		time.Sleep(time.Second * time.Duration(lfu.clearInterval))

		lfu.mu.Lock()
		for k, v := range lfu.bucket {
			if v.isExpired() {
				expiredKeys = append(expiredKeys, k)
			}
		}
		lfu.mu.Unlock()
		fmt.Println("got expired keys: ", expiredKeys, ", lru.cleaninterval=", lfu.clearInterval)

		// do clean work, there should consider the state of too many keys
		for len(expiredKeys) > 0 {
			lfu.mu.Lock()
			for idx, _ := range expiredKeys {
				lfu.delete(expiredKeys[idx])
			}
			lfu.mu.Unlock()
		}
	}
}

func (lfu *LFUCache) statistic() {
	template := fmt.Sprintf(`
------------LOCAL_CACHE------------
---USED_SLOTS     = %d
---MAX_SLOTS      = %d
---USED_MEMOTY    = %d
---MAX_MEMORY     = %d
---CLEAN_INTERVAL = %d
-----------------------------------
`, len(lfu.bucket), lfu.maxSlots, lfu.usedMemory, lfu.maxMemory, lfu.clearInterval)
	fmt.Println(template)
}

func (lfu *LFUCache) printBucket() {
	for _, item := range *lfu.minHeap {
		fmt.Printf("[%s:%s:%d], ", item.keyname, item.value.String(), item.frequency)
	}
	fmt.Println()
}

func (lfu *LFUCache) set(key string, val []byte, ttl time.Duration) (err error) {
	lfu.mu.Lock()
	defer lfu.mu.Unlock()

	// check if exists
	if oldLv, ok := lfu.bucket[key]; ok {
		buf := oldLv.value
		// reached the maximum memory, need shrink
		for {
			fmt.Println("exists check memory...")
			if lfu.checkMemory(buf.Len(), len(val)) {
				break
			}
			lfu.shrink()
		}
		buf.Reset()
		buf.Write(val)
		oldLv.value = buf
		oldLv.ttl = time.Now().Add(ttl)
		if string(val) != buf.String() {
			oldLv.frequency = 1
		}
		heap.Fix(lfu.minHeap, lfu.minHeap.Index(oldLv))
		lfu.bucket[key] = oldLv
		return
	}

	// check memory
	for {
		fmt.Println("new check memory..., cur-memory=", lfu.usedMemory)
		if lfu.checkMemory(0, len(val)) {
			break
		}
		lfu.shrink()
	}
	buf := cachePool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Write(val)
	newLv := &localValue{
		keyname:   key,
		value:     buf,
		ttl:       time.Now().Add(ttl),
		frequency: 1,
	}
	lfu.usedMemory += buf.Len()
	heap.Push(lfu.minHeap, newLv)
	lfu.bucket[key] = newLv
	// check slots
	for {
		fmt.Println("check slots ..., cur-length=", len(lfu.bucket))
		if len(lfu.bucket) > lfu.maxSlots {
			lfu.shrink()
		} else {
			break
		}
		//time.Sleep(time.Second * 5)
	}

	return nil
}

func (lfu *LFUCache) get(key string) (value []byte) {
	lfu.mu.RLock()
	lv, ok := lfu.bucket[key]
	lfu.mu.RUnlock()

	if !ok {
		return []byte{}
	}

	if lv.isExpired() {
		return []byte{}
	}
	return lv.value.Bytes()
}

func (lfu *LFUCache) delete(key string) {
	lv, ok := lfu.bucket[key]
	if !ok {
		return
	}

	heap.Remove(lfu.minHeap, lv.frequency)
	delete(lfu.bucket, lv.keyname)
	lfu.usedMemory -= lv.value.Len()
}

func (lfu *LFUCache) checkMemory(oldMem, newMem int) bool {
	if lfu.usedMemory+newMem-oldMem > lfu.maxMemory {
		return false
	}
	return true
}

func (lfu *LFUCache) shrink() {
	if lfu.minHeap.Len() <= 0 {
		return
	}

	lv, ok := lfu.minHeap.Top().(*localValue)
	fmt.Printf("shrink-shrink-shrink=%#v\n", lv)
	if !ok {
		return
	}
	lfu.minHeap.Pop()
	fmt.Println("shrink the name of =", lv.keyname)
	delete(lfu.bucket, lv.keyname)
	lfu.usedMemory -= lv.value.Len()
}
