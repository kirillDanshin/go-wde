package main

import (
	"fmt"
	"image/color"
	"image/gif"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/kirillDanshin/go-wde"
	_ "github.com/kirillDanshin/go-wde/init"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go wdetest()
	wde.Run()

	println("done")
}

func wdetest() {
	var wg sync.WaitGroup

	size := 300

	// emerland color
	green := color.RGBA{
		R: 46,
		G: 204,
		B: 113,
		A: 0,
	}

	// belize hole color
	blue := color.RGBA{
		R: 41,
		G: 128,
		B: 185,
		A: 0,
	}

	// carrot color
	yellow := color.RGBA{
		R: 230,
		G: 126,
		B: 34,
		A: 70,
	}

	amethyst := color.RGBA{
		R: 155,
		G: 89,
		B: 182,
		A: 70,
	}

	image, err := gif.Decode(gordon_gif())

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	ib := image.Bounds()
	imSize := ib.Size()
	imWidth := imSize.X
	imHeight := imSize.Y

	x := func() {
		// offset := time.Duration(8 * time.Millisecond)
		offset := time.Duration(0)

		dw, err := wde.NewWindow(size, size)
		if err != nil {
			fmt.Println(err)
			return
		}
		dw.SetTitle("hi!")
		dw.SetSize(size, size)

		events := dw.EventChan()

		done := make(chan bool)

		go func() {
		loop:
			for ei := range events {
				runtime.Gosched()
				switch e := ei.(type) {
				case wde.MouseDownEvent:
					fmt.Println("clicked", e.Where.X, e.Where.Y, e.Which)
					// dw.Close()
					// break loop
				case wde.MouseUpEvent:
				case wde.MouseMovedEvent:
				case wde.MouseDraggedEvent:
				case wde.MouseEnteredEvent:
					// fmt.Println("mouse entered", e.Where.X, e.Where.Y)
				case wde.MouseExitedEvent:
					// fmt.Println("mouse exited", e.Where.X, e.Where.Y)
				case wde.MagnifyEvent:
					fmt.Println("magnify", e.Where, e.Magnification)
				case wde.RotateEvent:
					fmt.Println("rotate", e.Where, e.Rotation)
				case wde.ScrollEvent:
					fmt.Println("scroll", e.Where, e.Delta)
				case wde.KeyDownEvent:
					// fmt.Println("KeyDownEvent", e.Glyph)
				case wde.KeyUpEvent:
					// fmt.Println("KeyUpEvent", e.Glyph)
				case wde.KeyTypedEvent:
					fmt.Printf("typed key %v, glyph %v, chord %v\n", e.Key, e.Glyph, e.Chord)
					switch e.Glyph {
					case "1":
						dw.SetCursor(wde.NormalCursor)
					case "2":
						dw.SetCursor(wde.CrosshairCursor)
					case "3":
						dw.SetCursor(wde.GrabHoverCursor)
					}
				case wde.CloseEvent:
					fmt.Println("close")
					dw.Close()
					break loop
				case wde.ResizeEvent:
					// i++
					// go fmt.Println("resize", e.Width, e.Height)
				}
			}
			done <- true
			fmt.Println("end of events")
		}()

		for i := 0; ; i++ {
			width, height := dw.Size()
			s := dw.Screen()
			// for x := 0; x < width; x++ {
			// 	for y := 0; y < height; y++ {
			// 		s.Set(x, y, color.White)
			// 	}
			// }
			var pxColor color.RGBA

			imLeft := width/2 - imWidth/2
			imRight := width/2 + imWidth/2
			imTop := height/2 - imHeight/2
			imBottom := height/2 + imHeight/2

			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {

					// top left
					if x < width/2 && y <= height/2 {
						pxColor = green
					}

					// top right
					if x > width/2 && y <= height/2 {
						pxColor = yellow
					}

					// bottom left
					if x < width/2 && y >= height/2 {
						pxColor = amethyst
					}

					// bottom right
					if x > width/2 && y >= height/2 {
						pxColor = blue
					}

					s.Set(x, y, pxColor)

					// if we're in the image part, draw the image
					if x > imLeft && x < imRight && y > imTop && y < imBottom {
						imPxColor := image.At(x-imLeft, y-imTop)
						if _, _, _, a := imPxColor.RGBA(); a != 0 {
							s.Set(
								x,
								y,
								imPxColor,
							)
						}
					}

				}
			}
			dw.FlushImage()
			dw.Show()
			select {
			case <-time.After(offset):
				dw.FlushImage()
			case <-done:
				wg.Done()
				return
			}
		}
	}
	wg.Add(1)
	go x()
	// wg.Add(1)
	// go x()

	wg.Wait()
	wde.Stop()
}
