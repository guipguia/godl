package util

import (
	"net/url"
	"path"
)

// GetNameBasedOnUrl will return the base name of the given URL.
func GetNameBasedOnUrl(fullURL string) string {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		panic(err)
	}
	return path.Base(parsedURL.Path)
}
