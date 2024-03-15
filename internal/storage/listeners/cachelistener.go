package listeners

import (
	"donation-service/internal/data"
	"donation-service/internal/storage"
)

type CacheListener struct {
	Cache *storage.Cache
}

func NewCacheListener(cache *storage.Cache) CacheListener {
	return CacheListener{
		Cache: cache,
	}
}

func (c CacheListener) NotifyNewDonation(donation data.Donation) {
	c.Cache.PushDonation(donation)
}
