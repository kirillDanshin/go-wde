package win

import (
	"fmt"

	"github.com/AllenDang/w32"
	"github.com/kirillDanshin/go-wde"
)

func (wnd *Window) checkKeyState() {
	if !wnd.keysStale {
		return
	}
	keyboard := make([]byte, 256)
	if w32.GetKeyboardState(&keyboard) {
		/* virtual keys before 0x08 are mouse buttons; skip them */
		for vk := uintptr(0x08); vk <= 0xff; vk++ {
			isDown := keyboard[vk]&0x80 != 0
			key := keyFromVirtualKeyCode(vk)
			_, wasDown := wnd.keysDown[key]
			if isDown != wasDown {
				if isDown {
					wnd.keysDown[key] = true
					wnd.events <- wde.KeyDownEvent(wde.KeyEvent{key})
				} else {
					delete(wnd.keysDown, key)
					wnd.events <- wde.KeyUpEvent(wde.KeyEvent{key})
				}
			}
		}
		wnd.keysStale = false
	}
}

func (wnd *Window) constructChord() string {
	wnd.checkKeyState()
	return wde.ConstructChord(wnd.keysDown)
}

// virtKeyToWDE maps virtual key code to WDE key
// if this key is not exists, fallbacks to
// vk-0xff format
var virtKeyToWDE = map[uintptr]string{
	w32.VK_BACK:       wde.KeyBackspace,
	w32.VK_TAB:        wde.KeyTab,
	w32.VK_RETURN:     wde.KeyReturn,
	w32.VK_SHIFT:      wde.KeyLeftShift,
	w32.VK_CONTROL:    wde.KeyLeftControl,
	w32.VK_MENU:       wde.KeyLeftAlt,
	w32.VK_CAPITAL:    wde.KeyCapsLock,
	w32.VK_ESCAPE:     wde.KeyEscape,
	w32.VK_SPACE:      wde.KeySpace,
	w32.VK_PRIOR:      wde.KeyPrior,
	w32.VK_NEXT:       wde.KeyNext,
	w32.VK_END:        wde.KeyEnd,
	w32.VK_HOME:       wde.KeyHome,
	w32.VK_LEFT:       wde.KeyLeftArrow,
	w32.VK_UP:         wde.KeyUpArrow,
	w32.VK_RIGHT:      wde.KeyRightArrow,
	w32.VK_DOWN:       wde.KeyDownArrow,
	w32.VK_INSERT:     wde.KeyInsert,
	w32.VK_DELETE:     wde.KeyDelete,
	w32.VK_LWIN:       wde.KeyLeftSuper,
	w32.VK_RWIN:       wde.KeyRightSuper,
	w32.VK_NUMPAD0:    wde.Key0,
	w32.VK_NUMPAD1:    wde.Key1,
	w32.VK_NUMPAD2:    wde.Key2,
	w32.VK_NUMPAD3:    wde.Key3,
	w32.VK_NUMPAD4:    wde.Key4,
	w32.VK_NUMPAD5:    wde.Key5,
	w32.VK_NUMPAD6:    wde.Key6,
	w32.VK_NUMPAD7:    wde.Key7,
	w32.VK_NUMPAD8:    wde.Key8,
	w32.VK_NUMPAD9:    wde.Key9,
	w32.VK_MULTIPLY:   wde.KeyPadStar,
	w32.VK_ADD:        wde.KeyPadPlus,
	w32.VK_SUBTRACT:   wde.KeyPadMinus,
	w32.VK_DECIMAL:    wde.KeyPadDot,
	w32.VK_DIVIDE:     wde.KeyPadSlash,
	w32.VK_F1:         wde.KeyF1,
	w32.VK_F2:         wde.KeyF2,
	w32.VK_F3:         wde.KeyF3,
	w32.VK_F4:         wde.KeyF4,
	w32.VK_F5:         wde.KeyF5,
	w32.VK_F6:         wde.KeyF5,
	w32.VK_F7:         wde.KeyF7,
	w32.VK_F8:         wde.KeyF8,
	w32.VK_F9:         wde.KeyF9,
	w32.VK_F10:        wde.KeyF10,
	w32.VK_F11:        wde.KeyF11,
	w32.VK_F12:        wde.KeyF12,
	w32.VK_F13:        wde.KeyF13,
	w32.VK_F14:        wde.KeyF14,
	w32.VK_F15:        wde.KeyF15,
	w32.VK_F16:        wde.KeyF16,
	w32.VK_NUMLOCK:    wde.KeyNumlock,
	w32.VK_LSHIFT:     wde.KeyLeftShift,
	w32.VK_RSHIFT:     wde.KeyRightShift,
	w32.VK_LCONTROL:   wde.KeyLeftControl,
	w32.VK_RCONTROL:   wde.KeyRightControl,
	w32.VK_LMENU:      wde.KeyLeftAlt,
	w32.VK_RMENU:      wde.KeyRightAlt,
	w32.VK_OEM_1:      wde.KeySemicolon,
	w32.VK_OEM_PLUS:   wde.KeyEqual,
	w32.VK_OEM_COMMA:  wde.KeyComma,
	w32.VK_OEM_MINUS:  wde.KeyMinus,
	w32.VK_OEM_PERIOD: wde.KeyPeriod,
	w32.VK_OEM_2:      wde.KeySlash,
	w32.VK_OEM_3:      wde.KeyBackTick,
	w32.VK_OEM_4:      wde.KeyLeftBracket,
	w32.VK_OEM_5:      wde.KeyBackslash,
	w32.VK_OEM_6:      wde.KeyRightBracket,
	w32.VK_OEM_7:      wde.KeyQuote,

	// the rest lack wde constants. the first few are xgb compatible
	w32.VK_PAUSE:  "Pause",
	w32.VK_APPS:   "Menu",
	w32.VK_SCROLL: "Scroll_Lock",

	/*
		// the rest fallthrough to the default format "vk-0xff"
		w32.VK_LBUTTON:
		w32.VK_RBUTTON:
		w32.VK_CANCEL:
		w32.VK_MBUTTON:
		w32.VK_XBUTTON1:
		w32.VK_XBUTTON2:
		w32.VK_CLEAR:
		w32.VK_HANGUL:
		w32.VK_JUNJA:
		w32.VK_FINAL:
		w32.VK_KANJI:
		w32.VK_CONVERT:
		w32.VK_NONCONVERT:
		w32.VK_ACCEPT:
		w32.VK_MODECHANGE:
		w32.VK_SELECT:
		w32.VK_PRINT:
		w32.VK_EXECUTE:
		w32.VK_SNAPSHOT:
		w32.VK_HELP:
		w32.VK_SLEEP:
		w32.VK_SEPARATOR:
		w32.VK_F17:
		w32.VK_F18:
		w32.VK_F19:
		w32.VK_F20:
		w32.VK_F21:
		w32.VK_F22:
		w32.VK_F23:
		w32.VK_F24:
		w32.VK_BROWSER_BACK:
		w32.VK_BROWSER_FORWARD:
		w32.VK_BROWSER_REFRESH:
		w32.VK_BROWSER_STOP:
		w32.VK_BROWSER_SEARCH:
		w32.VK_BROWSER_FAVORITES:
		w32.VK_BROWSER_HOME:
		w32.VK_VOLUME_MUTE:
		w32.VK_VOLUME_DOWN:
		w32.VK_VOLUME_UP:
		w32.VK_MEDIA_NEXT_TRACK:
		w32.VK_MEDIA_PREV_TRACK:
		w32.VK_MEDIA_STOP:
		w32.VK_MEDIA_PLAY_PAUSE:
		w32.VK_LAUNCH_MAIL:
		w32.VK_LAUNCH_MEDIA_SELECT:
		w32.VK_LAUNCH_APP1:
		w32.VK_LAUNCH_APP2:
		w32.VK_OEM_8:
		w32.VK_OEM_AX:
		w32.VK_OEM_102:
		w32.VK_ICO_HELP:
		w32.VK_ICO_00:
		w32.VK_PROCESSKEY:
		w32.VK_ICO_CLEAR:
		w32.VK_OEM_RESET:
		w32.VK_OEM_JUMP:
		w32.VK_OEM_PA1:
		w32.VK_OEM_PA2:
		w32.VK_OEM_PA3:
		w32.VK_OEM_WSCTRL:
		w32.VK_OEM_CUSEL:
		w32.VK_OEM_ATTN:
		w32.VK_OEM_FINISH:
		w32.VK_OEM_COPY:
		w32.VK_OEM_AUTO:
		w32.VK_OEM_ENLW:
		w32.VK_OEM_BACKTAB:
		w32.VK_ATTN:
		w32.VK_CRSEL:
		w32.VK_EXSEL:
		w32.VK_EREOF:
		w32.VK_PLAY:
		w32.VK_ZOOM:
		w32.VK_NONAME:
		w32.VK_PA1:
		w32.VK_OEM_CLEAR:
	*/
}

