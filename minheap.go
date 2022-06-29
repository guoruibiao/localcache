package localcache

type MinHeap []*localValue

func (mh MinHeap) Len() int {
	return len(mh)
}

func (mh MinHeap) Less(i, j int) bool {
	// first judge `frequency``
	if mh[j].frequency != mh[i].frequency {
		return mh[j].frequency > mh[i].frequency
	}

	// judge by character order with `keyname`
	if mh[j].keyname != mh[i].keyname {
		return mh[j].keyname < mh[i].keyname
	}

	// finally judge the `ttl` of expired time
	return mh[j].ttl.After(mh[i].ttl)
}

func (mh MinHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
}

func (mh *MinHeap) Push(lv interface{}) {
	*mh = append(*mh, lv.(*localValue))
}

func (mh *MinHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	x := old[n-1]
	*mh = old[0 : n-1]
	return x
}

func (mh *MinHeap) Top() interface{} {
	if mh.Len() <= 0 {
		return nil
	}
	return (*mh)[0]
}

func (mh *MinHeap) Index(ele *localValue) int {
	for idx, item := range *mh {
		if item.keyname == ele.keyname {
			return idx
		}
	}
	return -1
}
