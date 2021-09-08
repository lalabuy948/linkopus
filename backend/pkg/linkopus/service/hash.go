// +build !emoji

package service

import (
	"crypto/sha256"
	"encoding/base64"
)

// computeLinkHash generates SHA256 hash for given link and return first 6 chars.
// 6 chars should be enough for 56.800.235.584 amount of link maps.
// But.. emoji its more fun 🙆‍♀️ facade_emoji.go
func computeLinkHash(link string) string {
	hashFunc := sha256.New()
	hashFunc.Write([]byte(link)) //nolint
	sha := base64.URLEncoding.EncodeToString(hashFunc.Sum(nil))

	return sha[1:7]
}
