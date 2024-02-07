package main

import (
	"testing"
	"time"

	"github.com/Zekeriyyah/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// tm := time.Date(2024, 02, 1, 1, 35, 0, 0, time.UTC)
	// hd := humanDate(tm)

	// if hd != "01 Feb 2024 at 01:35" {
	// 	t.Errorf("got %q: want %q", hd, "01 Feb 2024 at 01:35")
	// }

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 2, 2, 4, 8, 0, 0, time.UTC),
			want: "02 Feb 2024 at 04:08",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 2, 2, 4, 8, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "02 Feb 2024 at 03:08",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			assert.Equal(t, hd, tt.want) // helper func in internal/assert
		})
	}
}
