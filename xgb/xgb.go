package xgb

import (
	"fmt"
	"image"
	"log"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/kirillDanshin/go-wde"
)

const frameDuration = 16 * time.Millisecond

func init() {
	wde.BackendNewWindow = func(width, height int) (w wde.Window, err error) {
		w, err = NewWindow(width, height)
		return
	}
	ch := make(chan struct{})
	wde.BackendRun = func() {
		<-ch
	}
	wde.BackendStop = func() {
		ch <- struct{}{}
	}
}

// AllEventsMask is a mask for Events
// that used for listening all events
const AllEventsMask = xproto.EventMaskKeyPress |
	xproto.EventMaskKeyRelease |
	xproto.EventMaskKeymapState |
	xproto.EventMaskButtonPress |
	xproto.EventMaskButtonRelease |
	xproto.EventMaskEnterWindow |
	xproto.EventMaskLeaveWindow |
	xproto.EventMaskPointerMotion |
	xproto.EventMaskStructureNotify

// Window struct
type Window struct {
	win           *xwindow.Window
	xu            *xgbutil.XUtil
	conn          *xgb.Conn
	buffer        *xgraphics.Image
	bufferLck     *sync.Mutex
	width, height int
	lockedSize    bool
	closed        bool
	showed        bool
	cursor        wde.Cursor // most recently set cursor

	events chan interface{}
}

// NewWindow creates a new window with provided width and height
func NewWindow(width, height int) (w *Window, err error) {

	w = new(Window)
	w.width, w.height = width, height
	w.showed = false

	w.xu, err = xgbutil.NewConn()
	if err != nil {
		return
	}

	w.conn = w.xu.Conn()
	// screen := w.xu.Screen()

	w.win, err = xwindow.Generate(w.xu)
	if err != nil {
		log.Printf("ERROR: %#+v\n", err)
		return
	}

	keybind.Initialize(w.xu)

	err = w.win.CreateChecked(w.xu.RootWin(), 0, 0, width, height, xproto.CwBackPixel, 0x606060ff)
	if err != nil {
		return
	}

	w.win.Listen(AllEventsMask)

	err = icccm.WmProtocolsSet(w.xu, w.win.Id, []string{"WM_DELETE_WINDOW"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		err = nil
	}

	w.bufferLck = &sync.Mutex{}
	w.buffer = xgraphics.New(w.xu, image.Rect(0, 0, width, height))
	w.buffer.XSurfaceSet(w.win.Id)

	// I /think/ XDraw actually sends data to server?
	w.buffer.XDraw()
	// I /think/ XPaint tells the server to paint image to window
	w.buffer.XPaint(w.win.Id)

	keyMap, modMap := keybind.MapsGet(w.xu)
	keybind.KeyMapSet(w.xu, keyMap)
	keybind.ModMapSet(w.xu, modMap)

	w.events = make(chan interface{})

	w.SetIcon(Gordon)
	w.SetIconName("Go")

	go w.handleEvents()

	return
}

// SetTitle sets Window title
func (w *Window) SetTitle(title string) {
	if w.closed {
		return
	}
	err := ewmh.WmNameSet(w.xu, w.win.Id, title)
	if err != nil {
		// TODO: log
	}
	return
}

// SetSize sets Window size
func (w *Window) SetSize(width, height int) {
	if w.closed {
		return
	}

	w.width, w.height = width, height
	if w.lockedSize {
		w.updateSizeHints()
	}
	w.win.Resize(width, height)
	return
}

// Size returns Window width and height
func (w *Window) Size() (width, height int) {
	if w.closed {
		return
	}
	width, height = w.width, w.height
	return
}

// LockSize locks Window size
func (w *Window) LockSize(lock bool) {
	w.lockedSize = lock
	w.updateSizeHints()
}

func (w *Window) updateSizeHints() {
	hints := new(icccm.NormalHints)
	if w.lockedSize {
		hints.Flags = icccm.SizeHintPMinSize | icccm.SizeHintPMaxSize
		hints.MinWidth = uint(w.width)
		hints.MaxWidth = uint(w.width)
		hints.MinHeight = uint(w.height)
		hints.MaxHeight = uint(w.height)
	}
	icccm.WmNormalHintsSet(w.xu, w.win.Id, hints)
}

// Show the window
func (w *Window) Show() {
	if w.closed || w.showed {
		return
	}
	w.showed = true
	w.win.Map()
}

// Screen returns wde.Image of screen
func (w *Window) Screen() (im wde.Image) {
	if w.closed {
		return
	}
	im = &Image{w.buffer}
	return im
}

// FlushImage and draw a new one
func (w *Window) FlushImage(bounds ...image.Rectangle) {

	if w.closed {
		return
	}
	if w.buffer.Pixmap == 0 {
		w.bufferLck.Lock()
		if err := w.buffer.XSurfaceSet(w.win.Id); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		w.bufferLck.Unlock()
	}
	if len(bounds) > 0 {
		w.buffer.XPaintRects(w.win.Id, bounds...)
	} else {
		w.buffer.XDraw()
		w.buffer.XPaint(w.win.Id)
	}
}

// Close the window
func (w *Window) Close() (err error) {
	if w.closed {
		return
	}
	w.win.Destroy()
	w.closed = true
	return
}

// Image Just an image.
type Image struct {
	*xgraphics.Image
}

// CopyRGBA of image
func (buffer Image) CopyRGBA(src *image.RGBA, r image.Rectangle) {
	// clip r against each image's bounds and move sp accordingly (see draw.clip())
	sp := src.Bounds().Min
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
