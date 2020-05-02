package cache

import (
	"time"

	"github.com/alka/supermart/store/model"
	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(5*time.Minute, 5*time.Minute)

func SetCache(key string, item []*model.Item) bool {
	Cache.Set(key, item, cache.NoExpiration)
	return true
}

func GetCache(key string) []*model.Item {
	var items []*model.Item
	data, found := Cache.Get(key)
	if found {
		for _, v := range data.([]*model.Item) {
			items = append(items, v)
		}
	}
	return items
}
