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

func getChar(r int) rune {
	//r := getGray(c)
	chars := []rune("@#%*+-:. ")
	id := r * len(chars) / 256
	return chars[id]
}

func processCell(img image.Image, y int, x int, sz_x int, sz_y int) rune {
	count := 0
	n, _ := minMax(x+sz_x, img.Bounds().Dx())
	m, _ := minMax(y+sz_y, img.Bounds().Dy())

	for i := x; i < n; i++ {
		for j := y; j < m; j++ {
			count += getGray(img.At(i, j))
		}
	}
	res := count / ((n - x) * (m - y))
	return getChar(res)
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

	//fmt.Println(*img_size_x, *img_size_y, delta_x, delta_y)
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

func convertToAscii(img image.Image) [][]rune {
	sz_x, sz_y := img.Bounds().Dx(), img.Bounds().Dy()

	delta_x, delta_y := getDeltas(img)
	// fmt.Println("deltas:", delta_x, delta_y)
	sz_y_new, sz_x_new := sz_y/delta_y+1, sz_x/delta_x+1

	textImg := make([][]rune, sz_y_new)
	for i := 0; i < sz_y_new; i++ {
		textImg[i] = make([]rune, sz_x_new)
	}

	// fmt.Println(sz_x, sz_y, sz_x_new, sz_y_new)

	for i := 0; i < sz_y; i += delta_y {
		for j := 0; j < sz_x; j += delta_x {
			textImg[i/delta_y][j/delta_x] = processCell(img, i, j, delta_x, delta_y)
		}
	}
	//for i := range textImg {
	//	for j := range textImg[i] {
	//		textImg[i][j] = processCell(img.At(j, i))
	//	}
	//}
	return textImg
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
	fmt.Println("base:", img_decoded.Bounds().Dx(), img_decoded.Bounds().Dy())

	textImg := convertToAscii(img_decoded)
	//fmt.Println(textImg)

	fmt.Println("res:", len(textImg), len(textImg[0]))

	for i := range textImg {
		for j := range textImg[i] {
			fmt.Printf("%c", textImg[i][j])
		}
		fmt.Println()
	}

}
