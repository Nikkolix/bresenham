package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/png"
	"os"
)

//Image is a custom implementation of golang's image.Image interface.

type Image struct {
	Pixel []uint32
	W     int
	H     int
}

func NewImage(w, h int) *Image {
	return &Image{
		Pixel: make([]uint32, w*h),
		W:     w,
		H:     h,
	}
}

func (i *Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i *Image) Bounds() image.Rectangle {
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

func (i *Image) At(x, y int) color.Color {
	c := i.Pixel[x+y*i.W]
	return color.RGBA{
		R: uint8(c & 0b11111111),
		G: uint8((c >> 8) & 0b11111111),
		B: uint8((c >> 16) & 0b11111111),
		A: uint8((c >> 24) & 0b11111111),
	}
}

func (i *Image) Set(x, y int, c uint32) {
	i.Pixel[x+y*i.W] = c
}

func (i *Image) Clear(c uint32) {
	for y := 0; y < i.H; y++ {
		for x := 0; x < i.W; x++ {
			i.Set(x, y, c)
		}
	}
}

func (i *Image) SaveToPNG(filename string) {
	file, _ := os.Create(filename)
	handle(png.Encode(file, i))
	handle(file.Close())
}

//Gif is a custom implementation of golang's image.Image interface.

type Gif struct {
	images []*image.Paletted
}

func NewGif() *Gif {
	return &Gif{
		images: []*image.Paletted{},
	}
}

func (g *Gif) AppendImage(img *Image) {
	paletteImg := image.NewPaletted(image.Rect(0, 0, img.W, img.H), palette.Plan9)
	for x := 0; x < img.W; x++ {
		for y := 0; y < img.H; y++ {
			paletteImg.Set(x, y, img.At(x, y))
		}
	}
	g.images = append(g.images, paletteImg)
}

func (g *Gif) Save(filename string) {
	file, _ := os.Create(filename)
	img := &gif.GIF{
		Image:           g.images,
		Delay:           delay(10, len(g.images)),
		LoopCount:       0,
		Disposal:        nil,
		Config:          image.Config{},
		BackgroundIndex: 0,
	}
	handle(gif.EncodeAll(file, img))
	handle(file.Close())
}
