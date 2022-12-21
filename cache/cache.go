package cache

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var resp sync.Map
var ticker *time.Ticker

var lifecycle TimeLimit

type Cached struct {
	Data              interface{}
	hash              int
	ExpireAtTimestamp time.Time
}

type TimeLimit struct {
	Seconds int
	Minutes int
	Hours   int
}

// InitCache
// cacheCleanUpCycle how often it will check if items need to be removed from map
// timeToStayInMemory how long you believe items show stay in memory
func InitCache(cacheCleanUpCycle *TimeLimit, timeToStayInMemory TimeLimit) {
	lifecycle = timeToStayInMemory
	resp = sync.Map{}
	if timeToStayInMemory.Seconds <= 0 && timeToStayInMemory.Minutes <= 0 && timeToStayInMemory.Hours <= 0 {
		log.Fatal("C")
		return
	}
	if cacheCleanUpCycle != nil {
		go schedule(cacheCleanUpCycle)
	}
}

// Add
// This will add to the internal cache
// key us what you want item to be looked up by
// value is the item you wish to store
// It will use the preset timeToStayInMemory to set how long the item should stay in memory
func Add(key string, value interface{}) {
	resp.Store(key, Cached{
		Data:              value,
		ExpireAtTimestamp: time.Now().Add(makeDuration(&lifecycle)),
		hash:              rand.Int(),
	})
}

// AddWithDifferentTimeLimit
// This will add to the internal cache
// key us what you want item to be looked up by
// value is the item you wish to store
// It will use the limit to set how long the item should stay in memory
func AddWithDifferentTimeLimit(key string, value interface{}, limit TimeLimit) {
	resp.Store(key, Cached{
		Data:              value,
		ExpireAtTimestamp: time.Now().Add(makeDuration(&limit)),
		hash:              rand.Int(),
	})
}

// UpdateCacheTimeLimit
// This will update how long all items will be stored in cache
func UpdateCacheTimeLimit(limit TimeLimit) {
	if limit.Seconds <= 0 && limit.Minutes <= 0 && limit.Hours <= 0 {
		lifecycle = limit
	}
}

// Get
// This will collect active cached item
// Will return if we have item in cache or will return nil if no item was found
// The bool will return if key is in cache
// Will removed and return nil is item has expired in cache
func Get(key string) (interface{}, bool) {
	var cache Cached
	now := time.Now()
	val, ok := resp.Load(key)
	if !ok {
		return nil, false
	}
	cache = val.(Cached)

	log.Println(now)
	if cache.ExpireAtTimestamp.Before(time.Now()) {
		remove(cache, key)
		return nil, false
	}
	return cache.Data, ok
}

func remove(cache Cached, key string) {
	val, _ := resp.Load(key)
	hash := val.(Cached).hash
	if cache.hash == hash {
		resp.Delete(key)
	}
}

func makeDuration(timeLimit *TimeLimit) time.Duration {
	return time.Duration(timeLimit.Hours)*time.Hour +
		time.Duration(timeLimit.Minutes)*time.Minute +
		time.Duration(timeLimit.Seconds)*time.Second
}

func schedule(purge *TimeLimit) {
	wait := makeDuration(purge)
	ticker := time.NewTicker(wait)

	for _ = range ticker.C {
		now := time.Now()
		resp.Range(func(key, value interface{}) bool {
			if value.(Cached).ExpireAtTimestamp.Before(now) {
				resp.Delete(key)
			}
			return true
		})
	}

}
