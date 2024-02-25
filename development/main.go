package main

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"

	"github.com/bits-and-blooms/bitset"
)

func main() {
	img := NewImage(100, 40)
	img.Clear(0xffffffff)
	col := uint32(0xff0000ff)
	BresenhamOptimized(5, 5, 90, 30, img, col)
	BresenhamGif(10, 15, 95, 35, img, col)
	img.SaveToPNG("out/bresenham.png")

	GenerateCompareImage(BresenhamFloat, Bresenham)

	img2 := NewImage(100, 40)
	img2.Clear(0xffffffff)
	col2 := uint32(0xff0000ff)
	BresenhamOptimized(5, 5, 90, 30, img2, col2)
	Wu(10, 15, 95, 35, img2, col2)
	img2.SaveToPNG("out/bresenham-wu.png")

	img3 := NewImage(100, 100)
	img3.Clear(0xffffffff)
	col3 := uint32(0xff0000ff)
	PrimitivLineDraw(5, 5, 90, 30, img3, col3)
	PrimitivLineDraw(5, 5, 90, 60, img3, col3)
	PrimitivLineDraw(5, 5, 90, 90, img3, col3)
	PrimitivLineDraw(5, 5, 60, 90, img3, col3)
	PrimitivLineDraw(5, 5, 30, 90, img3, col3)
	img3.SaveToPNG("out/line-draw-show.png")

	img4 := NewImage(100, 40)
	img4.Clear(0xffffffff)
	col4 := uint32(0xff0000ff)
	BresenhamOptimized(5, 5, 90, 30, img4, col4)
	BresenhamOptimizedGCD(10, 15, 95, 35, img4, col4)
	img4.SaveToPNG("out/gcd.png")

	img5 := NewImage(100, 40)
	img5.Clear(0xffffffff)
	col5 := uint32(0xff0000ff)
	BresenhamOptimized(5, 5, 90, 30, img5, col5)
	BresenhamOptimizedGCDMirror(10, 15, 95, 35, img5, col5)
	img5.SaveToPNG("out/gcd-mirror.png")

	img6 := NewImage(100, 40)
	img6.Clear(0xffffffff)
	col6 := uint32(0xff0000ff)
	BresenhamOptimized(5, 5, 90, 30, img6, col6)
	BresenhamOptimizedGCDMirrorBitset(10, 15, 95, 35, img6, col6)
	img6.SaveToPNG("out/gcd-mirror-bitset.png")
}

func GenerateCompareImage(alg1, alg2 func(x1, y1, x2, y2 int, img *Image, c uint32)) {
	c := uint32(0xff0000ff)

	img := NewImage(256, 256)
	alg1(20, 30, 200, 160, img, c)
	alg1(30, 30, 33, 31, img, c) // 4 - 2
	alg1(40, 30, 43, 32, img, c) // 4 - 3
	alg1(50, 30, 54, 33, img, c) // 5 - 4
	alg1(60, 30, 64, 32, img, c) // 5 - 3
	alg1(50, 220, 151, 222, img, c)
	alg1(50, 230, 151, 237, img, c)

	alg2(20, 40, 200, 170, img, c)
	alg2(130, 30, 133, 31, img, c) // 4 - 2
	alg2(140, 30, 143, 32, img, c) // 4 - 3
	alg2(150, 30, 154, 33, img, c) // 5 - 4
	alg2(160, 30, 164, 32, img, c) // 5 - 3
	alg2(50, 200, 151, 202, img, c)
	alg2(50, 210, 151, 217, img, c)

	alg1Name := last(strings.Split(runtime.FuncForPC(reflect.ValueOf(alg1).Pointer()).Name(), "."))
	alg2Name := last(strings.Split(runtime.FuncForPC(reflect.ValueOf(alg2).Pointer()).Name(), "."))

	img.SaveToPNG("out/compare-" + alg1Name + "-" + alg2Name + ".png")
}

// BresenhamGif integer algorithm
func BresenhamGif(x1, y1, x2, y2 int, img *Image, c uint32) {
	gif := NewGif()

	dx := x2 - x1
	dy := y2 - y1
	d := 2*dy - dx
	for x1 <= x2 {

		//copy img
		gif.AppendImage(img)

		img.Set(x1, y1, c)
		x1++
		if d <= 0 {
			d += 2 * dy
		} else {
			d += 2 * (dy - dx)
			y1++
		}
	}

	gif.Save("out/img.gif")
}