func keyFromVirtualKeyCode(vk uintptr) string {
	if vk >= '0' && vk <= 'Z' {
		/* alphanumeric range (windows doesn't use 0x3a-0x40) */
		if vk >= 'A' {
			return fmt.Sprintf("%c", vk-'A'+'a') // convert to lower case
		}
		return fmt.Sprintf("%c", vk)
	}

	if virtKeyToWDE[vk] != "" {
		return virtKeyToWDE[vk]
	}

	// switch vk {
	// case w32.VK_BACK:
	// 	return wde.KeyBackspace
	// case w32.VK_TAB:
	// 	return wde.KeyTab
	// case w32.VK_RETURN:
	// 	return wde.KeyReturn
	// case w32.VK_SHIFT:
	// 	return wde.KeyLeftShift
	// case w32.VK_CONTROL:
	// 	return wde.KeyLeftControl
	// case w32.VK_MENU:
	// 	return wde.KeyLeftAlt
	// case w32.VK_CAPITAL:
	// 	return wde.KeyCapsLock
	// case w32.VK_ESCAPE:
	// 	return wde.KeyEscape
	// case w32.VK_SPACE:
	// 	return wde.KeySpace
	// case w32.VK_PRIOR:
	// 	return wde.KeyPrior
	// case w32.VK_NEXT:
	// 	return wde.KeyNext
	// case w32.VK_END:
	// 	return wde.KeyEnd
	// case w32.VK_HOME:
	// 	return wde.KeyHome
	// case w32.VK_LEFT:
	// 	return wde.KeyLeftArrow
	// case w32.VK_UP:
	// 	return wde.KeyUpArrow
	// case w32.VK_RIGHT:
	// 	return wde.KeyRightArrow
	// case w32.VK_DOWN:
	// 	return wde.KeyDownArrow
	// case w32.VK_INSERT:
	// 	return wde.KeyInsert
	// case w32.VK_DELETE:
	// 	return wde.KeyDelete
	// case w32.VK_LWIN:
	// 	return wde.KeyLeftSuper
	// case w32.VK_RWIN:
	// 	return wde.KeyRightSuper
	// case w32.VK_NUMPAD0:
	// 	return wde.Key0
	// case w32.VK_NUMPAD1:
	// 	return wde.Key1
	// case w32.VK_NUMPAD2:
	// 	return wde.Key2
	// case w32.VK_NUMPAD3:
	// 	return wde.Key3
	// case w32.VK_NUMPAD4:
	// 	return wde.Key4
	// case w32.VK_NUMPAD5:
	// 	return wde.Key5
	// case w32.VK_NUMPAD6:
	// 	return wde.Key6
	// case w32.VK_NUMPAD7:
	// 	return wde.Key7
	// case w32.VK_NUMPAD8:
	// 	return wde.Key8
	// case w32.VK_NUMPAD9:
	// 	return wde.Key9
	// case w32.VK_MULTIPLY:
	// 	return wde.KeyPadStar
	// case w32.VK_ADD:
	// 	return wde.KeyPadPlus
	// case w32.VK_SUBTRACT:
	// 	return wde.KeyPadMinus
	// case w32.VK_DECIMAL:
	// 	return wde.KeyPadDot
	// case w32.VK_DIVIDE:
	// 	return wde.KeyPadSlash
	// case w32.VK_F1:
	// 	return wde.KeyF1
	// case w32.VK_F2:
	// 	return wde.KeyF2
	// case w32.VK_F3:
	// 	return wde.KeyF3
	// case w32.VK_F4:
	// 	return wde.KeyF4
	// case w32.VK_F5:
	// 	return wde.KeyF5
	// case w32.VK_F6:
	// 	return wde.KeyF5
	// case w32.VK_F7:
	// 	return wde.KeyF7
	// case w32.VK_F8:
	// 	return wde.KeyF8
	// case w32.VK_F9:
	// 	return wde.KeyF9
	// case w32.VK_F10:
	// 	return wde.KeyF10
	// case w32.VK_F11:
	// 	return wde.KeyF11
	// case w32.VK_F12:
	// 	return wde.KeyF12
	// case w32.VK_F13:
	// 	return wde.KeyF13
	// case w32.VK_F14:
	// 	return wde.KeyF14
	// case w32.VK_F15:
	// 	return wde.KeyF15
	// case w32.VK_F16:
	// 	return wde.KeyF16
	// case w32.VK_NUMLOCK:
	// 	return wde.KeyNumlock
	// case w32.VK_LSHIFT:
	// 	return wde.KeyLeftShift
	// case w32.VK_RSHIFT:
	// 	return wde.KeyRightShift
	// case w32.VK_LCONTROL:
	// 	return wde.KeyLeftControl
	// case w32.VK_RCONTROL:
	// 	return wde.KeyRightControl
	// case w32.VK_LMENU:
	// 	return wde.KeyLeftAlt
	// case w32.VK_RMENU:
	// 	return wde.KeyRightAlt
	// case w32.VK_OEM_1:
	// 	return wde.KeySemicolon
	// case w32.VK_OEM_PLUS:
	// 	return wde.KeyEqual
	// case w32.VK_OEM_COMMA:
	// 	return wde.KeyComma
	// case w32.VK_OEM_MINUS:
	// 	return wde.KeyMinus
	// case w32.VK_OEM_PERIOD:
	// 	return wde.KeyPeriod
	// case w32.VK_OEM_2:
	// 	return wde.KeySlash
	// case w32.VK_OEM_3:
	// 	return wde.KeyBackTick
	// case w32.VK_OEM_4:
	// 	return wde.KeyLeftBracket
	// case w32.VK_OEM_5:
	// 	return wde.KeyBackslash
	// case w32.VK_OEM_6:
	// 	return wde.KeyRightBracket
	// case w32.VK_OEM_7:
	// 	return wde.KeyQuote
	//
	// // the rest lack wde constants. the first few are xgb compatible
	// case w32.VK_PAUSE:
	// 	return "Pause"
	// case w32.VK_APPS:
	// 	return "Menu"
	// case w32.VK_SCROLL:
	// 	return "Scroll_Lock"
	//
	// // the rest fallthrough to the default format "vk-0xff"
	// case w32.VK_LBUTTON:
	// case w32.VK_RBUTTON:
	// case w32.VK_CANCEL:
	// case w32.VK_MBUTTON:
	// case w32.VK_XBUTTON1:
	// case w32.VK_XBUTTON2:
	// case w32.VK_CLEAR:
	// case w32.VK_HANGUL:
	// case w32.VK_JUNJA:
	// case w32.VK_FINAL:
	// case w32.VK_KANJI:
	// case w32.VK_CONVERT:
	// case w32.VK_NONCONVERT:
	// case w32.VK_ACCEPT:
	// case w32.VK_MODECHANGE:
	// case w32.VK_SELECT:
	// case w32.VK_PRINT:
	// case w32.VK_EXECUTE:
	// case w32.VK_SNAPSHOT:
	// case w32.VK_HELP:
	// case w32.VK_SLEEP:
	// case w32.VK_SEPARATOR:
	// case w32.VK_F17:
	// case w32.VK_F18:
	// case w32.VK_F19:
	// case w32.VK_F20:
	// case w32.VK_F21:
	// case w32.VK_F22:
	// case w32.VK_F23:
	// case w32.VK_F24:
	// case w32.VK_BROWSER_BACK:
	// case w32.VK_BROWSER_FORWARD:
	// case w32.VK_BROWSER_REFRESH:
	// case w32.VK_BROWSER_STOP:
	// case w32.VK_BROWSER_SEARCH:
	// case w32.VK_BROWSER_FAVORITES:
	// case w32.VK_BROWSER_HOME:
	// case w32.VK_VOLUME_MUTE:
	// case w32.VK_VOLUME_DOWN:
	// case w32.VK_VOLUME_UP:
	// case w32.VK_MEDIA_NEXT_TRACK:
	// case w32.VK_MEDIA_PREV_TRACK:
	// case w32.VK_MEDIA_STOP:
	// case w32.VK_MEDIA_PLAY_PAUSE:
	// case w32.VK_LAUNCH_MAIL:
	// case w32.VK_LAUNCH_MEDIA_SELECT:
	// case w32.VK_LAUNCH_APP1:
	// case w32.VK_LAUNCH_APP2:
	// case w32.VK_OEM_8:
	// case w32.VK_OEM_AX:
	// case w32.VK_OEM_102:
	// case w32.VK_ICO_HELP:
	// case w32.VK_ICO_00:
	// case w32.VK_PROCESSKEY:
	// case w32.VK_ICO_CLEAR:
	// case w32.VK_OEM_RESET:
	// case w32.VK_OEM_JUMP:
	// case w32.VK_OEM_PA1:
	// case w32.VK_OEM_PA2:
	// case w32.VK_OEM_PA3:
	// case w32.VK_OEM_WSCTRL:
	// case w32.VK_OEM_CUSEL:
	// case w32.VK_OEM_ATTN:
	// case w32.VK_OEM_FINISH:
	// case w32.VK_OEM_COPY:
	// case w32.VK_OEM_AUTO:
	// case w32.VK_OEM_ENLW:
	// case w32.VK_OEM_BACKTAB:
	// case w32.VK_ATTN:
	// case w32.VK_CRSEL:
	// case w32.VK_EXSEL:
	// case w32.VK_EREOF:
	// case w32.VK_PLAY:
	// case w32.VK_ZOOM:
	// case w32.VK_NONAME:
	// case w32.VK_PA1:
	// case w32.VK_OEM_CLEAR:
	// }
	return fmt.Sprintf("vk-0x%02x", vk)
}
