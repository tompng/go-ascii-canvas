package asciicanvas

import (
	"bytes"
	"image"
	_ "image/png"
	"io"
	"math"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func GetWinSize() (int, int) {
	type WinSize struct {
		H, W, _, _ int16
	}
	winsize := &WinSize{}
	syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(winsize)))
	return int(winsize.W), int(winsize.H)
}

type Image interface {
	ColorAt(x, y float64) (gray, alpha float64)
}

type ImageBuffer struct {
	Width  int
	Height int
	Gray   [][]float64
	Alpha  [][]float64
}

type SubImage struct {
	Source     Image
	X, Y, W, H float64
}

func (image *SubImage) ColorAt(x, y float64) (float64, float64) {
	if x < 0 || x > 1 || y < 0 || y > 1 {
		return 0, 0
	}
	return image.Source.ColorAt(image.X+image.W*x, image.Y+image.H*y)
}

func NewImageBufferFromBytes(b []byte) (*ImageBuffer, error) {
	return NewImageBufferFromReader(bytes.NewReader(b))
}
func NewImageBufferFromFile(fileName string) (*ImageBuffer, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewImageBufferFromReader(file)
}
func NewImageBufferFromReader(reader io.Reader) (*ImageBuffer, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return NewImageBufferFromImage(img)
}
func NewImageBufferFromImage(img image.Image) (*ImageBuffer, error) {
	rect := img.Bounds()
	image := NewImageBuffer(rect.Max.X-rect.Min.X, rect.Max.Y-rect.Min.Y)
	for x := 0; x < image.Width; x++ {
		for y := 0; y < image.Height; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			alpha := float64(a) / 0xffff
			gray := float64(r+g+b) / 3 / 0xffff
			image.Alpha[y][x] = alpha
			if alpha > 0 {
				image.Gray[y][x] = gray / alpha
			}
		}
	}
	return image, nil
}
func NewImageBuffer(width int, height int) *ImageBuffer {
	gray := make([][]float64, height)
	alpha := make([][]float64, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]float64, width)
		alpha[y] = make([]float64, width)
	}
	return &ImageBuffer{width, height, gray, alpha}
}
func (image *ImageBuffer) ColorAt(x, y float64) (float64, float64) {
	if x < 0 || y < 0 || x > 1 || y > 1 {
		return 0, 0
	}
	ix := int(x * float64(image.Width))
	if ix >= image.Width {
		ix = image.Width - 1
	}
	iy := int(y * float64(image.Height))
	if iy >= image.Height {
		iy = image.Height - 1
	}
	return image.Gray[iy][ix], image.Alpha[iy][ix]
}
func (image *ImageBuffer) Sub(x, y, w, h float64) *SubImage {
	return &SubImage{image, x, y, w, h}
}
func (image *ImageBuffer) Plot(x, y int, gray, alpha float64) {
	if x < 0 || y < 0 || x >= image.Width || y >= image.Height {
		return
	}
	dstGray, dstAlpha := image.Gray[y][x], image.Alpha[y][x]
	newAlpha := dstAlpha + alpha - dstAlpha*alpha
	image.Alpha[y][x] = newAlpha
	if newAlpha == 0 {
		image.Gray[y][x] = 0
	} else {
		image.Gray[y][x] = (dstGray*dstAlpha*(1-alpha) + gray*alpha) / newAlpha
	}
}
func (screen *ImageBuffer) Draw(image Image, x, y, w, h float64) {
	if x+w < 0 || y+h < 0 || float64(screen.Width) < x || float64(screen.Height) < y {
		return
	}
	x0, x1 := int(x), int(x+w)
	y0, y1 := int(y), int(y+h)
	if w < 0 {
		x0, x1 = x1, x0
	}
	if h < 0 {
		y0, y1 = y1, y0
	}
	for ix := x0; ix <= x1; ix++ {
		for iy := y0; iy <= y1; iy++ {
			gray, alpha := image.ColorAt((float64(ix)-x)/w, (float64(iy)-y)/h)
			screen.Plot(ix, iy, gray, alpha)
		}
	}
}
func (screen *ImageBuffer) RotateDraw(image Image, x, y, w, h, deg float64) {
	cx, cy := x+w/2, y+h/2
	theta := deg * math.Pi / 180
	sin, cos := math.Sin(theta), math.Cos(theta)
	xsize, ysize := math.Abs(w)/2, math.Abs(h)/2
	xmin := cx - math.Abs(cos)*xsize - math.Abs(sin)*ysize
	xmax := cx + math.Abs(cos)*xsize + math.Abs(sin)*ysize
	ymin := cy - math.Abs(sin)*xsize - math.Abs(cos)*ysize
	ymax := cy + math.Abs(sin)*xsize + math.Abs(cos)*ysize
	if xmax < 0 || ymax < 0 || xmin > float64(screen.Width) || ymin > float64(screen.Height) {
		return
	}
	x0, x1 := int(math.Ceil(xmin)), int(math.Ceil(xmax))
	for ix := x0; ix < x1; ix++ {
		basew, diffw := cy+(float64(ix)-cx)*sin/cos, math.Abs(ysize/cos)
		baseh, diffh := cy-(float64(ix)-cx)*cos/sin, math.Abs(xsize/sin)
		var y0, y1 int
		if cos == 0 {
			y0, y1 = int(math.Ceil(cy-xsize)), int(math.Ceil(cy+xsize))
		} else if sin == 0 {
			y0, y1 = int(math.Ceil(cy-ysize)), int(math.Ceil(cy+ysize))
		} else {
			y0 = int(math.Ceil(math.Max(basew-diffw, baseh-diffh)))
			y1 = int(math.Ceil(math.Min(basew+diffw, baseh+diffh)))
		}
		for iy := y0; iy < y1; iy++ {
			dx, dy := float64(ix)-cx, float64(iy)-cy
			gray, alpha := image.ColorAt((dx*cos+dy*sin+w/2)/w, (dy*cos-dx*sin+h/2)/h)
			screen.Plot(ix, iy, gray, alpha)
		}
	}
}