func BresenhamOptimizedGCDMirrorBitset(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 {
		return
	}

	if dy == 0 {
		i := img.W*y1 + x1
		end := img.W*y1 + x2
		for i <= end {
			img.Pixel[i] = c
			i++
			//img.Pixel[end] = c
			//end--
		}
		return
	}
	if dx == dy {
		i := img.W*y1 + x1
		end := img.W*y2 + x2
		off := 1 + img.W
		for i <= end {
			img.Pixel[i] = c
			i += off
			//img.Pixel[end] = c
			//end -= off
		}
		return
	}

	i := uint(0)

	g := gcd(dx, dy)

	ndy2 := -(dy << 1) // negated
	dx2 := dx << 1
	e := ndy2 + dx // e is negated

	end := uint(dx / g)

	bits := bitset.New(end)

	b := uint(uint64(e) >> 63)
	z := uint(^(uint64(-e)) >> 63)
	e += dx2*int(b) + ndy2

	end--

	bits.SetTo(i, b == 1)
	bits.SetTo(end, z == 1)

	i++
	end--

	for i <= end {
		b := int(uint64(e) >> 63)
		bits.SetTo(i, e < 0)
		bits.SetTo(end, e <= 0)
		i++
		end--
		e += ndy2 + dx2*b
	}

	index := x1 + img.W*y1
	img.Pixel[index] = c
	for a := 0; a < g; a++ {
		for d := uint(0); d < bits.Len(); d++ {
			if bits.Test(d) {
				index += img.W
			}
			index++
			img.Pixel[index] = c
		}
	}

	fmt.Println(bits.DumpAsBits(), "  ", bits.Len())
}

// BresenhamOptimizedGCDMirrorExp integer algorithm (optimized experimental)
func BresenhamOptimizedGCDMirrorExp(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 {
		return
	}

	i := x1
	end := img.W*dy + x2
	pixel := img.Pixel[img.W*y1 : img.W*(y2+1) : img.W*(y2+1)]

	if dy == 0 {
		for i <= end {
			pixel[i] = c
			pixel[end] = c
			i++
			end--
		}
		return
	}
	if dx == dy {
		off := 1 + img.W
		for i <= end {
			pixel[i] = c
			pixel[end] = c
			i += off
			end -= off
		}
		return
	}

	iMod := [2]int{1, img.W + 1}

	g := gcd(dx, dy)
	goff := dx + dy*img.W

	ndy2 := -(dy << 1) // negated
	dx2 := dx << 1
	e := ndy2 + dx // e is negated
	pixel[i] = c

	gpo := goff / g

	end = i + gpo

	for j := i + gpo; j <= i+goff; j += gpo {
		pixel[j] = c
	}

	pixel[end] = c

	b := int(uint64(e) >> 63)
	z := int(^(uint64(-e)) >> 63)
	e += dx2*b + ndy2
	i += iMod[b]
	end -= iMod[z]

	if g > 1 {
		for i <= end {
			b := int(uint64(e) >> 63)
			z := int(^(uint64(-e)) >> 63)
			pixel[i] = c
			pixel[end] = c
			off := end - i
			j := i + gpo
			m := i + goff
			for j < m {
				pixel[j] = c
				pixel[j+off] = c
				j += gpo
			}
			e += ndy2
			i += iMod[b]
			end -= iMod[z]
			e += dx2 * b
		}
	} else {
		for i <= end {
			pixel[i] = c
			pixel[end] = c
			if e < 0 {
				i += img.W + 1
				end -= img.W + 1
				e += dx2 + ndy2
			} else {
				i++
				e += ndy2
				if e == 0 {
					end -= img.W + 1
				} else {
					end--
				}
			}
		}
	}
}

