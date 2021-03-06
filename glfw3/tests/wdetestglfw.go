// +build ignore

package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/kirillDanshin/go-wde"
	_ "github.com/kirillDanshin/go-wde/glfw3"
)

func main() {
	go wdetest()
	wde.Run()

	println("done")
}

func wdetest() {
	var wg sync.WaitGroup

	size := 200

	x := func() {
		offset := time.Duration(rand.Intn(1e9))

		dw, err := wde.NewWindow(size, size)
		if err != nil {
			fmt.Println(err)
			return
		}
		dw.SetTitle("hi GLFW!")
		dw.SetSize(size, size)
		dw.Show()

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
					fmt.Println("mouse entered", e.Where.X, e.Where.Y)
				case wde.MouseExitedEvent:
					fmt.Println("mouse exited", e.Where.X, e.Where.Y)
				case wde.KeyDownEvent:
					// fmt.Println("KeyDownEvent", e.Key)
				case wde.KeyUpEvent:
					// fmt.Println("KeyUpEvent", e.Key)
				case wde.KeyTypedEvent:
					fmt.Printf("typed key %v, glyph %v chord %v\n", e.Key, e.Glyph, e.Chord)
				case wde.CloseEvent:
					fmt.Println("close")
					dw.Close()
					break loop
				case wde.ResizeEvent:
					fmt.Println("resize", e.Width, e.Height)
				}
			}
			done <- true
			fmt.Println("end of events")
		}()

		for i := 0; ; i++ {
			width, height := dw.Size()
			s := dw.Screen()
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					s.Set(x, y, color.White)
				}
			}
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					var r uint8
					if x > width/2 {
						r = 255
					}
					var g uint8
					if y >= height/2 {
						g = 255
					}
					var b uint8
					if y < height/4 || y >= height*3/4 {
						b = 255
					}
					if i%2 == 1 {
						r = 255 - r
					}

					if y > height-10 {
						r = 255
						g = 255
						b = 255
					}

					if x == y {
						r = 100
						g = 100
						b = 100
					}

					s.Set(x, y, color.RGBA{r, g, b, 255})
				}
			}
			dw.FlushImage()
			select {
			case <-time.After(5e8 + offset):
			case <-done:
				wg.Done()
				return
			}
		}
	}
	wg.Add(1)
	go x()
	wg.Add(1)
	go x()

	wg.Wait()
	wde.Stop()
}
