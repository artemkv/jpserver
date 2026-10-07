package app

import (
	"fmt"
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
	cp, err := getCapturer()
	if err != nil {
		fmt.Printf("Could not initialize capturer: %v", err)
		toInternalServerError(c, fmt.Sprintf("Could not initialize capturer: %v", err))
	}

	err = cp.CaptureFrame()
	if err != nil {
		releaseCapturer()
		fmt.Printf("Failed to capture the frame: %v", err)
		toInternalServerError(c, fmt.Sprintf("Failed to capture the frame: %v", err))
	}

	png, err := cp.ConvertToPNG()
	if err != nil {
		fmt.Printf("Failed to convert to png: %v", err)
		toInternalServerError(c, fmt.Sprintf("Failed to convert to png: %v", err))
	}

	toBinary(c, png)
}
