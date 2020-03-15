package main

import (
	"errors"
	"github.com/ian-kent/go-log/log"
	"time"
)

// Release A single instance of a release
type Release struct {
	ID            int
	Name          string
	Version       string
	DateSubmitted time.Time
	Released      bool
	DateReleased  time.Time
}

// Releases The releases contained by the service
var Releases map[int]Release

// MarkReleased Marks the release as released this second
func (r *Release) MarkReleased() {
	if r.Released {
		return
	}
	r.Released = true
	r.DateReleased = time.Now()
}

// AddRelease Adds a release to the todo list
func AddRelease(r Release) error {
	if Releases == nil {
		Releases = make(map[int]Release)
	}
	_, ok := Releases[r.ID]
	if ok {
		log.Error("ID %d already exists in the system!", r.ID)
		return errors.New("ID already exists in the system")
	}
	Releases[r.ID] = r
	return nil
}
