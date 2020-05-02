package main

import (
	"flag"
	"fmt"
	Color "github.com/gookit/color"
	"image"
	"image/color"
	"sync"
	"time"

	//"time"

	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var (
	img_size_y = flag.Int("h", 0, "Image height")
	img_size_x = flag.Int("w", 0, "Image width")
	while_bg   = flag.Bool("bg", false, "Background color")
	text_color = flag.Bool("color", false, "Ascii text color")
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
	chars := []rune(".:+*%#@")
	id := r * len(chars) / 256
	return chars[id]
}

func processCell(img image.Image, y int, x int, sz_x int, sz_y int) (rune, Color.RGBStyle) {
	res := 0
	n, _ := minMax(x+sz_x, img.Bounds().Dx())
	m, _ := minMax(y+sz_y, img.Bounds().Dy())

	for i := x; i < n; i++ {
		for j := y; j < m; j++ {
			res += getGray(img.At(i, j))
		}
	}
	cnt := (n - x) * (m - y)
	res1 := res / cnt

	var rc, gc, bc int
	for i := x; i < n; i++ {
		for j := y; j < m; j++ {
			r, g, b := getColor(img.At(i, j))
			rc += r >> 8
			gc += g >> 8
			bc += b >> 8
		}
	}
	txt_color := Color.RGB(255, 255, 255)
	bg_color := Color.RGB(0, 0, 0)

	if *text_color {
		txt_color = Color.RGB(uint8(rc/cnt), uint8(gc/cnt), uint8(bc/cnt))
	}
	if *while_bg {
		bg_color = Color.RGB(255, 255, 255)
	}
	c := *Color.NewRGBStyle(txt_color, bg_color)
	return getChar(res1), c
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

	return sz_x, sz_y
}

func convertToAscii(img image.Image) ([][]rune, [][]Color.RGBStyle) {
	sz_x, sz_y := img.Bounds().Dx(), img.Bounds().Dy()

	delta_x, delta_y := getDeltas(img)

	sz_y_new, sz_x_new := sz_y/delta_y+1, sz_x/delta_x+1

	textImg := make([][]rune, sz_y_new)
	colorImg := make([][]Color.RGBStyle, sz_y_new)

	for i := 0; i < sz_y_new; i++ {
		textImg[i] = make([]rune, sz_x_new)
		colorImg[i] = make([]Color.RGBStyle, sz_x_new)
	}

	for i := 0; i < sz_y; i += delta_y {
		for j := 0; j < sz_x; j += delta_x {
			textImg[i/delta_y][j/delta_x], colorImg[i/delta_y][j/delta_x] = processCell(img, i, j, delta_x, delta_y)
		}
	}
	return textImg, colorImg
}

func call(a Color.RGBStyle, b rune, wg *sync.WaitGroup) {
	defer wg.Done()
	a.Printf("%c", b)

}

// TODO merge normal photo and sur photo to make more distinguishable
func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <imagename>")
		os.Exit(0)
	}
	img := flag.Arg(0)
	img_decoded, err := decodeImageFile(img)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	textImg, colorImg := convertToAscii(img_decoded)
	fmt.Println(len(textImg), len(textImg[0]))
	for i := range textImg {
		var wg sync.WaitGroup
		wg.Add(len(textImg[i]))
		for j := range textImg[i] {
			//colorImg[i][j].Printf("%c",textImg[i][j])
			go call(colorImg[i][j], textImg[i][j], &wg)
			// TODO убрать этот костыль
			time.Sleep(time.Second / 5000000)
		}
		wg.Wait()
		fmt.Println()
	}
}
