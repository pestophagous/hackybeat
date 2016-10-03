package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/pestophagous/hackybeat/beater"
)

func main() {
	err := beat.Run("hackybeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
