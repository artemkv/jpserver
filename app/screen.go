package app

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"

	dda "github.com/shinkar94/godesktopdup"
	xdraw "golang.org/x/image/draw"
)

const PRIMARY_MONITOR = 0
const FRAME_CAPTURE_TIMEOUT = 1000
const RESIZING_FACTOR = 2

type Capturer struct {
	dd     *dda.DesktopDuplication
	width  int
	height int
	buffer []byte
	img    *image.RGBA
}

func NewCapturer() (*Capturer, error) {
	dd, err := dda.New(PRIMARY_MONITOR)
	if err != nil {
		return nil, fmt.Errorf("failed to create desktop duplication: %w", err)
	}

	width, height, err := dd.GetSize()
	if err != nil {
		dd.Release()
		return nil, fmt.Errorf("failed to detect screen size: %w", err)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	return &Capturer{
		dd:     dd,
		width:  width,
		height: height,
		buffer: make([]byte, width*height*4),
		img:    img,
	}, nil
}

func (c *Capturer) Close() {
	if c.dd != nil {
		c.dd.Release()
		c.dd = nil
	}
}

func (c *Capturer) CaptureFrame() error {
	for {
		err := c.dd.GetFrameBGRA(c.buffer, FRAME_CAPTURE_TIMEOUT)
		if err != nil {
			if err.Error() == "no image yet" {
				continue
			}
			return fmt.Errorf("failed to capture frame: %w", err)
		}
		return nil
	}
}

func (c *Capturer) ConvertToPNG() ([]byte, error) {
	for i := 0; i < len(c.buffer); i += 4 {
		c.img.Pix[i+0] = c.buffer[i+2] // R
		c.img.Pix[i+1] = c.buffer[i+1] // G
		c.img.Pix[i+2] = c.buffer[i+0] // B
		c.img.Pix[i+3] = 255
	}

	// TODO: resize
	resizedRect := image.Rect(0, 0, c.width/RESIZING_FACTOR, c.height/RESIZING_FACTOR)
	resizedImg := image.NewRGBA(resizedRect)
	xdraw.NearestNeighbor.Scale(resizedImg, resizedRect, c.img, c.img.Bounds(), draw.Over, nil)

	var out bytes.Buffer
	if err := png.Encode(&out, c.img); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return out.Bytes(), nil
}
