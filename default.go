package localcache

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type Default struct {
	// statistic info
	maxSlots int
	maxMemory int
	usedMemory int
	usedSlots int

	// read-write locker
	mu sync.RWMutex
	// interval
	clearInterval int
	// container
	bucket map[string]localValue
}



var cachePool = sync.Pool{
	New: func() interface{}{
		return &bytes.Buffer{}
	},
}

func NewDefaultLocalCache(maxSlots, maxMemory, interval int) Cache {
	cacher := &Default{
		maxSlots: maxSlots,
		maxMemory: maxMemory,
		usedMemory: 0,
		usedSlots: 0,
		mu: sync.RWMutex{},
		clearInterval: interval,
		bucket:make(map[string]localValue),
	}


	// background clean goroutine
	go cacher.backgroundClean()

	return cacher
}

func (d *Default)statistic() {
	template := fmt.Sprintf(`
------------LOCAL_CACHE------------
---USED_SLOTS     = %d
---MAX_SLOTS      = %d
---USED_MEMOTY    = %d
---MAX_MEMORY     = %d
---CLEAN_INTERVAL = %d
-----------------------------------
`, d.usedSlots, d.maxSlots, d.usedMemory, d.maxMemory, d.clearInterval)
    fmt.Println(template)
}

func (d *Default) backgroundClean() {
	expiredKeys := make([]string, 0, d.maxSlots)
	for {
		// time to sleep
		time.Sleep(time.Duration(d.clearInterval) * time.Second)
		//fmt.Println("background goroutine running...")
		// got expired keys
		d.mu.RLock()
		for k, v := range d.bucket {
			if v.isExpired() {
				expiredKeys = append(expiredKeys, k)
			}
		}
		d.mu.RUnlock()
		fmt.Println("--------------")
		fmt.Printf("expired keys= %#v\n", expiredKeys)

		// do clean work, there should consider the state of too many keys
		for len(expiredKeys) > 0 {
			d.mu.Lock()
			for idx, _ := range expiredKeys {
				d.delete(expiredKeys[idx])
			}
			d.mu.Unlock()
		}
	}
}

func (d *Default) set(key string, val []byte, ttl time.Duration) (err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// check the number of slots
	if d.usedSlots > d.maxSlots {
		return fmt.Errorf("reached max slots")
	}
	fmt.Println("usedslots=", d.usedSlots, ", maxslots=", d.maxSlots, ", usedmemory=", d.usedMemory, " , max_memory=", d.maxMemory)

	// check the memory used
	var oldMem int
	if oldValue, ok := d.bucket[key]; ok {
		oldMem = oldValue.value.Len()
	}
	if d.usedMemory + len(val) - oldMem > d.maxMemory {
		return fmt.Errorf("reached max memory")
	}

	// update system info
	d.usedSlots += 1
	d.usedMemory += len(val) - oldMem
	buf := cachePool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Write(val)
	d.bucket[key] = localValue{
		ttl: time.Now().Add(ttl),
		value: buf,
	}

	return nil
}

func (d *Default)get(key string) (value []byte) {
	d.mu.RLock()
	v, ok := d.bucket[key]
	d.mu.RUnlock()

	if !ok {
		return []byte{}
	}

	// check if expired or not
	if v.isExpired() {
		d.delete(key)
		return []byte{}
	}
	return v.value.Bytes()
}

func (d *Default)delete(key string) {
	val, ok := d.bucket[key]
	if !ok {
		return
	}

	// update system info
	oldMem := val.value.Len()
	d.usedMemory -= oldMem
	d.usedSlots -= 1
	buf := val.value
	if buf != nil {
		buf.Reset()
		cachePool.Put(buf)
	}
	delete(d.bucket, key)
}