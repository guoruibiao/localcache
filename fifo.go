package localcache

import (
	"bytes"
	"container/list"
	"fmt"
	"sync"
	"time"
)

type FIFOCache struct {
	// statistic info
	maxSlots   int
	maxMemory  int
	usedMemory int
	usedSlots  int

	// read-write locker
	mu *sync.RWMutex
	// interval
	clearInterval int
	queue         *list.List
	// container
	bucket map[string]*list.Element
}

func NewFIFOCache(maxSlots, maxMemory int, interval int) Cache {
	cache := &FIFOCache{
		maxSlots:  maxSlots,
		maxMemory: maxMemory,

		clearInterval: interval,

		mu:     new(sync.RWMutex),
		queue:  list.New(),
		bucket: make(map[string]*list.Element),
	}

	go cache.backgroundCleanForFIFOCache()
	return cache
}

func (f *FIFOCache) backgroundCleanForFIFOCache() {
	expiredKeys := make([]string, 0, f.maxSlots)
	for {
		// time to sleep
		time.Sleep(time.Second * time.Duration(f.clearInterval))

		f.mu.Lock()
		for k, v := range f.bucket {
			if v.Value.(*localValue).isExpired() {
				expiredKeys = append(expiredKeys, k)
			}
		}
		f.mu.Unlock()
		fmt.Println("got expired keys: ", expiredKeys, ", lru.cleaninterval=", f.clearInterval)

		// do clean work, there should consider the state of too many keys
		for len(expiredKeys) > 0 {
			f.mu.Lock()
			for idx, _ := range expiredKeys {
				f.delete(expiredKeys[idx])
			}
			f.mu.Unlock()
		}
	}
}

func (f *FIFOCache) statistic() {
	template := fmt.Sprintf(`
------------LOCAL_CACHE------------
---USED_SLOTS     = %d
---MAX_SLOTS      = %d
---USED_MEMOTY    = %d
---MAX_MEMORY     = %d
---CLEAN_INTERVAL = %d
-----------------------------------
`, f.usedSlots, f.maxSlots, f.usedMemory, f.maxMemory, f.clearInterval)
	fmt.Println(template)
}

func (f *FIFOCache) checkMemory(oldMem, newMem int) bool {
	if f.usedMemory+newMem-oldMem > f.maxMemory {
		return false
	}
	return true
}

func (f *FIFOCache) set(key string, val []byte, ttl time.Duration) (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// if key exists, then update the value
	if ele, ok := f.bucket[key]; ok {
		lv := ele.Value.(*localValue)
		for {
			if f.checkMemory(lv.value.Len(), len(val)) {
				break
			}
			f.shrink()
		}
		buf := lv.value
		buf.Reset()
		buf.Write(val)
		f.queue.PushBack(ele)
		f.usedMemory += len(val) - ele.Value.(*localValue).value.Len()
		return
	}

	// check memory
	for {
		fmt.Println(" checking memory ...")
		if memoryEnough := f.checkMemory(0, len(val)); !memoryEnough {
			fmt.Println("memory not enough, now=", f.usedMemory)
			f.shrink()
		} else {
			fmt.Println("memory enough")
			break
		}
	}

	buf := cachePool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Write(val)
	element := f.queue.PushBack(&localValue{
		value:   buf,
		keyname: key,
		ttl:     time.Now().Add(ttl),
	})
	f.bucket[key] = element
	f.usedMemory += buf.Len()

	// check slots
	for {
		fmt.Println(" checking slots...")
		if len(f.bucket) > f.maxSlots {
			fmt.Println("triggered max slots...")
			f.shrink()
		} else {
			break
		}
	}

	return nil
}

func (f *FIFOCache) get(key string) (value []byte) {
	f.mu.RLock()
	ele, ok := f.bucket[key]
	f.mu.RUnlock()

	if !ok {
		return []byte{}
	}

	lv := ele.Value.(*localValue)
	if lv.isExpired() {
		return []byte{}
	}

	return lv.value.Bytes()
}

func (f *FIFOCache) delete(key string) {
	ele, ok := f.bucket[key]
	if !ok {
		return
	}

	lv := ele.Value.(*localValue)
	buf := lv.value
	f.usedMemory -= buf.Len()
	if buf != nil {
		buf.Reset()
		cachePool.Put(buf)
	}
	f.queue.Remove(ele)
	delete(f.bucket, lv.keyname)
}

func (f *FIFOCache) shrink() {
	if ele := f.queue.Front(); ele != nil {
		fmt.Println("shrinking for ", ele.Value.(*localValue).keyname)
		f.delete(ele.Value.(*localValue).keyname)
	}
}
