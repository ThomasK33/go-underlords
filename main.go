package main

import (
	"log"
)

var (
	// Version - Version number
	Version string
	// Build - Build number
	Build string
)

func main() {
	log.Println("Starting " + Version + " (" + Build + ")")
}
