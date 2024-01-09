package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/png"
	"math"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"time"
)

func main() {
	/*img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	BresenhamOptimized(5, 20, 90, 50, img, c)
	BresenhamOptimized(10, 15, 95, 30, img, c)
	SaveToPng(img, "final.png")*/

	for i := 16; i < 25; i++ {
		fmt.Println(i, ":")
		//benchmark(1<<i, 1<<8, LineDraw, IncrementalLineDraw, BresenhamFloat, BresenhamOptimized, BresenhamSetRGBA)
		benchmark(1<<i, 1<<8, BresenhamOptimized)
		benchmarkCustomRGBA(1<<i, 1<<8, BresenhamOptimizedCustomRGBA)
		fmt.Println("")
	}

	/*img = image.NewRGBA(image.Rect(0, 0, 100, 100))
	c = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	BresenhamOptimized(5, 20, 90, 50, img, c)
	Wu(10, 15, 95, 30, c, img)
	SaveToPng(img, "final.png")*/
}

// BresenhamGif integer algorithm
func BresenhamGif(x1, y1, x2, y2 int, img *image.RGBA, c *color.RGBA) {
	var images []*image.Paletted
	dx := x2 - x1
	dy := y2 - y1
	d := 2*dy - dx
	for x1 <= x2 {

		//copy img
		imgCopy := image.NewPaletted(img.Rect, palette.Plan9)
		for x := 0; x < img.Rect.Dx(); x++ {
			for y := 0; y < img.Rect.Dy(); y++ {
				imgCopy.Set(x, y, img.At(x, y))
			}
		}
		images = append(images, imgCopy)

		img.Set(x1, y1, c)
		x1++
		if d <= 0 {
			d += 2 * dy
		} else {
			d += 2 * (dy - dx)
			y1++
		}
	}

	SaveToGif(images, "img.gif")
}

type CustomRGBA struct {
	Pixel []uint32
	W     int
	H     int
}

func (i *CustomRGBA) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *CustomRGBA) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: i.W,
			Y: i.H,
		},
	}
}

func (i *CustomRGBA) At(x, y int) color.Color {
	c := i.Pixel[x+y*i.W]
	return color.RGBA{
		R: uint8(c & 0b11111111),
		G: uint8((c >> 8) & 0b11111111),
		B: uint8((c >> 16) & 0b11111111),
		A: uint8((c >> 24) & 0b11111111),
	}
}

// BresenhamOptimized integer algorithm (optimized)
func BresenhamOptimizedCustomRGBA(x1, y1, x2, y2 int, img *CustomRGBA, c uint32) {
	iMod := [2]int{1, img.W + 1}

	i := img.W*y1 + x1
	img.Pixel[i] = c
	end := img.W*y2 + x2
	img.Pixel[end] = c

	dx := x2 - x1
	dy2 := (y2 - y1) << 1
	dx2 := dx << 1
	e := -dy2 + dx // e is negated

	b := int(uint64(e) >> 63)
	e += dx2*b - dy2
	i += iMod[b]

	for i < end {
		b := int(uint64(e) >> 63)
		img.Pixel[i] = c
		e -= dy2
		i += iMod[b]
		e += dx2 * b
	}
}

// BresenhamOptimized integer algorithm (optimized)
func BresenhamOptimized(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	i := img.PixOffset(x1, y1)
	s := img.Pix[i : i+4 : i+4]
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
	end := img.PixOffset(x2, y2)
	s = img.Pix[end : end+4 : end+4]
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
	dx := x2 - x1
	dy2 := (y2 - y1) << 1
	dx2 := dx << 1
	e := -dy2 + dx // e is negated

	b := int(uint64(e) >> 63)
	i += 4 + b*img.Stride
	e += dx2*b - dy2

	for i < end {
		s := img.Pix[i : i+4 : i+4]
		s[0] = c.R
		s[1] = c.G
		s[2] = c.B
		s[3] = c.A
		b := int(uint64(e) >> 63)
		i += 4 + b*img.Stride
		e += dx2*b - dy2
	}
}

