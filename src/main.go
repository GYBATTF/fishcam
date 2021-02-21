package main

import (
	"log"
	"time"
)

func main() {
	service()
}

// service sets a photo as the background at a regular interval specified in the config,
// printing an error to log if any occur.
func service() {
	set := func() {
		err := setBackground(getPicDir(), func(filename string) {
			log.Printf("Downloaded '%s' at '%s'", filename, time.Now())
		})

		if err != nil {
			log.Println(err)
		}
	}

	durationChannel := make(chan time.Duration)
	go readConfig(durationChannel, func(err error) {
		if err != nil {
			log.Println(err)
		}
	})

	var duration *time.Duration = nil
	for {
		if duration == nil {
			go set()
			d := <-durationChannel
			duration = &d
		}

		select {
		case <-time.After(*duration):
			go set()
		case d := <-durationChannel:
			duration = &d
		}
	}
}
