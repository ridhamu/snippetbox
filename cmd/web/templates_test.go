package main

import (
	"testing"
	"time"

	"github.com/ridhamu/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	testTable := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2026, 0o6, 14, 21, 0, 0, 0, time.UTC),
			want: "14 Jun 2026 at 21:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2026, 0o6, 14, 22, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "14 Jun 2026 at 21:00",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			humanDateResult := humanDate(test.tm)

			assert.Equal(t, humanDateResult, test.want)
		})
	}
}
