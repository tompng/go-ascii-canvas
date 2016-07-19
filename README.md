# Golang Asciiart Canvas Library

```
      ,gg.      .,,,.      ....       ..     ..    .......    ......         ...
   ,gM#TTMr    pMMTMM,    :MMMMMM;   `MM     MM   `MMMMMMM... MMMMMMM,       MM!
  pM^        (M!     MM.  :M:   .M#` `MM,,,,,MM   -MQ55;;:7YbjMM`  ,gT`      rM~
 `MM   ggg,  rM.     MM:  :Mi,ggMT`  `MMMMMMMMM,, 4MMMMMM:-:?bMMMMMMT        !M`
  #M,  `~Mr  :M;     MM~  :M?"^`  ,,,4MMsb5Y4MMbbwdMd;;bQuu-:bMM` #MM        :M`
  'MQg  `Mr  `TQ,  .gM*   :M:,,,sxJr75MM++<!!MM66bdMQddMMMY---MM. `TMg,      .g.
   "TMMppMr    TMMpMM"    !MysbYY:---|MM--:(cMM"`'*MMMMMMM<-:]MM`   'MM;     MM4
                        ~u4kr++/------~~;])'``      .rpaybAuub?                 
                       ubbc:-----------(j;`        ,mMMMMQQmbbb4,,              
                     (wbc/------------;J*`        :PQMMMMMMmDr<?r6u,~           
                    ubb+-------------/uv`          rMMMMMMMr1(:-:7:6s,~         
                  ,oL?/--------------:wu~          ^~rVT#r^^(1-----7:6i,        
                 ,(b?7---------------+j4,            ``    ~v;-------7:b(,.     
                ,bJ7/-----------:-----74x                 .cc----------?bb;,,,,c
               4b5+--------;;;cv)4ci]::44m;             .,;r------------/-?bLb5"
             `bkY------:r**" `'  `*`74irtdAQ;  .      ,,u4//--------------:-<bL~
            `,br;----:r"''     ..rpaarrbuj<KKKg(uau,xxcr::------------------/|b.
            ,YbY/----1``      `rQMMMMQpr4x]-:?<amdww?-;:---------------------:7r
           "ub5;----;~        (PQMMMMMMy<x4(-gpQMMMpLLLbu;----------------------
           ;bs?----:|.         rMMMMMMMr-7bApQMMMMMPY11Jbu1---------------------
           *x6:---:+r~         ^~rVT#r^``ummMMMMMMNu111tjk!:--------------------
            4u/-----r,           ``     ;6dvMMM#TT1111{usu{;--------------------
         ,,rLbu------r,.              ..cY//1xJJ11111{oo"/jcc{-;----------------
```

## How to use

### Install
```sh
$ go get github.com/tompng/go-ascii-canvas
```

### Basic
```go
package main

import (
	"fmt"
	"github.com/tompng/go-ascii-canvas"
)

func main() {
	img, err := asciicanvas.NewImageBufferFromFile("gopher.png")
	if err != nil {
		panic("cannot read file")
	}
	canvas := asciicanvas.NewImageBuffer(80, 48)
	canvas.Draw(img, 0, 0, 40, 40) // img, x, y, w, h
	fmt.Println(canvas.String())
}
```

### Animation
```go
for {
	terminalWidth, terminalHeight := asciicanvas.GetWinSize()
	canvas := asciicanvas.NewImageBuffer(terminalWidth, 2*terminalHeight)
	canvas.Draw(img, x, y, 80, 80)
	x += 0.1
	y += 0.2
	canvas.Print()
	time.Sleep(50 * time.Millisecond)
}
```

### Other
```go
// Size
fmt.Print(canvas.Width, canvas.Height, img.Width, img.Height)

// Sub image
subImage := img.Sub(0.1, 0.2, 0.5, 0.5) // x, y, w, h (0~1)
canvas.Draw(subImage, 0, 0, 40, 40)

// Draw with rotation
angle := 120.0
canvas.RotateDraw(img, x, y, w, h, angle)

// Draw canvas into another canvas
canvas1.Draw(img, 0, 0, 100, 100)
canvas2.Draw(canvas1, 0, 0, 40, 40)

// Jpeg image
import (
	"github.com/tompng/go-ascii-canvas"
	_ "image/jpeg"
)
...
img, _ := asciicanvas.NewImageBufferFromFile("image.jpg")

// Plot
for x, y := 0, 0; x < 20; x, y = x+1, y+1 {
  gray, alpha := float64(x)/20, 0.5
  canvas.Plot(x, y, gray, alpha)
}
```

### Drawing text
- Prepare a bitmap font data image
- Create a subimage corresponds to each chars
- Draw it
```go
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
...
DrawText(canvas, "Hello World", 0, 0, 16)
```
