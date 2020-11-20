package main_test

import (
	"log"
	"screenshot"
	"testing"
)

func TestCapture(t *testing.T) {
	shot := screenshot.NewScreenshotDefault()
	log.Println(shot.Capture())
}
