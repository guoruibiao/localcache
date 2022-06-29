# localcache
localcache in golang


1. customize ur conf
```toml
# customize your conf name
name="LOCAL_CACHE"
# maximum slots or maximum keys you want
max_slots=10
# maximum memory you want
max_memory=1024

# choose one specific strategy, such as default/lru/lfu/fifo
cache_strategy="default"

# clean interval, and the unit is `time.Second`
clean_interval=1
```

cache_strategy contains `default`, `lru`, `fifo`, `lfu`, ... select one which you need.

2 init the instance
```golang
package main

import (
    "time"
	
    "github.com/guoruibiao/localcache"
)

func main() {
    confPath := "./conf/localcache.toml"
    localcache.Init(confPath)
    
    localcache.Set("name", []byte("tiger"), time.Second*2)
    println(localcache.GetString("name"))
    time.Sleep(time.Second*3)
    println(localcache.GetString("name"))
}
```