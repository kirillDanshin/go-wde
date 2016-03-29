package cocoa

// #include "cursor.h"
import "C"

import (
	"unsafe"

	"github.com/kirillDanshin/go-wde"
)

var cursors map[wde.Cursor]unsafe.Pointer

func init() {
	C.initMacCursor()

	cursors = map[wde.Cursor]unsafe.Pointer{
		wde.NoneCursor:     nil,
		wde.NormalCursor:   C.cursors.arrow,
		wde.ResizeNCursor:  C.cursors.resizeUp,
		wde.ResizeECursor:  C.cursors.resizeRight,
		wde.ResizeSCursor:  C.cursors.resizeDown,
		wde.ResizeWCursor:  C.cursors.resizeLeft,
		wde.ResizeEWCursor: C.cursors.resizeLeftRight,
		wde.ResizeNSCursor: C.cursors.resizeUpDown,

		// might be able to improve the diagonal arrow cursors:
		// http://stackoverflow.com/questions/10733228/native-osx-lion-resize-cursor-for-custom-nswindow-or-nsview
		wde.ResizeNECursor: C.cursors.pointingHand,
		wde.ResizeSECursor: C.cursors.pointingHand,
		wde.ResizeSWCursor: C.cursors.pointingHand,
		wde.ResizeNWCursor: C.cursors.pointingHand,

		wde.CrosshairCursor:  C.cursors.crosshair,
		wde.IBeamCursor:      C.cursors.IBeam,
		wde.GrabHoverCursor:  C.cursors.openHand,
		wde.GrabActiveCursor: C.cursors.closedHand,
		wde.NotAllowedCursor: C.cursors.operationNotAllowed,
	}
}

func setCursor(c wde.Cursor) {
	nscursor := cursors[c]
	if nscursor != nil {
		C.setCursor(nscursor)
	}
}

func (w *Window) SetCursor(cursor wde.Cursor) {
	if w.cursor == cursor {
		return
	}
	if w.hasMouse {
		/* the osx set cursor is application wide rather than window specific */
		setCursor(cursor)
	}
	w.cursor = cursor
}
