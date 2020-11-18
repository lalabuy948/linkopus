// +build emoji

package service

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/entity"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/go-redis/cache/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// Facade holds read, write and utility services to hide complexity.
type Facade struct {
	query      *Query
	command    *Command
	redisCache *Cache
}

// NewFacade returns instance of linkopus facade service.
func NewFacade(q *Query, c *Command, rc *Cache) *Facade {
	return &Facade{q, c, rc}
}

// HandleLinkMapInsert generates 6 digits hash according to given link.
// First of all this function will check if url is valid.
// Second check if this link already exists in database.
// Finally store it in links collection and cache it for 5 minutes.
func (f *Facade) HandleLinkMapInsert(link string) (*entity.LinkMap, error) {
	if !isValidUrl(link) {
		return nil, errors.New("service: invalid url")
	}

	result, err := f.query.QueryLinkMap(link, "")
	if err != mongo.ErrNoDocuments && result == nil {
		return nil, err
	}

	if result != nil {
		return nil, errors.New("service: mapping already exists")
		//return result, nil // dragons on line 47
	}

	linkHash := computeLinkHash(link)
	result, err = f.query.QueryLinkMap("", linkHash)
	if err != mongo.ErrNoDocuments {
		return f.HandleLinkMapInsert(link)
	}

	err = f.command.saveLinkMap(link, linkHash)
	if err != nil {
		return nil, errors.New("service: failed to save mapping")
	}

	linkMap := &entity.LinkMap{Link: link, LinkHash: linkHash}

	return linkMap, err
}

// HandleLinkExtraction extracting stored link according to given hash.
// First check if link exists in cache and return the result, otherwise
// check in links collection and store it in cache for 5 minutes.
func (f *Facade) HandleLinkExtraction(linkHash string) (*entity.LinkMap, error) {
	linkMap, err := f.redisCache.GetCachedLinkMap(linkHash)
	if err != cache.ErrCacheMiss && linkMap != nil {
		return linkMap, nil
	}

	linkMap, err = f.query.QueryLinkMap("", linkHash)
	if err != nil {
		return nil, err
	}

	if linkMap != nil {
		go f.redisCache.CacheLinkMap(*linkMap)
	}

	return linkMap, nil
}

// HandleLinkExtraction extracting stored link according to today's date.
// First check if stats exist in cache and return the result, otherwise
// check in views collection and store it in cache for 1 minute.
func (f *Facade) HandleTodayTopLinksViewsExtraction() (*[]entity.LinkView, error) {
	todayDate := time.Now().Format("2006-01-02")
	dateSplit := strings.Split(todayDate, "-")
	year := dateSplit[0]
	month := dateSplit[1]
	day := dateSplit[2]

	topLinkViews, err := f.redisCache.GetCachedTopViews(todayDate)
	if err != cache.ErrCacheMiss && topLinkViews != nil {
		return topLinkViews, nil
	}

	topLinkViews, err = f.query.QueryTodayTopLinksViews(year, month, day)
	if err != nil {
		return nil, err
	}

	if len(*topLinkViews) > 0 {
		go f.redisCache.CacheTopViews(todayDate, *topLinkViews)
	}

	return topLinkViews, err
}

// StopMessageProducer extracted in separate method not to expose producer to the Container.
func (f *Facade) StopMessageProducer() {
	f.command.producer.Stop()
}

// computeLinkHash generates SHA256 hash for given link and return first 6 chars.
// 6 chars should be enough for 56.800.235.584 amount of link maps.
// But.. emoji its more fun 🙆‍♀️
func computeLinkHash(link string) string {
	return gofakeit.Emoji() + gofakeit.Emoji() + gofakeit.Emoji() + gofakeit.Emoji()
}

// isValidUrl checks if given string is a valid url.
func isValidUrl(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	u, err := url.Parse(link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
