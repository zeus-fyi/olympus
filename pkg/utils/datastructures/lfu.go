package datastructures

import (
	"container/list"
	"sync"
)

type Eviction struct {
	Key   string
	Value interface{}
}

type LFU struct {
	// If len > UpperBound, cache will automatically evict
	// down to LowerBound.  If either value is 0, this behavior
	// is disabled.
	UpperBound      int
	LowerBound      int
	values          map[string]*cacheEntry
	freqs           *list.List
	len             int
	Lock            *sync.Mutex
	EvictionChannel chan<- Eviction
}

type cacheEntry struct {
	key      string
	value    interface{}
	freqNode *list.Element
}

type listEntry struct {
	entries map[*cacheEntry]byte
	freq    int
}

func New() *LFU {
	c := new(LFU)
	c.values = make(map[string]*cacheEntry)
	c.freqs = list.New()
	c.Lock = new(sync.Mutex)
	return c
}

func (c *LFU) Get(key string) interface{} {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value
	}
	return nil
}

func (c *LFU) Set(key string, value interface{}) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if e, ok := c.values[key]; ok {
		// value already exists for key.  overwrite
		e.value = value
		c.increment(e)
	} else {
		// value doesn't exist.  insert
		e := new(cacheEntry)
		e.key = key
		e.value = value
		c.values[key] = e
		c.increment(e)
		c.len++
		// bounds mgmt
		if c.UpperBound > 0 && c.LowerBound > 0 {
			if c.len > c.UpperBound {
				c.evict(c.len - c.LowerBound)
			}
		}
	}
}

func (c *LFU) Len() int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.len
}

func (c *LFU) Evict(count int) int {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	return c.evict(count)
}

func (c *LFU) GetLeastFrequentValue() (interface{}, interface{}) {
	if place := c.freqs.Front(); place != nil {
		for entry, _ := range place.Value.(*listEntry).entries {
			return entry.key, entry.value
		}
	}
	return nil, nil
}

func (c *LFU) evict(count int) int {
	// No Lock here so it can be called
	// from within the Lock (during Set)
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry, _ := range place.Value.(*listEntry).entries {
				if i < count {
					if c.EvictionChannel != nil {
						c.EvictionChannel <- Eviction{
							Key:   entry.key,
							Value: entry.value,
						}
					}
					delete(c.values, entry.key)
					c.remEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *LFU) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		// new entry
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		// move up
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next()
	}

	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		// create a new list entry
		li := new(listEntry)
		li.freq = nextFreq
		li.entries = make(map[*cacheEntry]byte)
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry).entries[e] = 1
	if currentPlace != nil {
		// remove from current position
		c.remEntry(currentPlace, e)
	}
}

func (c *LFU) remEntry(place *list.Element, entry *cacheEntry) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}
