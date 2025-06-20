package multi_tier

import (
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	length time.Duration
}

type MultiTier struct {
	primary, secondary, third string
	tracks                    []*Track
}

func (m *MultiTier) Len() int {
	return len(m.tracks)
}

func (m *MultiTier) Swap(i, j int) {
	m.tracks[i], m.tracks[j] = m.tracks[j], m.tracks[i]
}

func (m *MultiTier) Less(i, j int) bool {
	itemI, itemJ := m.tracks[i], m.tracks[j]
	for _, f := range []string{
		m.primary, m.secondary, m.third,
		"title", "artist", "album", "year", "length"} {
		switch f {
		case "title":
			if itemI.Title != itemJ.Title {
				return itemI.Title < itemJ.Title
			}
		case "artist":
			if itemI.Artist != itemJ.Artist {
				return itemI.Artist < itemJ.Artist
			}
		case "album":
			if itemI.Album != itemJ.Album {
				return itemI.Album < itemJ.Album
			}
		case "year":
			if itemI.Year != itemJ.Year {
				return itemI.Year < itemJ.Year
			}
		case "length":
			if itemI.length != itemJ.length {
				return itemI.length < itemJ.length
			}
		}
	}

	return true
}

func setPrimary(x *MultiTier, p string) {
	x.primary, x.secondary, x.third = p, x.primary, x.secondary
}

func SetPrimary(x sort.Interface, p string) {
	if x, ok := x.(*MultiTier); ok {
		setPrimary(x, p)
	}
}

func NewMultiTier(t []*Track, p, s, th string) sort.Interface {
	return &MultiTier{
		primary:   p,
		secondary: s,
		third:     th,
		tracks:    t,
	}
}
