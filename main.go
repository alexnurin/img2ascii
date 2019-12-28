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

func processCell(c color.Color) rune {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	chars := []rune("@#%*+-:. ")
	id := int(r) * len(chars) / 256
	return chars[id]
}

func getCell(img image.Image, y int, x int) rune {
	res :=  processCell(img.At(x, y))
	return res
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

	textImg := make([][]rune, sz_y)
	fmt.Println(sz_x, sz_y)

	for i := 0; i < sz_y; i++ {
		textImg[i] = make([]rune, sz_x)
	}

	for i := 0; i < sz_y; i++ {
		for j := 0; j < sz_x; j++ {
			textImg[i][j] = getCell(img, i, j)
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
	for i := range textImg {
		for j := range textImg[i] {
			fmt.Printf("%c", textImg[i][j])
		}
		fmt.Println()
	}
}
