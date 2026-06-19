//go:build !solution

package lrucache

import (
	"container/list"
)

type pair struct {
	value int
	pos   *list.Element
}

type LRUCache struct {
	cap    int
	values map[int]pair
	order  *list.List
}

func (c *LRUCache) Get(key int) (int, bool) {
	v, ok := c.values[key]
	if !ok {
		return 0, false
	}
	c.order.MoveToFront(v.pos)
	return v.value, true
}

func (c *LRUCache) Set(key, value int) {
	if c.cap == 0 {
		return
	}
	v, ok := c.values[key]
	if ok {
		c.values[key] = pair{value: value, pos: v.pos}
		c.order.MoveToFront(v.pos)
		return
	}
	if len(c.values) == c.cap {
		last := c.order.Back()
		delete(c.values, last.Value.(int))
		c.order.Remove(last)
	}
	pos := c.order.PushFront(key)
	c.values[key] = pair{value: value, pos: pos}
}

func (c *LRUCache) Range(f func(key, value int) bool) {
	for cur := c.order.Back(); cur != nil; {
		key := cur.Value.(int)
		if !f(key, c.values[key].value) {
			return
		}
		cur = cur.Prev()
	}
}

func (c *LRUCache) Clear() {
	for k := range c.values {
		delete(c.values, k)
	}
	for cur := c.order.Front(); cur != nil; {
		next := cur.Next()
		c.order.Remove(cur)
		cur = next
	}
}

func New(cap int) Cache {
	if cap < 0 {
		panic("cache size can't be less then 0")
	}
	return &LRUCache{
		cap:    cap,
		values: make(map[int]pair, cap),
		order:  list.New(),
	}
}
