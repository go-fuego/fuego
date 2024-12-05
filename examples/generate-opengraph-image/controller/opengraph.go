package controller

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	Title         = "Ewen's Notes"
	SubtitleColor = color.RGBA{200, 200, 200, 0xff}
	TitleColor    = color.RGBA{255, 255, 255, 0xff}
	LineColor     = color.RGBA{0x3D, 0xAE, 0xE3, 0xff}

	Font     *truetype.Font = nil
	FontFile                = "./Raleway-Regular.ttf"

	width  = 400
	height = 200
	lineY  = 120
	startX = 32
)

func OpenGraphHandler(w http.ResponseWriter, r *http.Request) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case lineY < y && y < lineY+4 && startX <= x && x < 140:
				img.Set(x, y, LineColor)
			default:
				img.Set(x, y, color.Black)
			}
		}
	}
	addLabel(img, startX, lineY-12, 36, r.PathValue("title"), TitleColor)
	addLabel(img, startX, lineY+32, 24, Title, SubtitleColor)

	png.Encode(w, img)
}

func addLabel(img *image.RGBA, x, y int, size float64, label string, color color.RGBA) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	if Font == nil {
		fontBytes, err := os.ReadFile(FontFile)
		if err != nil {
			log.Println(err)
			return
		}
		f, err := truetype.Parse(fontBytes)
		if err != nil {
			log.Println(err)
			return
		}
		Font = f
	}

	ffont := truetype.NewFace(Font, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: ffont,
		Dot:  point,
	}

	for d.MeasureString(label).Ceil() > width-(2*startX) || size <= 0 {
		size--
		d.Face = truetype.NewFace(Font, &truetype.Options{
			Size:    size,
			DPI:     72,
			Hinting: font.HintingNone,
		})
	}
	d.DrawString(label)
}
