// +build emoji

package service

import "github.com/brianvoe/gofakeit/v5"

// computeLinkHash generates SHA256 hash for given link and return first 6 chars.
// 6 chars should be enough for 56.800.235.584 amount of link maps.
// But.. emoji its more fun 🙆‍♀️
func computeLinkHash(link string) string {
	return gofakeit.Emoji() + gofakeit.Emoji() + gofakeit.Emoji() + gofakeit.Emoji()
}