// BresenhamOptimizedGCDMirror integer algorithm (optimized)
func BresenhamOptimizedGCDMirror(x1, y1, x2, y2 int, img *Image, c uint32) {
	gif := NewGif()

	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 {
		return
	}
	if dy == 0 {
		i := img.W*y1 + x1
		end := img.W*y2 + x2
		for i <= end {
			img.Pixel[i] = c
			i++
			img.Pixel[end] = c
			end--
		}
		return
	}
	if dx == dy {
		i := img.W*y1 + x1
		end := img.W*y1 + x2
		off := 1 + img.W
		for i <= end {
			img.Pixel[i] = c
			i += off
			img.Pixel[end] = c
			end -= off
		}
		return
	}

	iMod := [2]int{1, img.W + 1}

	i := img.W*y1 + x1

	g := gcd(dx, dy)
	goff := dx + dy*img.W

	ndy2 := -(dy << 1) // negated
	dx2 := dx << 1
	e := ndy2 + dx // e is negated
	img.Pixel[i] = c

	gpo := goff / g

	end := i + gpo

	for j := i + gpo; j <= i+goff; j += gpo {
		img.Pixel[j] = c
	}

	img.Pixel[end] = c

	b := int(uint64(e) >> 63)
	z := int(^(uint64(-e)) >> 63)
	e += dx2*b + ndy2
	i += iMod[b]
	end -= iMod[z]

	if g > 1 {
		for i <= end {
			b := int(uint64(e) >> 63)
			z := int(^(uint64(-e)) >> 63)
			img.Pixel[i] = c
			img.Pixel[end] = c
			off := end - i
			j := i + gpo
			m := i + goff
			for j < m {
				img.Pixel[j] = c
				img.Pixel[j+off] = c
				j += gpo
			}
			e += ndy2
			i += iMod[b]
			end -= iMod[z]
			e += dx2 * b
			gif.AppendImage(img)
		}
	} else {
		for i <= end {
			b := int(uint64(e) >> 63)
			z := int(^(uint64(-e)) >> 63)
			img.Pixel[i] = c
			img.Pixel[end] = c

			e += ndy2
			i += iMod[b]
			end -= iMod[z]
			e += dx2 * b
			gif.AppendImage(img)
		}
	}

	gif.Save("out/gcd-mirror.gif")
}

// BresenhamOptimizedGCD integer algorithm (optimized)
func BresenhamOptimizedGCD(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1

	gif := NewGif()

	if dx == 0 {
		return
	}
	if dy == 0 {
		i := img.W*y1 + x1
		end := img.W*y1 + x2
		for i <= end {
			img.Pixel[i] = c
			i++
		}
		return
	}
	if dx == dy {
		i := img.W*y1 + x1
		end := img.W*y2 + x2
		for i <= end {
			img.Pixel[i] = c
			i += 1 + img.W
		}
		return
	}

	iMod := [2]int{1, img.W + 1}

	g := gcd(dx, dy)
	goff := dx + dy*img.W
	gpo := goff / g

	i := img.W*y1 + x1
	img.Pixel[i] = c
	end := i + gpo
	img.Pixel[end] = c

	for j := i + gpo; j <= i+goff; j += gpo {
		img.Pixel[j] = c
	}

	ndy2 := -(dy << 1) // negated
	dx2 := dx << 1
	e := ndy2 + dx // e is negated

	b := int(uint64(e) >> 63)
	e += dx2*b + ndy2
	i += iMod[b]

	if g > 1 {
		for i < end {
			b := int(uint64(e) >> 63)
			img.Pixel[i] = c
			j := i + gpo
			m := i + goff
			for j < m {
				img.Pixel[j] = c
				j += gpo
			}
			e += ndy2
			i += iMod[b]
			e += dx2 * b

			gif.AppendImage(img)
		}
	} else {
		for i < end {
			b := int(uint64(e) >> 63)
			img.Pixel[i] = c

			e += ndy2
			i += iMod[b]
			e += dx2 * b

			gif.AppendImage(img)
		}
	}

	gif.Save("out/gcd.gif")
}

