package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	mu       sync.Mutex // guards capturer
	capturer *Capturer
)

func getCapturer() (*Capturer, error) {
	mu.Lock()
	defer mu.Unlock()

	if capturer == nil {
		log.Print("Initializing new capturer")

		c, err := NewCapturer()
		if err != nil {
			return nil, err
		}
		capturer = c
	}

	return capturer, nil
}

func releaseCapturer() {
	mu.Lock()
	defer mu.Unlock()

	if capturer == nil {
		return
	}

	capturer.Close()
	capturer = nil
}

func handleFrame(c *gin.Context) {
	borders := Borders{
		left:   getQueryParamAsInt(c, "left", 0),
		right:  getQueryParamAsInt(c, "right", 0),
		top:    getQueryParamAsInt(c, "top", 0),
		bottom: getQueryParamAsInt(c, "bottom", 0),
	}
	resizeFactor := getQueryParamAsInt(c, "resize", 1)

	cp, err := getCapturer()
	if err != nil {
		fmt.Printf("Could not initialize capturer: %v", err)
		toInternalServerError(c, fmt.Sprintf("Could not initialize capturer: %v", err))
		return
	}

	err = cp.CaptureFrame()
	if err != nil {
		releaseCapturer()
		fmt.Printf("Failed to capture the frame: %v", err)
		toInternalServerError(c, fmt.Sprintf("Failed to capture the frame: %v", err))
		return
	}

	png, err := cp.ConvertToPNG(borders, resizeFactor)
	if err != nil {
		fmt.Printf("Failed to convert to png: %v", err)
		toInternalServerError(c, fmt.Sprintf("Failed to convert to png: %v", err))
		return
	}

	toBinary(c, png)
}
