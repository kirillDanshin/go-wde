package wde

import (
	"image"
	"image/draw"
)

// Window interface describes API
// for backend's Window struct
type Window interface {
	SetTitle(title string)
	SetSize(width, height int)
	Size() (width, height int)
	LockSize(lock bool)
	Show()
	Screen() (im Image)
	FlushImage(bounds ...image.Rectangle)
	EventChan() (events <-chan interface{})
	Close() (err error)
	SetCursor(cursor Cursor)
}

// Image interface describes API
// for backend's Image struct
type Image interface {
	draw.Image
	// CopyRGBA() copies the source image to this image, translating
	// the source image to the provided bounds.
	CopyRGBA(src *image.RGBA, bounds image.Rectangle)
}

/*
Run needs to be called in the main thread to make
your code as cross-platform as possible.

Some wde backends (cocoa) require that this function be called in the
main thread. To make your code as cross-platform as possible, it is
recommended that your main function look like the the code below.

	func main() {
		go theRestOfYourProgram()
		wde.Run()
	}

wde.Run() will return when wde.Stop() is called.

For this to work, you must import one of the go.wde backends. For
instance,

	import _ "github.com/kirillDanshin/go-wde/xgb"

or

	import _ "github.com/kirillDanshin/go-wde/win"

or

	import _ "github.com/kirillDanshin/go-wde/cocoa"


will register a backend with go.wde, allowing you to call
wde.Run(), wde.Stop() and wde.NewWindow() without referring to the
backend explicitly.

If you pupt the registration import in a separate file filtered for
the correct platform, your project will work on all three major
platforms without configuration.

That is, if you import go.wde/xgb in a file named "wde_linux.go",
go.wde/win in a file named "wde_windows.go" and go.wde/cocoa in a
file named "wde_darwin.go", the go tool will import the correct one.

*/
func Run() {
	BackendRun()
}

// BackendRun throws a panic if
// no backend were imported.
// It's overwritten by backends.
var BackendRun = func() {
	panic("no wde backend imported")
}

/*
Stop trigger. Call this when you want wde.Run() to return.
Usually to allow your program to exit gracefully.
*/
func Stop() {
	BackendStop()
}

// BackendStop throws a panic if
// no backend were imported.
// It's overwritten by backends.
var BackendStop = func() {
	panic("no wde backend imported")
}

/*
NewWindow creates a new window
with the specified width and height.
*/
func NewWindow(width, height int) (Window, error) {
	return BackendNewWindow(width, height)
}

// BackendNewWindow throws a panic if
// no backend were imported.
// It's overwritten by backends.
var BackendNewWindow = func(width, height int) (Window, error) {
	panic("no wde backend imported")
}
