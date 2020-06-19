package main

import (
    "testing"
    "time"
)

func TestHumanDate(t *testing.T) {
    tests := []struct {
        name string
        time time.Time
        want string
    }{
        {
            name: "UTC",
            time: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
            want: "17 Dec 2020 at 10:00",
        },
        {
            name: "Empty",
            time: time.Time{},
            want: "",
        },
        {
            name: "CET",
            time: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
            want: "17 Dec 2020 at 09:00",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := humanDate(tt.time)

            if got != tt.want {
                t.Errorf("want %q; got %q", tt.want, got)
            }
        })
    }
}