// BresenhamSetRGBA integer algorithm (optimized)
func BresenhamSetRGBA(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	img.SetRGBA(x1, y1, c)
	img.SetRGBA(x2, y2, c)
	dx := x2 - x1
	x2--
	dy2 := (y2 - y1) << 1
	dx2 := dx << 1
	e := -dy2 + dx // e is negated
	for x1 < x2 {
		x1++
		b := int(uint64(e) >> 63)
		e += dx2*b - dy2
		y1 += b
		img.SetRGBA(x1, y1, c)
	}
}

// BresenhamFloat algorithm
func BresenhamFloat(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	slope := float64(dy) / float64(dx)
	e := slope - 0.5
	for x := x1; x <= x2; x++ {
		img.SetRGBA(x, y1, c)
		e += slope
		if e > 0 {
			e--
			y1++
		}
	}
}

// IncrementalLineDraw algorithm
func IncrementalLineDraw(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	y := float64(y1)
	slope := float64(dy) / float64(dx)
	for x := x1; x <= x2; x++ {
		img.SetRGBA(x, round(y), c)
		y = y + slope
	}
}

// LineDraw algorithm
func LineDraw(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	slope := float64(dy) / float64(dx)
	for x := x1; x <= x2; x++ {
		y := slope*float64(x-x1) + float64(y1)
		img.SetRGBA(x, round(y), c)
	}
}

// wu inspired by wikipedia line draw wu ("https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm")

func fpart(x float64) float64 {
	return x - math.Floor(x)
}

func rfpart(x float64) float64 {
	return 1 - fpart(x)
}

func plot(x, y float64, c float64, rgba color.RGBA, img *image.RGBA) {
	c = max(0.0, min(c, 1.0))
	rgba.R = uint8(float64(rgba.R) * c)
	rgba.G = uint8(float64(rgba.G) * c)
	rgba.B = uint8(float64(rgba.B) * c)
	img.SetRGBA(round(x), round(y), rgba)
}

func Wu(x0, y0, x1, y1 float64, rgba color.RGBA, img *image.RGBA) {
	var images []*image.Paletted

	steep := math.Abs(y1-y0) > math.Abs(x1-x0)

	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0

	}

	dx := x1 - x0
	dy := y1 - y0

	gradient := 1.0
	if dx != 0 {
		gradient = dy / dx
	}

	xend := math.Round(x0)
	yend := y0 + gradient*(xend-x0)
	xgap := rfpart(x0 + 0.5)
	xpxl1 := xend
	ypxl1 := math.Floor(yend)
	if steep {
		plot(ypxl1, xpxl1, rfpart(yend)*xgap, rgba, img)
		plot(ypxl1+1, xpxl1, fpart(yend)*xgap, rgba, img)
	} else {
		plot(xpxl1, ypxl1, rfpart(yend)*xgap, rgba, img)
		plot(xpxl1, ypxl1+1, fpart(yend)*xgap, rgba, img)
	}
	intery := yend + gradient

	xend = math.Round(x1)
	yend = y1 + gradient*(xend-x1)
	xgap = fpart(x1 + 0.5)
	xpxl2 := xend
	ypxl2 := math.Floor(yend)
	if steep {
		plot(ypxl2, xpxl2, rfpart(yend)*xgap, rgba, img)
		plot(ypxl2+1, xpxl2, fpart(yend)*xgap, rgba, img)
	} else {
		plot(xpxl2, ypxl2, rfpart(yend)*xgap, rgba, img)
		plot(xpxl2, ypxl2+1, fpart(yend)*xgap, rgba, img)
	}

	if steep {
		for x := xpxl1 + 1; x <= xpxl2-1; x++ {
			//copy img
			imgCopy := image.NewPaletted(img.Rect, palette.Plan9)
			for x := 0; x < img.Rect.Dx(); x++ {
				for y := 0; y < img.Rect.Dy(); y++ {
					imgCopy.Set(x, y, img.At(x, y))
				}
			}
			images = append(images, imgCopy)

			plot(math.Floor(intery), x, rfpart(intery), rgba, img)
			plot(math.Floor(intery)+1, x, fpart(intery), rgba, img)
			intery = intery + gradient
		}
	} else {
		for x := xpxl1 + 1; x <= xpxl2-1; x++ {
			//copy img
			imgCopy := image.NewPaletted(img.Rect, palette.Plan9)
			for x := 0; x < img.Rect.Dx(); x++ {
				for y := 0; y < img.Rect.Dy(); y++ {
					imgCopy.Set(x, y, img.At(x, y))
				}
			}
			images = append(images, imgCopy)

			plot(x, math.Floor(intery), rfpart(intery), rgba, img)
			plot(x, math.Floor(intery)+1, fpart(intery), rgba, img)
			intery = intery + gradient
		}
	}

	SaveToGif(images, "imgwu.gif")
}

