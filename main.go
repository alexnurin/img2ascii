package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/png"
	"os"
)

var (
	img_size_y = flag.Int("h", 0, "Image height")
	img_size_x = flag.Int("w", 0, "Image width")
	quadro     = flag.Bool("q", false, "Convert to quadro size")
	normalize  = flag.Bool("n", false, "Convert to normal size (3:2)")
	colors     = []string{"\033[0;31m%c\033[0m", "\033[0;32m%c\033[0m", "\033[0;34m%c\033[0m"}
)

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func getGray(c color.Color) int {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	return int(r)
}

func getColor(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(r), int(g), int(b)
}

func getChar(r int) rune {
	chars := []rune("@#%*+-:. ")
	id := r * len(chars) / 256
	return chars[id]
}

func processCell(img image.Image, y int, x int, sz_x int, sz_y int) (rune, int) {
	count := 0
	n, _ := minMax(x+sz_x, img.Bounds().Dx())
	m, _ := minMax(y+sz_y, img.Bounds().Dy())

	for i := x; i < n; i++ {
		for j := y; j < m; j++ {
			count += getGray(img.At(i, j))
		}
	}
	res := count / ((n - x) * (m - y))

	clr := []int{0, 0, 0}
	for i := x; i < n; i++ {
		for j := y; j < m; j++ {
			r, g, b := getColor(img.At(i, j))
			clr[0] += r
			clr[1] += g
			clr[2] += b
		}
	}
	mx := 0
	mx_i := 0
	for i := range clr {
		if clr[i] > mx {
			mx = clr[i]
			mx_i = i
		}
	}

	return getChar(res), mx_i
}

func minMax(a int, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func getDeltas(img image.Image) (int, int) {
	delta_x := 64
	delta_y := 64
	if *img_size_y > 0 {
		delta_y = *img_size_y
	}
	if *img_size_x > 0 {
		delta_x = *img_size_x
	}

	sz_x := img.Bounds().Dx() / delta_x
	sz_y := img.Bounds().Dy() / delta_y
	if *normalize {
		_, sz_x = minMax(sz_x, sz_y)
		_, sz_y = minMax(sz_x, sz_y)
		sz_y = sz_y * 4 / 3
	} else if *quadro {
		_, sz_x = minMax(sz_x, sz_y)
		_, sz_y = minMax(sz_x, sz_y)
	}
	return sz_x, sz_y
}

func convertToAscii(img image.Image) ([][]rune, [][]int) {
	sz_x, sz_y := img.Bounds().Dx(), img.Bounds().Dy()

	delta_x, delta_y := getDeltas(img)

	sz_y_new, sz_x_new := sz_y/delta_y+1, sz_x/delta_x+1

	textImg := make([][]rune, sz_y_new)
	colorImg := make([][]int, sz_y_new)

	for i := 0; i < sz_y_new; i++ {
		textImg[i] = make([]rune, sz_x_new)
		colorImg[i] = make([]int, sz_x_new)
	}

	for i := 0; i < sz_y; i += delta_y {
		for j := 0; j < sz_x; j += delta_x {
			textImg[i/delta_y][j/delta_x], colorImg[i/delta_y][j/delta_x] = processCell(img, i, j, delta_x, delta_y)
		}
	}
	return textImg, colorImg
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <imagename.jpg>")
		os.Exit(0)
	}
	img := flag.Arg(0)
	img_decoded, err := decodeImageFile(img)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	textImg, colorImg := convertToAscii(img_decoded)

	for i := range textImg {
		for j := range textImg[i] {
			fmt.Printf(colors[colorImg[i][j]], textImg[i][j])
		}
		fmt.Println()
	}
}
