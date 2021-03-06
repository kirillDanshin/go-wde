package xgb

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/kirillDanshin/go-wde"
)

func buttonForDetail(detail xproto.Button) wde.Button {
	switch detail {
	case 1:
		return wde.LeftButton
	case 2:
		return wde.MiddleButton
	case 3:
		return wde.RightButton
	case 4:
		return wde.WheelUpButton
	case 5:
		return wde.WheelDownButton
	}
	return 0
}

func (w *Window) handleEvents() {
	defer func() {
		recover()
		w.events <- wde.CloseEvent{}
		close(w.events)
	}()
	var noX int32 = 1<<31 - 1
	noX++
	var lastX, lastY int32 = noX, 0
	var button wde.Button

	downKeys := map[string]bool{}

	var (
		newWidth, newHeight = w.width, w.height
	)
	//
	var resize = func(*Window) {
		// w.buffer = xgraphics.New(w.xu, image.Rect(0, 0, newWidth, newHeight))
		if w.width == newWidth &&
			w.height == newHeight &&
			w.buffer.Bounds().Dx() == newWidth &&
			w.buffer.Bounds().Dy() == newHeight {
			return
		}
		var newBuf *xgraphics.Image
		newBuf = xgraphics.New(w.xu, image.Rect(0, 0, newWidth, newHeight))
		w.bufferLck.Lock()
		w.buffer.Destroy()
		w.buffer = newBuf
		w.bufferLck.Unlock()
		w.win.Resize(newWidth, newHeight)
		w.width, w.height = newWidth, newHeight
		newBuf = nil
		var re wde.ResizeEvent
		re.Width = newWidth
		re.Height = newHeight
		w.events <- re
		go log.Println("resized")
	}

	go func() {
		var e interface{}
		for {
			select {
			case e = <-w.events:
				switch e.(type) {
				case wde.CloseEvent:
					return
				default:
					continue
				}
			case <-time.After(1 * frameDuration):
				resize(w)
			}
		}
	}()

	for {
		e, err := w.conn.WaitForEvent()

		if err != nil {
			fmt.Fprintln(os.Stderr, "[go.wde X error] ", err)
			continue
		}

		switch e := e.(type) {

		case xproto.ButtonPressEvent:
			button = button | buttonForDetail(e.Detail)
			var bpe wde.MouseDownEvent
			bpe.Which = buttonForDetail(e.Detail)
			bpe.Where.X = int(e.EventX)
			bpe.Where.Y = int(e.EventY)
			lastX = int32(e.EventX)
			lastY = int32(e.EventY)
			w.events <- bpe

		case xproto.ButtonReleaseEvent:
			button = button & ^buttonForDetail(e.Detail)
			var bue wde.MouseUpEvent
			bue.Which = buttonForDetail(e.Detail)
			bue.Where.X = int(e.EventX)
			bue.Where.Y = int(e.EventY)
			lastX = int32(e.EventX)
			lastY = int32(e.EventY)
			w.events <- bue

		case xproto.LeaveNotifyEvent:
			var wee wde.MouseExitedEvent
			wee.Where.X = int(e.EventX)
			wee.Where.Y = int(e.EventY)
			if lastX != noX {
				wee.From.X = int(lastX)
				wee.From.Y = int(lastY)
			} else {
				wee.From.X = wee.Where.X
				wee.From.Y = wee.Where.Y
			}
			lastX = int32(e.EventX)
			lastY = int32(e.EventY)
			w.events <- wee
		case xproto.EnterNotifyEvent:
			var wee wde.MouseEnteredEvent
			wee.Where.X = int(e.EventX)
			wee.Where.Y = int(e.EventY)
			if lastX != noX {
				wee.From.X = int(lastX)
				wee.From.Y = int(lastY)
			} else {
				wee.From.X = wee.Where.X
				wee.From.Y = wee.Where.Y
			}
			lastX = int32(e.EventX)
			lastY = int32(e.EventY)
			w.events <- wee

		case xproto.MotionNotifyEvent:
			var mme wde.MouseMovedEvent
			mme.Where.X = int(e.EventX)
			mme.Where.Y = int(e.EventY)
			if lastX != noX {
				mme.From.X = int(lastX)
				mme.From.Y = int(lastY)
			} else {
				mme.From.X = mme.Where.X
				mme.From.Y = mme.Where.Y
			}
			lastX = int32(e.EventX)
			lastY = int32(e.EventY)
			if button == 0 {
				w.events <- mme
			} else {
				var mde wde.MouseDraggedEvent
				mde.MouseMovedEvent = mme
				mde.Which = button
				w.events <- mde
			}

		case xproto.KeyPressEvent:
			var ke wde.KeyEvent
			code := keybind.LookupString(w.xu, e.State, e.Detail)
			ke.Key = keyForCode(code)
			w.events <- wde.KeyDownEvent(ke)
			downKeys[ke.Key] = true
			kpe := wde.KeyTypedEvent{
				KeyEvent: ke,
				Glyph:    letterForCode(code),
				Chord:    wde.ConstructChord(downKeys),
			}
			w.events <- kpe

		case xproto.KeyReleaseEvent:
			var ke wde.KeyUpEvent
			ke.Key = keyForCode(keybind.LookupString(w.xu, e.State, e.Detail))
			delete(downKeys, ke.Key)
			w.events <- ke

		case xproto.KeymapNotifyEvent:
			newDownKeys := make(map[string]bool)
			for i := 0; i < len(e.Keys); i++ {
				mask := e.Keys[i]
				for j := 0; j < 8; j++ {
					if mask&(1<<uint(j)) != 0 {
						keycode := xproto.Keycode(8*(i+1) + j)
						key := keyForCode(keybind.LookupString(w.xu, 0, keycode))
						newDownKeys[key] = true
					}
				}
			}
			/* remove keys that are no longer pressed */
			for key := range downKeys {
				if _, ok := newDownKeys[key]; !ok {
					var ke wde.KeyUpEvent
					ke.Key = key
					delete(downKeys, key)
					w.events <- ke
				}
			}
			/* add keys that are newly pressed */
			for key := range newDownKeys {
				if _, ok := downKeys[key]; !ok {
					var ke wde.KeyDownEvent
					ke.Key = key
					downKeys[key] = true
					w.events <- ke
				}
			}

		case xproto.ConfigureNotifyEvent:
			newWidth = int(e.Width)
			newHeight = int(e.Height)

		case xproto.ClientMessageEvent:
			if icccm.IsDeleteProtocol(
				w.xu,
				// it's uses unkeyed fields
				// but it's ok for now.
				// it you know how to do it better
				// please feel free to send a PR
				xevent.ClientMessageEvent{&e},
			) {
				w.events <- wde.CloseEvent{}
			}
		case xproto.DestroyNotifyEvent:
		case xproto.ReparentNotifyEvent:
		case xproto.MapNotifyEvent:
		case xproto.UnmapNotifyEvent:
		case xproto.PropertyNotifyEvent:

		default:
			fmt.Fprintf(os.Stderr, "unhandled event: type %T\n%+v\n", e, e)
		}

	}
}

// EventChan returns the chan
func (w *Window) EventChan() (events <-chan interface{}) {
	return w.events
}