//helper

func round(f float64) int {
	if f > 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Uint32()),
		G: uint8(rand.Uint32()),
		B: uint8(rand.Uint32()),
		A: 255,
	}
}

func benchmark(N int, res int, rasterizer ...func(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA)) {
	start := time.Now()

	for index := range rasterizer {
		img := image.NewRGBA(image.Rect(0, 0, res, res))
		s := time.Now()

		for i := 0; i < N; i++ {
			c := randomColor()
			r1 := rand.Intn(res)
			r2 := rand.Intn(res)
			r3 := rand.Intn(res)
			r4 := rand.Intn(res)
			x1 := min(r1, r2)
			x2 := max(r1, r2)
			y1 := min(r3, r4)
			y2 := max(r3, r4)
			dx := x2 - x1
			dy := y2 - y1
			if dx < dy {
				y1, x1 = x1, y1
				y2, x2 = x2, y2
			}
			rasterizer[index](x1, y1, x2, y2, img, c)
		}

		fmt.Println(runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name(), time.Now().Sub(s).Milliseconds())
		SaveToPng(img, runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name()+".png")

	}

	fmt.Println(time.Now().Sub(start))
}

func benchmarkCustomRGBA(N int, res int, rasterizer ...func(x1, y1, x2, y2 int, img *CustomRGBA, c uint32)) {
	start := time.Now()

	for index := range rasterizer {
		img := &CustomRGBA{
			Pixel: make([]uint32, res*res),
			W:     res,
			H:     res,
		}
		s := time.Now()

		for i := 0; i < N; i++ {
			c := rand.Uint32() | (0b11111111 << 24)
			r1 := rand.Intn(res)
			r2 := rand.Intn(res)
			r3 := rand.Intn(res)
			r4 := rand.Intn(res)
			x1 := min(r1, r2)
			x2 := max(r1, r2)
			y1 := min(r3, r4)
			y2 := max(r3, r4)
			dx := x2 - x1
			dy := y2 - y1
			if dx < dy {
				y1, x1 = x1, y1
				y2, x2 = x2, y2
			}
			rasterizer[index](x1, y1, x2, y2, img, c)
		}

		fmt.Println(runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name(), time.Now().Sub(s).Milliseconds())
		SaveToPng(img, runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name()+".png")

	}

	fmt.Println(time.Now().Sub(start))
}

func SaveToPng(img image.Image, filename string) {
	file, _ := os.Create(filename)
	handle(png.Encode(file, img))
	handle(file.Close())
}

func Delay(d, n int) []int {
	var ds = make([]int, n)
	for i := 0; i < n; i++ {
		ds[i] = d
	}
	return ds
}

func SaveToGif(images []*image.Paletted, filename string) {
	file, _ := os.Create(filename)
	img := &gif.GIF{
		Image:           images,
		Delay:           Delay(10, len(images)),
		LoopCount:       1,
		Disposal:        nil,
		Config:          image.Config{},
		BackgroundIndex: 0,
	}
	handle(gif.EncodeAll(file, img))
	handle(file.Close())
}
