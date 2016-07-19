# Golang Asciiart Canvas Library

```
           ,(w6?-----------(j*"`  `'*UdMM-:bJ                                   
          ,57/------------o*`       ,,vj7-:b~                                   
       ,;57--------------;"       ,pMMQpAu?"                                    
      ubY/--:;|+""*7v(:--x        <MMMMMD54,                                    
     ;6+---:"`       *4i:4.       `*##V<1:<6u.                                  
     4---:]      .yQMQyrb{Q;          ;v----?b                                  
LcJjLb----r      rMMMMMr4x?Kg(,    .;J7------?.                                 
7;-;bb----j      ^*T#T*/s4ggpQQ?7777/---------s,                 ,.,            
-ggpb5----?.          .umQMMMMMLLu:-----------?bLLJ7ru;.   .,,54675YxJ,,        
;QMbv:-----x,.       ,yYvMMMMNt1Jb)------------<b*   ~bu,,7Y7/--------Y5s,      
/vMb|-------:+?5xx45V:-ub311111usv:------------:5;   LbY/-----------~--!4b,,u,,
--:kj------------------3bo11{{o"/jj;-------------?;,?+--::;;------]**"*"bb5:-?b
u;;ob(-----------------:+77?7**~,-7L---------u6v3]b7-::"` ` 7x--;7      /UgG:u"
"Y66bb---------------------:/1, ,;jv--------c:;-u7--;*     , c-]*    .ugucT-u*  
    <b------------------------/jv/--------]7-gps7--(     gMM(wJY    .pMMMw]Y"   
     4;----------------------------------(4;8MY:--:/    (MMMj<Q`    *MMM0k      
     /{----------------------------------b:7EY----xc    ^T*u7-P;     .(/:t      
      wu--------------------;;-;u(:-----jc-;w-----3w     .apQpQD,,,,,/--];      
      5b]-----------------;Lb`  ``*;----:6Lb/------+xcu?Y1MMMME(-:/-----6"      
      `6b:---------------(5Y6;    .k------bY-----------]uvv1Y1JY-------u(       
       `6bi--------------:-----+77/------ik------------rr11{uv?--------b`       
        ^6b-----------------------------:b:------------/?*o+ij--------:b"7,.    
         `*u:---------------------------b<--------------1 ^/vt-------:b`  b`    
           5s:-------------------------Lb---------------/,117--------j* ,;"     
```

## How to use

### basic
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

### animation
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
angle := 120
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

### Drawing Text
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
