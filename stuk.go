// Copyright 2022-2022 The jdh99 Authors. All rights reserved.
// 带过期功能的map
// Authors: jdh99 <jdh821@163.com>
// 设计:准备map和过期链表.map节点中持有过期链表节点的地址,链表节点中有key.
// 当数据更新时,map节点直接将链表节点提到链表首部
// 开辟线程定时查询链表节点是否过期,过期则删除

package stuk

import (
	"container/list"
	"sync"
	"time"
)

type tItemMap struct {
	object   any
	itemList *list.Element
}

type tItemList struct {
	key uint64
	// 过期时间.单位:ns
	expire int64
}

// TimeoutCallback 超时回调函数
type TimeoutCallback func(k uint64, v any)

type tCache struct {
	items           map[uint64]*tItemMap
	timeList        *list.List
	lock            sync.RWMutex
	expire          time.Duration
	timeoutCallback TimeoutCallback
}

type Cache struct {
	*tCache
}

// New 新建map
func New(expire time.Duration, timeoutCallback TimeoutCallback) *Cache {
	var c tCache
	c.expire = expire
	c.items = make(map[uint64]*tItemMap)
	c.timeList = new(list.List)
	c.timeoutCallback = timeoutCallback

	var C Cache
	C.tCache = &c
	go checkExpire(&c)
	return &C
}

func checkExpire(c *tCache) {
	for {
		select {
		case <-time.After(time.Second):
		}

		c.lock.Lock()
		checkList(c)
		c.lock.Unlock()
	}
}

func checkList(c *tCache) {
	now := time.Now().UnixNano()
	node := c.timeList.Back()
	var nodeNext *list.Element
	var item *tItemList
	for {
		if node == nil {
			break
		}
		nodeNext = node.Next()

		item = node.Value.(*tItemList)
		if now > item.expire {
			var value *tItemMap = nil
			if c.timeoutCallback != nil {
				value, _ = c.items[item.key]
			}

			delete(c.items, item.key)
			c.timeList.Remove(node)
			if c.timeoutCallback != nil && value != nil {
				c.timeoutCallback(item.key, value.object)
			}
		}

		node = nodeNext
	}
}

// Set 设置键值对
func (c *tCache) Set(k uint64, v any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.items[k]
	if ok == false {
		var item tItemList
		item.key = k
		item.expire = time.Now().Add(c.expire).UnixNano()

		c.items[k] = &tItemMap{
			object:   v,
			itemList: c.timeList.PushFront(&item),
		}
	} else {
		value.object = v
		c.timeList.MoveToFront(value.itemList)
	}
}

// Get 读取键值对.返回nil表示读取失败
func (c *tCache) Get(k uint64) any {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.items[k]
	if ok == false {
		return nil
	}
	return value.object
}

// Pull 读取键值对并刷新过期时间.返回nil表示读取失败
func (c *tCache) Pull(k uint64) any {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.items[k]
	if ok == false {
		return nil
	}

	itemList := value.itemList.Value.(*tItemList)
	itemList.expire = time.Now().Add(c.expire).UnixNano()

	c.timeList.MoveToFront(value.itemList)
	return value.object
}

// Delete 删除键值对
func (c *tCache) Delete(k uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.items[k]
	if ok == false {
		return
	}

	c.timeList.Remove(value.itemList)
	delete(c.items, k)
}
