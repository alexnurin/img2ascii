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

func processPixel(r int) rune {
	//r := getGray(c)
	chars := []rune("@#%*+-:. ")
	id := r * len(chars) / 256
	return chars[id]
}

func processCell(img image.Image, y int, x int, sz_x int, sz_y int) rune {
	count := 0
	for i := x; i < x+sz_x; i++ {
		for j := y; j < y+sz_y; j++ {
			count += getGray(img.At(i, j))
		}
	}
	res := count / (sz_x * sz_y)
	return processPixel(res)
}

func convertToAscii(img image.Image) [][]rune {
	sz_x := img.Bounds().Dx()
	sz_y := img.Bounds().Dy()

	if sz_x > *img_size_x && *img_size_x > 0 {
		sz_x = *img_size_x
	}
	if sz_y > *img_size_y && *img_size_y > 0 {
		sz_y = *img_size_y
	}

	delta_x := 10
	delta_y := 15
	sz_y_new := sz_y/delta_y + 1
	sz_x_new := sz_x/delta_x + 1

	textImg := make([][]rune, sz_y_new)
	for i := 0; i < sz_y_new; i++ {
		textImg[i] = make([]rune, sz_x_new)
	}

	//fmt.Println(sz_x, sz_y, sz_x_new, sz_y_new)

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

	textImg := convertToAscii(img_decoded)
	//fmt.Println(textImg)
	for i := range textImg {
		for j := range textImg[i] {
			fmt.Printf("%c", textImg[i][j])
		}
		fmt.Println()
	}

}
