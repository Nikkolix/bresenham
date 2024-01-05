package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/png"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"time"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	/*img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Bresenham(5, 20, 90, 50, img, &c)
	Bresenham(10, 15, 95, 30, img, &c)
	SaveToPng(img, "final.png")*/

	for i := 0; i < 32; i++ {
		fmt.Println(i, ":")
		benchmark(1<<i, 1<<6, LineDraw, IncrementalLineDraw, BresenhamFloat, Bresenham)
		fmt.Println("")
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

		fmt.Println(time.Now().Sub(s).Milliseconds())
		SaveToPng(img, runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name()+".png")

	}

	fmt.Println(time.Now().Sub(start))
}

func SaveToPng(img *image.RGBA, filename string) {
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

// Bresenham integer algorithm
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

// Bresenham integer algorithm
func Bresenham(x1, y1, x2, y2 int, img *image.RGBA, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	dx2 := dx << 1
	dy2 := dy << 1
	e := dy2 - dx
	img.SetRGBA(x1, y1, c)
	for x1 < x2 {
		x1++
		img.SetRGBA(x1, y1, c)
		b := ((-e) >> 63) & 1
		e = e + dy2 - dx2*b
		y1 = y1 + b
	}
	img.SetRGBA(x2, y2, c)
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
		//b := (math.Float64bits(-e) >> 63) ^ 1
		if e > 0 {
			e--
			y1++
		}
	}
}

func round(f float64) int {
	if f > 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
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
