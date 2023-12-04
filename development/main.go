package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/png"
	"os"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	c := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Bresenham(5, 20, 90, 50, img, &c)
	Bresenham(10, 15, 95, 30, img, &c)
	SaveToPng(img, "final.png")
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

// Bresenham algorithm
func Bresenham(x1, y1, x2, y2 int, img *image.RGBA, c *color.RGBA) {
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