func (image *ImageBuffer) StringLines() []string {
	lines := make([]string, image.Height/2)
	buf := make([]byte, image.Width)
	for y := 0; y < image.Height/2; y++ {
		for x := 0; x < image.Width; x++ {
			ug, ua := image.Gray[2*y][x], image.Alpha[2*y][x]
			up := int(16 * (ug*ua + 1*(1-ua)))
			dg, da := image.Gray[2*y+1][x], image.Alpha[2*y+1][x]
			down := int(16 * (dg*da + 1*(1-da)))
			if up < 0 {
				up = 0
			}
			if down < 0 {
				down = 0
			}
			if up > 0xf {
				up = 0xf
			}
			if down > 0xf {
				down = 0xf
			}
			buf[x] = charTable[up][down]
		}
		lines[y] = string(buf)
	}
	return lines
}
func (image *ImageBuffer) String() string {
	return strings.Join(image.StringLines(), "\n")
}
func (image *ImageBuffer) Print() {
	os.Stdout.WriteString("\x1B[1;1H" + image.String())
}

var charTable []string = []string{
	"MMMMMM###TTTTTTT",
	"QQBMMNW##TTTTTV*",
	"QQQBBEK@PTTTVVV*",
	"QQQmdE88P9VVVV**",
	"QQQmdGDU0YVV77**",
	"pQQmAbk65YY?7***",
	"ppgAww443vv?7***",
	"pggyysxcJv??7***",
	"pggyaLojrt<<+**\"",
	"gggaauuj{11!//\"\"",
	"gggaauui])|!/~~\"",
	"ggaauui]((;::~~^",
	"ggaauu](;;::-~~'",
	"ggauu(;;;;---~``",
	"gaau;;,,,,,...``",
	"gau,,,,,,,,...  "}
