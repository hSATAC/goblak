package main

import (
	"image"
	"image/color"
	_ "image/jpeg" // Register JPEG format
	"image/png"    // Register PNG  format
	"log"
	"net/http"
	"os"
)

// Converted implements image.Image, so you can
// pretend that it is the converted image.
type Converted struct {
	Img image.Image
	Mod color.Model
}

// We return the new color model...
func (c *Converted) ColorModel() color.Model {
	return c.Mod
}

// ... but the original bounds
func (c *Converted) Bounds() image.Rectangle {
	return c.Img.Bounds()
}

// At forwards the call to the original image and
// then asks the color model to convert it.
func (c *Converted) At(x, y int) color.Color {
	return c.Mod.Convert(c.Img.At(x, y))
}

func main() {
	http.Handle("/gray", http.HandlerFunc(convertGray))
	http.Handle("/bw", http.HandlerFunc(convertBw))
	err := http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func convertGray(w http.ResponseWriter, req *http.Request) {
	img := imageFromURL(req.FormValue("url"))

	gr := &Converted{img, color.GrayModel}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, gr)
}

func convertBw(w http.ResponseWriter, req *http.Request) {
	img := imageFromURL(req.FormValue("url"))

	bw := []color.Color{color.Black, color.White}
	gr := &Converted{img, color.Palette(bw)}

	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, gr)
}

func imageFromURL(imgUrl string) image.Image {
	if imgUrl == "" {
		return nil // Handle empty parameter here
	}

	res, err := http.Get(imgUrl)
	if err != nil || res.StatusCode != 200 {
		log.Fatalln(err)
	}

	img, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return img
}
