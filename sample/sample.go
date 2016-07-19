package main

import (
	"github.com/tompng/go-ascii-canvas"
	"math/rand"
	"time"
)

var fontImage *asciicanvas.ImageBuffer

func DrawText(canvas *asciicanvas.ImageBuffer, text string, x, y, size float64) {
	if fontImage == nil {
		fontImage, _ = asciicanvas.NewImageBufferFromFile("font.png")
	}
	for i, code := range text {
		charImage := fontImage.Sub(float64(code%16)/16.0, float64(code/16)/8.0, 1/16.0, 1/8.0)
		canvas.Draw(charImage, x+float64(i)*size/2, y, size/2, size)
	}
}

func main() {
	img, err := asciicanvas.NewImageBufferFromFile("gopher.png")
	if err != nil {
		panic("cannot read file")
	}
	for {
		terminalWidth, terminalHeight := asciicanvas.GetWinSize()
		canvas := asciicanvas.NewImageBuffer(terminalWidth, 2*terminalHeight)
		canvas.RotateDraw(img, 0, 0, 80, 80, -30+60*rand.Float64())
		for x := 0; x < canvas.Width; x++ {
			for y := 0; y < canvas.Height; y++ {
				canvas.Plot(x, y, float64(x)/float64(canvas.Width), 0.5*float64(y)/float64(canvas.Height))
			}
		}
		DrawText(canvas, "Gopher!", 4*rand.Float64(), 4*rand.Float64(), 20)
		canvas.Print()
		time.Sleep(50 * time.Millisecond)
	}
}
