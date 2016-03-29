package xgb

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/kirillDanshin/go-wde"
)

var cursorCache map[wde.Cursor]xproto.Cursor
var cursorXIds map[wde.Cursor]uint16

func init() {
	cursorCache = make(map[wde.Cursor]xproto.Cursor)
	// the default cursor is always cursor 0 - no need to CreateCursor so cache it up front
	cursorCache[wde.NormalCursor] = 0

	cursorXIds = map[wde.Cursor]uint16{
		wde.ResizeNCursor:    xcursor.TopSide,
		wde.ResizeECursor:    xcursor.RightSide,
		wde.ResizeSCursor:    xcursor.BottomSide,
		wde.ResizeWCursor:    xcursor.LeftSide,
		wde.ResizeEWCursor:   xcursor.SBHDoubleArrow,
		wde.ResizeNSCursor:   xcursor.SBVDoubleArrow,
		wde.ResizeNECursor:   xcursor.TopRightCorner,
		wde.ResizeSECursor:   xcursor.BottomRightCorner,
		wde.ResizeSWCursor:   xcursor.BottomLeftCorner,
		wde.ResizeNWCursor:   xcursor.TopLeftCorner,
		wde.CrosshairCursor:  xcursor.Crosshair,
		wde.IBeamCursor:      xcursor.XTerm,
		wde.GrabHoverCursor:  xcursor.Hand2,
		wde.GrabActiveCursor: xcursor.Hand2,
		// xcursor defines this but no crossed-circle or similar. GUMBY. dafuq?
		wde.NotAllowedCursor: xcursor.Gumby,
	}
}

// SetCursor sets cursor
func (w *Window) SetCursor(cursor wde.Cursor) {
	if w.cursor != cursor {
		w.cursor = cursor
		w.win.Change(xproto.CwCursor, uint32(xCursor(w, cursor)))
	}
}

func xCursor(w *Window, c wde.Cursor) xproto.Cursor {
	xc, ok := cursorCache[c]
	if !ok {
		xc = createCursor(w, c)
		cursorCache[c] = xc
	}
	return xc
}

func createCursor(w *Window, c wde.Cursor) xproto.Cursor {
	xid, ok := cursorXIds[c]
	if ok {
		xc, err := xcursor.CreateCursor(w.win.X, xid)
		if err == nil {
			return xc
		}
	}
	return 0 // fallback to cursor 0
}
