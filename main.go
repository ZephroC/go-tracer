package main
import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"github.com/faiface/pixel/imdraw"
)
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Iain World!",
		Bounds: pixel.R(0, 0, 1280, 720),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(color.RGBA{0,0,0,0})
	canvas := pixelgl.NewCanvas(win.Bounds())
	var size = win.Bounds().Size()
	var width = int(size.X)
	var height = int(size.Y)
	var stride = 4
	buffer := make([]uint8,width * height * stride)
	DrawToBuffer(buffer,width,height,stride)
	imd := imdraw.New(nil)
	for !win.Closed() {
		imd.Clear()
		canvas.SetPixels(buffer)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}
func main() {
	pixelgl.Run(run)
}