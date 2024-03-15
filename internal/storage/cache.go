package storage

import (
	"container/list"
	"sync"

	"donation-service/internal/data"
)

type Cache struct {
	mu        sync.Mutex
	Donations *list.List
}

func NewCache() *Cache {
	return &Cache{
		Donations: list.New(),
	}
}

func (c *Cache) PushDonation(d data.Donation) {
	c.mu.Lock()

	c.Donations.PushBack(d)

	c.mu.Unlock()
}

func (c *Cache) PopDonation() data.Donation {
	var d data.Donation

	c.mu.Lock()

	if c.Donations.Len() > 0 {
		e := c.Donations.Front()
		c.Donations.Remove(e)

		d = e.Value.(data.Donation)
	}

	c.mu.Unlock()

	return d
}

func (c *Cache) FlushAll() []data.Donation {
	donations := make([]data.Donation, c.Donations.Len())

	c.mu.Lock()

	for c.Donations.Len() > 0 {
		e := c.Donations.Front()
		c.Donations.Remove(e)

		d := e.Value.(data.Donation)
		donations = append(donations, d)
	}

	c.mu.Unlock()

	return donations
}
