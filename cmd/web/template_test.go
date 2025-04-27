package main

import (
	"testing"
	"time"

	"dream.website/internal/assert"
)

func TestHumandate(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 March 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 March 2022 at 09:15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := HumanDate(test.tm)

			assert.Equal(t, test.want, got)

		})
	}

}
