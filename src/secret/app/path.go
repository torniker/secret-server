package app

import (
	"fmt"
	"net/url"
	"secret/log"
	"strings"
)

// Path defines structure for routing
type Path struct {
	url      *url.URL
	segments []string
	index    int
}

// NewPath creates path object from url
func NewPath(url *url.URL) *Path {
	return &Path{
		url:      url,
		segments: strings.Split(url.Path, "/"),
		index:    0,
	}
}

// URL returns url
func (p *Path) URL() *url.URL {
	return p.url
}

// Next returns next segment of the path
func (p *Path) Next() string {
	if len(p.segments) < p.index+2 {
		return ""
	}
	log.With("route", p.segments[p.index+1]).With("index", p.index+1).With("segments", p.segments).Info("next route")
	return p.segments[p.index+1]
}

// Current returns current segment of the path
func (p *Path) Current() string {
	if len(p.segments) < p.index+1 {
		return ""
	}
	return p.segments[p.index]
}

// Increment increases path index by 1
func (p *Path) Increment() {
	p.index++
}

func (p Path) String() string {
	return fmt.Sprintf("url schema: %s, host: %s, path: %s, index: %v", p.url.Scheme, p.url.Host, p.url.Path, p.index)
}
