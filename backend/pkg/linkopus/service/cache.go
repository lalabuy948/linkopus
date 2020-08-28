// cache.go provides type safe redis cache functionality.

package service

import (
	"context"
	"time"

	"github.com/go-redis/cache/v8"
)

// Cache struct hold cache service.
type Cache struct {
	redisCache *cache.Cache
}

// NewCache constructs cache service, takes cache.Cache as first function argument and return Cache struct.
func NewCache(rc *cache.Cache) *Cache {
	return &Cache{rc}
}

var ctx = context.Background()

// CacheLinkMap stores given LinkMap using LinkMap.LinkHash as a key.
func (c *Cache) CacheLinkMap(linkMap LinkMap) error {
	err := c.redisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   linkMap.LinkHash,
		Value: linkMap,
		TTL:   1 * time.Minute,
	})

	if err != nil {
		return err
	}

	return nil
}

// CacheLinkMap extracts LinkMap by given LinkMap.LinkHash as first function argument.
func (c *Cache) GetCachedLinkMap(key string) (*LinkMap, error) {
	var linkMap LinkMap
	if err := c.redisCache.Get(ctx, key, &linkMap); err != nil {
		return nil, err
	}

	return &linkMap, nil
}

// CacheTopViews stores given [] LinkView using today's date (2006-01-02) as a key.
func (c *Cache) CacheTopViews(todayDate string, linkViews []LinkView) error {
	err := c.redisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   todayDate,
		Value: linkViews,
		TTL:   1 * time.Minute,
	})

	if err != nil {
		return err
	}

	return nil
}

// GetCachedTopViews extracts [] LinkView by given date (2006-01-02) as first function argument.
func (c *Cache) GetCachedTopViews(todayDate string) (*[]LinkView, error) {
	var linkViews []LinkView
	if err := c.redisCache.Get(ctx, todayDate, &linkViews); err != nil {
		return nil, err
	}

	return &linkViews, nil
}
