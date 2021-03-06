package glfw3

import (
	//	"fmt"

	"image"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	glfw "github.com/grd/glfw3"
	"github.com/kirillDanshin/go-wde"
)

func init() {
	wde.BackendNewWindow = func(width, height int) (w wde.Window, err error) {
		var window *Window
		window, err = NewWindow(width, height)
		w = window
		return
	}

	wde.BackendRun = glfw.Main

	wde.BackendStop = func() {
		glfw.Terminate()
		os.Exit(0)
	}

	go doRun()

	// Don't show the window context before using the Show function.
	go glfw.WindowHint(glfw.Visible, glfw.False)

	go flushBuffer()

	go setGlyph()
}

func doRun() {
	for {
		// Poll for and process events
		glfw.PollEvents()
	}
}

type Window struct {
	win    *glfw.Window
	buffer Image
	events chan interface{}
}

var windowMap = make(map[uintptr]*Window)

func NewWindow(width, height int) (w *Window, err error) {

	w = new(Window)

	w.win, err = glfw.CreateWindow(width, height, "", nil, nil)
	if err != nil {
		return nil, err
	}

	windowMap[w.win.C()] = w

	w.buffer.RGBA = image.NewRGBA(image.Rect(0, 0, width, height))

	// Events and callback functions for events

	w.events = make(chan interface{})
	w.win.SetMouseButtonCallback(mouseButtonCallback)
	w.win.SetCursorEnterCallback(cursorEnterCallback)
	w.win.SetCursorPositionCallback(cursorPositionCallback)
	w.win.SetFramebufferSizeCallback(framebufferSizeCallback)
	w.win.SetCharacterCallback(characterCallback)
	w.win.SetKeyCallback(keyCallback)
	w.checkShouldClose()

	return
}

func (w *Window) SetTitle(title string) {
	w.win.SetTitle(title)
}

func (w *Window) SetCursor(cursor wde.Cursor) {
	// TODO
}

func (w *Window) SetSize(width, height int) {
	w.win.SetSize(width, height)
}

func (w *Window) Size() (width, height int) {
	return w.win.GetSize()
}

func (w *Window) LockSize(lock bool) {
	// glfw supports only window size locking when the
	// parameter is set before the creation of the window.
	// It doesn't work on an existing window.
	if lock {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	}
}

func (w *Window) Show() {
	w.win.Show()
}

func (w *Window) Screen() wde.Image {
	return w.buffer
}

func (w *Window) FlushImage(bounds ...image.Rectangle) {

	if w.win.ShouldClose() {
		return
	}

	// TODO: Howto implement ...image.Rectangle

	// flush the buffer
	windowFlushBuffer <- w

	// waiting for the flushing is done before filling the buffer again
	<-windowFlushBufferDone
}

func (w *Window) Close() (err error) {
	w.win.Destroy()
	return
}

func (w *Window) checkShouldClose() {
	go func() {
		for {
			if w.win.ShouldClose() {
				var cev wde.CloseEvent
				w.events <- cev
				break
			}
		}
	}()
}

type Image struct {
	*image.RGBA
}

func (buffer Image) CopyRGBA(src *image.RGBA, r image.Rectangle) {
	// clip r against each image's bounds and move sp accordingly (see draw.clip())
	sp := image.ZP
	orig := r.Min
	r = r.Intersect(buffer.Bounds())
	r = r.Intersect(src.Bounds().Add(orig.Sub(sp)))
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	(sp).X += dx
	(sp).Y += dy

	i0 := (r.Min.X - buffer.Rect.Min.X) * 4
	i1 := (r.Max.X - buffer.Rect.Min.X) * 4
	si0 := (sp.X - src.Rect.Min.X) * 4
	yMax := r.Max.Y - buffer.Rect.Min.Y

	y := r.Min.Y - buffer.Rect.Min.Y
	sy := sp.Y - src.Rect.Min.Y
	for ; y != yMax; y, sy = y+1, sy+1 {
		dpix := buffer.Pix[y*buffer.Stride:]
		spix := src.Pix[sy*src.Stride:]

		for i, si := i0, si0; i < i1; i, si = i+4, si+4 {
			dpix[i+0] = spix[si+2]
			dpix[i+1] = spix[si+1]
			dpix[i+2] = spix[si+0]
			dpix[i+3] = spix[si+3]
		}
	}
}

var (
	windowFlushBuffer     = make(chan *Window)
	windowFlushBufferDone = make(chan bool)
)

func flushBuffer() {
	for {

		w := <-windowFlushBuffer

		w.win.MakeContextCurrent()

		imgWidth := w.buffer.Rect.Max.X
		imgHeight := w.buffer.Rect.Max.Y

		windowWidth, windowHeight := w.Size()

		// gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.RasterPos2f(-1, 1)
		gl.PixelZoom(1, -1)
		gl.Viewport(0, 0, windowWidth, windowHeight)
		gl.DrawPixels(imgWidth, imgHeight, gl.RGBA, gl.UNSIGNED_BYTE, &w.buffer.Pix[0])

		w.win.SwapBuffers()

		windowFlushBufferDone <- true
	}
}
