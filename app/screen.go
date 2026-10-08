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

// In % of the screen
type Borders struct {
	left   int
	right  int
	top    int
	bottom int
}

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

func (c *Capturer) ConvertToPNG(borders Borders, resizeFactor int) ([]byte, error) {
	// Copy the buffer into image
	for i := 0; i < len(c.buffer); i += 4 {
		c.img.Pix[i+0] = c.buffer[i+2] // R
		c.img.Pix[i+1] = c.buffer[i+1] // G
		c.img.Pix[i+2] = c.buffer[i+0] // B
		c.img.Pix[i+3] = 255
	}

	// Clipping area
	srcWidth := c.img.Bounds().Dx()
	srcHeight := c.img.Bounds().Dy()
	cutoutRect := image.Rect(
		srcWidth*borders.left/100,
		srcHeight*borders.top/100,
		srcWidth-srcWidth*borders.right/100,
		srcHeight-srcHeight*borders.bottom/100)

	// Target image
	resizedRect := image.Rect(
		0,
		0,
		cutoutRect.Dx()/resizeFactor,
		cutoutRect.Dy()/resizeFactor)

	// Resize
	resizedImg := image.NewRGBA(resizedRect)
	xdraw.NearestNeighbor.Scale(
		resizedImg, resizedRect, c.img, cutoutRect, draw.Over, nil)

	// PNG encode the result
	var out bytes.Buffer
	if err := png.Encode(&out, resizedImg); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return out.Bytes(), nil
}
