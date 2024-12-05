package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/freetype/truetype"

	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/middleware/cache"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

// A custom option to add a custom response to the OpenAPI spec.
// The route returns a PNG image.
var optionReturnsPNG = func(br *fuego.BaseRoute) {
	response := openapi3.NewResponse()
	response.WithDescription("Generated image")
	response.WithContent(openapi3.NewContentWithSchema(nil, []string{"image/png"}))
	br.Operation.AddResponse(200, response)
}

func main() {
	s := fuego.NewServer(
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			PrettyFormatJson: true,
		}),
	)

	fuego.GetStd(s, "/{title}", imageGen,
		optionReturnsPNG,
		option.Description("Generate an image with a title. Useful for Opengraph."),
		option.Path("title", "The title to write on the image", param.Example("example", "My awesome article!")),
		option.Middleware(cache.New()),
	)

	s.Run()
}

var (
	darkGray = color.RGBA{50, 50, 50, 0xff}
	red      = color.RGBA{0xE3, 0x42, 0x34, 0xff}
	yellow   = color.RGBA{0xFF, 0xBA, 0x08, 0xff}
	green    = color.RGBA{0x84, 0xBD, 0x00, 0xff}
	blue     = color.RGBA{0x3D, 0xAE, 0xE3, 0xff}
	blue2    = color.RGBA{0x1D, 0x99, 0xF3, 0xff}
	width    = 400
	height   = 200
	lineY    = 120
	startX   = 40
)

func imageGen(w http.ResponseWriter, r *http.Request) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case lineY+5 < y && y < lineY+10 && startX <= x && x < 60:
				img.Set(x, y, red)
			case lineY+5 < y && y < lineY+10 && 60 <= x && x < 80:
				img.Set(x, y, yellow)
			case lineY+5 < y && y < lineY+10 && 80 <= x && x < 100:
				img.Set(x, y, green)
			case lineY+5 < y && y < lineY+10 && 100 <= x && x < 120:
				img.Set(x, y, blue)
			case lineY+5 < y && y < lineY+10 && 120 <= x && x < 140:
				img.Set(x, y, blue2)
			default:
				img.Set(x, y, color.White)
			}
		}
	}
	addLabel(img, startX, lineY, 36, r.PathValue("title"))

	png.Encode(w, img)
}

func addLabel(img *image.RGBA, x, y int, size float64, label string) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	fontBytes, err := os.ReadFile("./Raleway-Regular.ttf")
	if err != nil {
		log.Println(err)
		return
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	ffont := truetype.NewFace(f, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(darkGray),
		Face: ffont,
		Dot:  point,
	}
	d.DrawString(label)
}
