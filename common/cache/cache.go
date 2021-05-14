package cache

// Explicit lock free LRU cache

import (
	"container/list"
	"context"
)


type LRUCache struct {
	inCh chan *inData
	size int
	mp map[string]*list.Element
	queue *list.List
}


type cacheData struct {
	key string
	value interface{}
}

type inData struct {
	outCh chan *cacheData
	fn func(chan<- *cacheData)
}

func (cache *LRUCache) processor(ctx context.Context) {
	for {
		select {
		case ind := <- cache.inCh:
			ind.fn(ind.outCh)
		case <- ctx.Done():
			break
		}
	}
}

func NewLRUCache(size int) *LRUCache {
	cache := &LRUCache{
		inCh : make(chan *inData, 5),
		size: size,
		mp : map[string]*list.Element{},
		queue : list.New(),
	}
	go cache.processor(context.TODO())
	return cache
}

func (cache *LRUCache) Set(key string, value interface{}) {
	fn := func(outCh chan<- *cacheData) {
		for cache.queue.Len() >= cache.size {
			elem := cache.queue.Front()
			cache.queue.Remove(elem)
			data := elem.Value.(*cacheData)
			delete (cache.mp, data.key)
		}
		data := &cacheData{key: key, value: value}
		cache.queue.PushBack(data)
		cache.mp[key] = cache.queue.Back()
		outCh <- data
	}
	outCh := make(chan *cacheData)
	defer close(outCh)
	cache.inCh <- &inData{fn : fn, outCh: outCh}
	<- outCh
}

func (cache *LRUCache) Get(key string) (interface{}, bool) {
	fn := func(outCh chan<- *cacheData) {
		elem, ok := cache.mp[key]
		if !ok {
			outCh <- nil
			return
		}
		cache.queue.Remove(elem)
		data := elem.Value.(*cacheData)
		cache.queue.PushBack(data)
		cache.mp[key] = cache.queue.Back()
		outCh <- data
	}
	outCh := make(chan *cacheData)
	defer close(outCh)
	cache.inCh <- &inData{fn : fn, outCh: outCh}
	data := <- outCh
	return data, data != nil
}

func (cache *LRUCache) Delete(key string) {
	fn := func(outCh chan<- *cacheData) {
		elem, ok := cache.mp[key]
		if !ok {
			outCh <- nil
			return
		}
		cache.queue.Remove(elem)
		data := elem.Value.(*cacheData)
		delete(cache.mp, key)
		outCh <- data
	}
	outCh := make(chan *cacheData)
	defer close(outCh)
	cache.inCh <- &inData{fn : fn, outCh: outCh}
	<- outCh
}