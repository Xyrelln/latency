//go:build windows
// +build windows

package input

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"golang.org/x/sys/windows"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

// WindowsInputEv ...
type WindowsInputEv struct {
}

// const mouse event
const (
	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MOUSEWHEEL  = 0x020A
	WM_MOUSEHWHEEL = 0x020E
)

var bootTime time.Time

func init() {
	Modkernel32 := windows.NewLazySystemDLL("kernel32.dll")
	procGetTickCount64 := Modkernel32.NewProc("GetTickCount64")

	procGetTickCount := procGetTickCount64
	err := procGetTickCount64.Find()
	if err != nil {
		panic(err)
	}
	r1, _, lastErr := syscall.Syscall(procGetTickCount.Addr(), 0, 0, 0, 0)
	if lastErr != 0 {
		panic(lastErr)
	}
	bootTime = time.Now().Add(-time.Duration(r1) * time.Millisecond)
}

// WaitMouseEvent ...
func (wev WindowsInputEv) WaitMouseEvent(ctx context.Context, mouseEvent int) (eventTime time.Time, err error) {
	mouseChan := make(chan types.MouseEvent, 100)
	err = mouse.Install(nil, mouseChan)
	if err != nil {
		return
	}
	defer mouse.Uninstall()
	// timeoutC := time.After(timeout)
	for {
		select {
		case mev := <-mouseChan:
			if mev.Message == types.Message(mouseEvent) {
				fmt.Printf("mouse event time: %d\n", mev.Time)
				eventTime = bootTime.Add(time.Duration(mev.Time) * time.Millisecond)
				return
			}
			continue
		case <-ctx.Done():
			err = fmt.Errorf("time out")
			return
		}
	}
}

// WaitKeyBoardEvent ...
func (wev WindowsInputEv) WaitKeyBoardEvent(ctx context.Context, keyCode types.VKCode) (eventTime time.Time, err error) {
	keyboardChan := make(chan types.KeyboardEvent, 100)
	err = keyboard.Install(nil, keyboardChan)
	if err != nil {
		return
	}
	defer keyboard.Uninstall()
	for {
		select {
		case kev := <-keyboardChan:
			if kev.VKCode == keyCode {
				eventTime = bootTime.Add(time.Duration(kev.Time) * time.Millisecond)
				return
			}
			continue
		case <-ctx.Done():
			err = fmt.Errorf("time out")
			return
		}
	}
}

