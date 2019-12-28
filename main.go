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

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func processPixel(c color.Color) rune {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	chars := []rune("@#%*+-:. ")
	//fmt.Println(r)
	id := int(r) * len(chars) / 256
	return chars[id]
}

func convertToAscii(img image.Image) [][]rune {
	textImg := make([][]rune, img.Bounds().Dy())
	// i := 0; i < i.Bound().Dy; i++
	for i := range textImg {
		textImg[i] = make([]rune, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
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
