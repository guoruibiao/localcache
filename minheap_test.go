package localcache

import (
	"bytes"
	"container/heap"
	"testing"
)

func TestMinHeap(t *testing.T) {
	mh := &MinHeap{}
	heap.Init(mh)
	heap.Push(mh, &localValue{
		keyname:   "lv2",
		value:     bytes.NewBuffer([]byte("tiger2")),
		frequency: 2,
	})
	heap.Push(mh, &localValue{
		keyname:   "lv3",
		value:     bytes.NewBuffer([]byte("tiger3")),
		frequency: 3,
	})
	heap.Push(mh, &localValue{
		keyname:   "lv1",
		value:     bytes.NewBuffer([]byte("tiger1")),
		frequency: 1,
	})
	heap.Push(mh, &localValue{
		keyname:   "lv-1",
		value:     bytes.NewBuffer([]byte("tiger-1")),
		frequency: -1,
	})

	t.Logf("data=%#v\n", mh.Pop())
	t.Logf("data=%#v\n", mh.Pop())
	t.Logf("data=%#v\n", mh.Pop())
	t.Logf("data=%#v\n", mh.Pop())
}
