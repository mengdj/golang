package main_test

import (
	"log"
	"screenshot"
	"testing"
)

func TestCapture(t *testing.T) {
	shot := screenshot.NewScreenshotDefault()
	if nil == shot.Capture() {
		log.Println(shot.Resize(1920, 0, screenshot.Lanczos3, 90))
	}
}