// BresenhamOptimized integer algorithm (optimized)
func BresenhamOptimized(x1, y1, x2, y2 int, img *Image, c uint32) {
	iMod := [2]int{1, img.W + 1}

	i := img.W*y1 + x1
	img.Pixel[i] = c
	end := img.W*y2 + x2
	img.Pixel[end] = c

	dx := x2 - x1
	dy := y2 - y1
	dy2 := dy << 1
	dx2 := dx << 1
	e := dx - dy2 // e is negated

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

// Bresenham integer algorithm
func Bresenham(x1, y1, x2, y2 int, img *Image, c uint32) {
	img.Set(x1, y1, c)
	img.Set(x2, y2, c)

	dx := x2 - x1
	dy := y2 - y1

	dy2 := dy << 1
	dx2 := dx << 1
	e := dy2 - dx

	x2--

	for x1 < x2 {
		x1++
		b := int(uint64(-e) >> 63)
		e += dy2 - dx2*b
		y1 += b
		img.Set(x1, y1, c)
	}
}

// BresenhamFloat algorithm
func BresenhamFloat(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1
	slope := float64(dy) / float64(dx)
	e := slope - 0.5
	for x := x1; x <= x2; x++ {
		img.Set(x, y1, c)
		if e > 0 {
			e--
			y1++
		}
		e += slope
	}
}

// IncrementalLineDraw algorithm
func IncrementalLineDraw(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1
	y := float64(y1)
	m := float64(dy) / float64(dx)
	for x := x1; x <= x2; x++ {
		img.Set(x, round(y), c)
		y += m
	}
}

// PrimitivLineDraw algorithm
func PrimitivLineDraw(x1, y1, x2, y2 int, img *Image, c uint32) {
	dx := x2 - x1
	dy := y2 - y1
	slope := float64(dy) / float64(dx)
	for x := x1; x <= x2; x++ {
		y := slope*float64(x-x1) + float64(y1)
		img.Set(x, round(y), c)
	}
}

// wu inspired by wikipedia line draw wu ("https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm")

func fpart(x float64) float64 {
	return x - math.Floor(x)
}

func rfpart(x float64) float64 {
	return 1 - fpart(x)
}

func plot(x, y float64, c float64, rgba uint32, img *Image) {
	c = max(0.0, min(c, 1.0))
	r := uint8(0b11111111 & rgba)
	g := uint8(0b11111111 & (rgba >> 8))
	b := uint8(0b11111111 & (rgba >> 16))

	cr, cg, cb, _ := img.At(round(x), round(y)).RGBA()

	r = uint8(min(max(float64(r)*c+float64(cr>>8)*(1-c), 0), 255))
	g = uint8(min(max(float64(g)*c+float64(cg>>8)*(1-c), 0), 255))
	b = uint8(min(max(float64(b)*c+float64(cb>>8)*(1-c), 0), 255))

	rgba = uint32(r) | uint32(g)<<8 | uint32(b)<<16 | 0xff000000

	img.Set(round(x), round(y), rgba)
}

func Wu(x0, y0, x1, y1 float64, img *Image, c uint32) {

	gif := NewGif()

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
		plot(ypxl1, xpxl1, rfpart(yend)*xgap, c, img)
		plot(ypxl1+1, xpxl1, fpart(yend)*xgap, c, img)
	} else {
		plot(xpxl1, ypxl1, rfpart(yend)*xgap, c, img)
		plot(xpxl1, ypxl1+1, fpart(yend)*xgap, c, img)
	}
	intery := yend + gradient

	xend = math.Round(x1)
	yend = y1 + gradient*(xend-x1)
	xgap = fpart(x1 + 0.5)
	xpxl2 := xend
	ypxl2 := math.Floor(yend)
	if steep {
		plot(ypxl2, xpxl2, rfpart(yend)*xgap, c, img)
		plot(ypxl2+1, xpxl2, fpart(yend)*xgap, c, img)
	} else {
		plot(xpxl2, ypxl2, rfpart(yend)*xgap, c, img)
		plot(xpxl2, ypxl2+1, fpart(yend)*xgap, c, img)
	}

	if steep {
		for x := xpxl1 + 1; x <= xpxl2; x++ {
			//copy img
			gif.AppendImage(img)

			plot(math.Floor(intery), x, rfpart(intery), c, img)
			plot(math.Floor(intery)+1, x, fpart(intery), c, img)
			intery = intery + gradient
		}
	} else {
		for x := xpxl1 + 1; x <= xpxl2; x++ {
			//copy img
			gif.AppendImage(img)

			plot(x, math.Floor(intery), rfpart(intery), c, img)
			plot(x, math.Floor(intery)+1, fpart(intery), c, img)
			intery = intery + gradient
		}
	}

	gif.Save("out/bresenham-wu.gif")
}
