package handlers

import "time"

type ShortenedUrls struct {
	Id        int
	Longurl   string
	Shorturl  string
	ExpiresAt *time.Time
}
