# localcache
localcache in golang


1. customize ur conf
```toml
name="LOCAL_CACHE"
max_slots=10
# byte
max_memory=1024

cache_strategy="default"
clean_interval=1
```

cache_strategy contains `default`, `lru`, `fifo`, `lfu`, ... select one you need.

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