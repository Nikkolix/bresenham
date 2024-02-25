package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

func benchmark(N int, res int, rasterizer ...func(x1, y1, x2, y2 int, img *Image, c uint32)) {
	start := time.Now()

	for index := range rasterizer {
		img := &Image{
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
		img.SaveToPNG(runtime.FuncForPC(reflect.ValueOf(rasterizer[index]).Pointer()).Name() + ".png")

	}

	fmt.Println(time.Now().Sub(start))
}

func benchmarkSame(N int, res int, rasterizer ...func(x1, y1, x2, y2 int, img *Image, c uint32)) {
	img := &Image{
		Pixel: make([]uint32, res*res),
		W:     res,
		H:     res,
	}
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
		for index := range rasterizer {
			rasterizer[index](x1, y1, x2, y2, img, c)
		}
	}
}
