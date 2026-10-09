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

type ImageCache struct {
	rect image.Rectangle
	img  *image.RGBA
}

type Capturer struct {
	dd     *dda.DesktopDuplication
	width  int
	height int
	buffer []byte
	// re-using stuff across calls to avoid re-allocating
	img       *image.RGBA
	resizeImg ImageCache
	pngBuffer bytes.Buffer
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

	// to be re-used
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	imgCache := ImageCache{
		rect: image.Rect(0, 0, 0, 0),
		img:  nil,
	}

	return &Capturer{
		dd:        dd,
		width:     width,
		height:    height,
		buffer:    make([]byte, width*height*4),
		img:       img,
		resizeImg: imgCache,
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
	// Calculate clipping area
	x0cut := c.width * borders.left / 100
	y0cut := c.height * borders.top / 100
	x1cut := c.width - c.width*borders.right/100
	y1cut := c.height - c.height*borders.bottom/100
	cutoutRect := image.Rect(x0cut, y0cut, x1cut, y1cut)

	// Copy the buffer into image, converting BGRA -> RGBA
	// Only copy the cutout part
	for y := y0cut; y < y1cut; y++ {
		for x := x0cut; x < x1cut; x++ {
			i := (y*c.width + x) * 4
			c.img.Pix[i+0] = c.buffer[i+2] // R
			c.img.Pix[i+1] = c.buffer[i+1] // G
			c.img.Pix[i+2] = c.buffer[i+0] // B
			c.img.Pix[i+3] = 255           // A
		}
	}

	// Target image rectangle
	resizedRect := image.Rect(0, 0,
		cutoutRect.Dx()/resizeFactor, cutoutRect.Dy()/resizeFactor)

	// Resize target
	var resizedImg *image.RGBA
	if c.resizeImg.img != nil &&
		c.resizeImg.rect == resizedRect {
		// from cache
		resizedImg = c.resizeImg.img
	} else {
		resizedImg = image.NewRGBA(resizedRect)
		// cache
		c.resizeImg.rect = resizedRect
		c.resizeImg.img = resizedImg
	}

	// Resize
	xdraw.NearestNeighbor.Scale(
		resizedImg, resizedRect, c.img, cutoutRect, draw.Src, nil)

	// PNG encode the result
	c.pngBuffer.Reset()
	if err := png.Encode(&c.pngBuffer, resizedImg); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return c.pngBuffer.Bytes(), nil
}