var (
	keyCodeMap = map[string]types.VKCode{
		// "": types.VK_LBUTTON             , // Left mouse button
		// "": types.VK_RBUTTON             , // Right mouse button
		// "": types.VK_CANCEL              , // Control-break processing
		// "": types.VK_MBUTTON             , // Middle mouse button (three-button mouse)
		// "": types.VK_XBUTTON1            , // X1 mouse button
		// "": types.VK_XBUTTON2            , // X2 mouse button
		"Backspace": types.VK_BACK,   // BACKSPACE key
		"Tab":       types.VK_TAB,    // TAB key
		"Enter":     types.VK_RETURN, // ENTER key
		// "":          types.VK_CLEAR,   // CLEAR key
		// "shift":     types.VK_SHIFT,   // SHIFT key
		// "ctrl":      types.VK_CONTROL, // CTRL key
		// "alt":       types.VK_MENU,    // ALT key
		"Pause":    types.VK_PAUSE,   // PAUSE key
		"CapsLock": types.VK_CAPITAL, // CAPS LOCK key
		// "": types.VK_KANA                , // IME Kana mode
		// "": types.VK_HANGUEL             , // IME Hanguel mode (maintained for compatibility; use VK_HANGUL)
		// "": types.VK_HANGUL              , // IME Hangul mode
		// "": types.VK_IME_ON              , // IME On
		// "": types.VK_JUNJA               , // IME Junja mode
		// "": types.VK_FINAL               , // IME final mode
		// "": types.VK_HANJA               , // IME Hanja mode
		// "": types.VK_KANJI               , // IME Kanji mode
		// "": types.VK_IME_OFF             , // IME Off
		"Escape": types.VK_ESCAPE, // ESC key
		// "": types.VK_CONVERT             , // IME convert
		// "": types.VK_NONCONVERT          , // IME nonconvert
		// "": types.VK_ACCEPT              , // IME accept
		// "": types.VK_MODECHANGE          , // IME mode change request
		"Space":    types.VK_SPACE, // SPACEBAR
		"PageUp":   types.VK_PRIOR, // PAGE UP key
		"PageDown": types.VK_NEXT,  // PAGE DOWN key
		"End":      types.VK_END,   // END key
		"Home":     types.VK_HOME,  // HOME key

		"ArrowLeft":  types.VK_LEFT,  // LEFT ARROW key
		"ArrowUp":    types.VK_UP,    // UP ARROW key
		"ArrowRight": types.VK_RIGHT, // RIGHT ARROW key
		"ArrowDown":  types.VK_DOWN,  // DOWN ARROW key
		// "": types.VK_SELECT              , // SELECT key
		// "": types.VK_PRINT               , // PRINT key
		// "": types.VK_EXECUTE             , // EXECUTE key
		"PrintScreen": types.VK_SNAPSHOT, // PRINT SCREEN key
		"Insert":      types.VK_INSERT,   // INS key
		"Delete":      types.VK_DELETE,   // DEL key
		"Help":        types.VK_HELP,     // HELP key

		"Digit0": types.VK_0, // 0 key
		"Digit1": types.VK_1, // 1 key
		"Digit2": types.VK_2, // 2 key
		"Digit3": types.VK_3, // 3 key
		"Digit4": types.VK_4, // 4 key
		"Digit5": types.VK_5, // 5 key
		"Digit6": types.VK_6, // 6 key
		"Digit7": types.VK_7, // 7 key
		"Digit8": types.VK_8, // 8 key
		"Digit9": types.VK_9, // 9 key
		"KeyA":   types.VK_A, // A key
		"KeyB":   types.VK_B, // B key
		"KeyC":   types.VK_C, // C key
		"KeyD":   types.VK_D, // D key
		"KeyE":   types.VK_E, // E key
		"KeyF":   types.VK_F, // F key
		"KeyG":   types.VK_G, // G key
		"KeyH":   types.VK_H, // H key
		"KeyI":   types.VK_I, // I key
		"KeyJ":   types.VK_J, // J key
		"KeyK":   types.VK_K, // K key
		"KeyL":   types.VK_L, // L key
		"KeyM":   types.VK_M, // M key
		"KeyN":   types.VK_N, // N key
		"KeyO":   types.VK_O, // O key
		"KeyP":   types.VK_P, // P key
		"KeyQ":   types.VK_Q, // Q key
		"KeyR":   types.VK_R, // R key
		"KeyS":   types.VK_S, // S key
		"KeyT":   types.VK_T, // T key
		"KeyU":   types.VK_U, // U key
		"KeyV":   types.VK_V, // V key
		"KeyW":   types.VK_W, // W key
		"KeyX":   types.VK_X, // X key
		"KeyY":   types.VK_Y, // Y key
		"KeyZ":   types.VK_Z, // Z key

		// "lwin":  types.VK_LWIN,     // Left Windows key (Natural keyboard)
		// "rwin":  types.VK_RWIN,     // Right Windows key (Natural keyboard)

		// "": types.VK_APPS                , // Applications key (Natural keyboard)
		// "": types.VK_SLEEP               , // Computer Sleep key
		"Numpad0": types.VK_NUMPAD0, // Numeric keypad 0 key
		"Numpad1": types.VK_NUMPAD1, // Numeric keypad 1 key
		"Numpad2": types.VK_NUMPAD2, // Numeric keypad 2 key
		"Numpad3": types.VK_NUMPAD3, // Numeric keypad 3 key
		"Numpad4": types.VK_NUMPAD4, // Numeric keypad 4 key
		"Numpad5": types.VK_NUMPAD5, // Numeric keypad 5 key
		"Numpad6": types.VK_NUMPAD6, // Numeric keypad 6 key
		"Numpad7": types.VK_NUMPAD7, // Numeric keypad 7 key
		"Numpad8": types.VK_NUMPAD8, // Numeric keypad 8 key
		"Numpad9": types.VK_NUMPAD9, // Numeric keypad 9 key

		"NumpadMultiply": types.VK_MULTIPLY, // Multiply key
		"NumpadAdd":      types.VK_ADD,      // Add key
		"NumpadSubtract": types.VK_SUBTRACT, // Subtract key
		"NumpadDecimal":  types.VK_DECIMAL,  // Decimal key
		"NumpadDivide":   types.VK_DIVIDE,   // Divide key

		"F1":  types.VK_F1,  // F1 key
		"F2":  types.VK_F2,  // F2 key
		"F3":  types.VK_F3,  // F3 key
		"F4":  types.VK_F4,  // F4 key
		"F5":  types.VK_F5,  // F5 key
		"F6":  types.VK_F6,  // F6 key
		"F7":  types.VK_F7,  // F7 key
		"F8":  types.VK_F8,  // F8 key
		"F9":  types.VK_F9,  // F9 key
		"F10": types.VK_F10, // F10 key
		"F11": types.VK_F11, // F11 key
		"F12": types.VK_F12, // F12 key
		"F13": types.VK_F13, // F13 key
		"F14": types.VK_F14, // F14 key
		"F15": types.VK_F15, // F15 key
		"F16": types.VK_F16, // F16 key
		"F17": types.VK_F17, // F17 key
		"F18": types.VK_F18, // F18 key
		"F19": types.VK_F19, // F19 key
		"F20": types.VK_F20, // F20 key
		"F21": types.VK_F21, // F21 key
		"F22": types.VK_F22, // F22 key
		"F23": types.VK_F23, // F23 key
		"F24": types.VK_F24, // F24 key

		"NumLock":      types.VK_NUMLOCK,  // NUM LOCK key
		"ScrollLock":   types.VK_SCROLL,   // SCROLL LOCK key
		"ShiftLeft":    types.VK_LSHIFT,   // Left SHIFT key
		"ShiftRight":   types.VK_RSHIFT,   // Right SHIFT key
		"ControlLeft":  types.VK_LCONTROL, // Left CONTROL key
		"ControlRight": types.VK_RCONTROL, // Right CONTROL key
		"AltLeft":      types.VK_LMENU,    // Left MENU key
		"AltRight":     types.VK_RMENU,    // Right MENU key

		"Equal":        types.VK_OEM_PLUS,   // For any country/region, the '+' key
		"Comma":        types.VK_OEM_COMMA,  // For any country/region, the ',' key
		"Minus":        types.VK_OEM_MINUS,  // For any country/region, the '-' key
		"Period":       types.VK_OEM_PERIOD, // For any country/region, the '.' key
		"Semicolon":    types.VK_OEM_1,      // For any country/region, the ';' key
		"Slash":        types.VK_OEM_2,      // For any country/region, the '/?' key
		"Backquote":    types.VK_OEM_3,      // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '`~' key
		"BracketLeft":  types.VK_OEM_4,      // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '[{' key
		"Backslash":    types.VK_OEM_5,      // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '\|' key
		"BracketRight": types.VK_OEM_6,      // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key
		"Quote":        types.VK_OEM_7,      // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the 'single-quote/double-quote' key
	}
)

// KeyToVKCode ...
func KeyToVKCode(k string) (types.VKCode, error) {
	if code, ok := keyCodeMap[k]; ok {
		return code, nil
	}
	return 0, fmt.Errorf("invalid key")
}
