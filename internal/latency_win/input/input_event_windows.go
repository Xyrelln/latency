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
		"backspace": types.VK_BACK,   // BACKSPACE key
		"tab":       types.VK_TAB,    // TAB key
		"enter":     types.VK_RETURN, // ENTER key
		// "":          types.VK_CLEAR,   // CLEAR key
		// "shift":     types.VK_SHIFT,   // SHIFT key
		// "ctrl":      types.VK_CONTROL, // CTRL key
		// "alt":       types.VK_MENU,    // ALT key
		"pause":    types.VK_PAUSE,   // PAUSE key
		"capslock": types.VK_CAPITAL, // CAPS LOCK key
		// "": types.VK_KANA                , // IME Kana mode
		// "": types.VK_HANGUEL             , // IME Hanguel mode (maintained for compatibility; use VK_HANGUL)
		// "": types.VK_HANGUL              , // IME Hangul mode
		// "": types.VK_IME_ON              , // IME On
		// "": types.VK_JUNJA               , // IME Junja mode
		// "": types.VK_FINAL               , // IME final mode
		// "": types.VK_HANJA               , // IME Hanja mode
		// "": types.VK_KANJI               , // IME Kanji mode
		// "": types.VK_IME_OFF             , // IME Off
		"esc": types.VK_ESCAPE, // ESC key
		// "": types.VK_CONVERT             , // IME convert
		// "": types.VK_NONCONVERT          , // IME nonconvert
		// "": types.VK_ACCEPT              , // IME accept
		// "": types.VK_MODECHANGE          , // IME mode change request
		"space":    types.VK_SPACE, // SPACEBAR
		"pageup":   types.VK_PRIOR, // PAGE UP key
		"pagedown": types.VK_NEXT,  // PAGE DOWN key
		"end":      types.VK_END,   // END key
		"home":     types.VK_HOME,  // HOME key
		"left":     types.VK_LEFT,  // LEFT ARROW key
		"up":       types.VK_UP,    // UP ARROW key
		"right":    types.VK_RIGHT, // RIGHT ARROW key
		"down":     types.VK_DOWN,  // DOWN ARROW key
		// "": types.VK_SELECT              , // SELECT key
		// "": types.VK_PRINT               , // PRINT key
		// "": types.VK_EXECUTE             , // EXECUTE key
		"prtsc": types.VK_SNAPSHOT, // PRINT SCREEN key
		"ins":   types.VK_INSERT,   // INS key
		"del":   types.VK_DELETE,   // DEL key
		"help":  types.VK_HELP,     // HELP key
		"0":     types.VK_0,        // 0 key
		"1":     types.VK_1,        // 1 key
		"2":     types.VK_2,        // 2 key
		"3":     types.VK_3,        // 3 key
		"4":     types.VK_4,        // 4 key
		"5":     types.VK_5,        // 5 key
		"6":     types.VK_6,        // 6 key
		"7":     types.VK_7,        // 7 key
		"8":     types.VK_8,        // 8 key
		"9":     types.VK_9,        // 9 key
		"a":     types.VK_A,        // A key
		"b":     types.VK_B,        // B key
		"c":     types.VK_C,        // C key
		"d":     types.VK_D,        // D key
		"e":     types.VK_E,        // E key
		"f":     types.VK_F,        // F key
		"g":     types.VK_G,        // G key
		"h":     types.VK_H,        // H key
		"i":     types.VK_I,        // I key
		"j":     types.VK_J,        // J key
		"k":     types.VK_K,        // K key
		"l":     types.VK_L,        // L key
		"m":     types.VK_M,        // M key
		"n":     types.VK_N,        // N key
		"o":     types.VK_O,        // O key
		"p":     types.VK_P,        // P key
		"q":     types.VK_Q,        // Q key
		"r":     types.VK_R,        // R key
		"s":     types.VK_S,        // S key
		"t":     types.VK_T,        // T key
		"u":     types.VK_U,        // U key
		"v":     types.VK_V,        // V key
		"w":     types.VK_W,        // W key
		"x":     types.VK_X,        // X key
		"y":     types.VK_Y,        // Y key
		"z":     types.VK_Z,        // Z key
		"lwin":  types.VK_LWIN,     // Left Windows key (Natural keyboard)
		"rwin":  types.VK_RWIN,     // Right Windows key (Natural keyboard)
		// "": types.VK_APPS                , // Applications key (Natural keyboard)
		// "": types.VK_SLEEP               , // Computer Sleep key
		"num0": types.VK_NUMPAD0,  // Numeric keypad 0 key
		"num1": types.VK_NUMPAD1,  // Numeric keypad 1 key
		"num2": types.VK_NUMPAD2,  // Numeric keypad 2 key
		"num3": types.VK_NUMPAD3,  // Numeric keypad 3 key
		"num4": types.VK_NUMPAD4,  // Numeric keypad 4 key
		"num5": types.VK_NUMPAD5,  // Numeric keypad 5 key
		"num6": types.VK_NUMPAD6,  // Numeric keypad 6 key
		"num7": types.VK_NUMPAD7,  // Numeric keypad 7 key
		"num8": types.VK_NUMPAD8,  // Numeric keypad 8 key
		"num9": types.VK_NUMPAD9,  // Numeric keypad 9 key
		"num*": types.VK_MULTIPLY, // Multiply key
		"num+": types.VK_ADD,      // Add key
		"num-": types.VK_SUBTRACT, // Subtract key
		"num.": types.VK_DECIMAL,  // Decimal key
		"num/": types.VK_DIVIDE,   // Divide key
		"F1":   types.VK_F1,       // F1 key
		"F2":   types.VK_F2,       // F2 key
		"F3":   types.VK_F3,       // F3 key
		"F4":   types.VK_F4,       // F4 key
		"F5":   types.VK_F5,       // F5 key
		"F6":   types.VK_F6,       // F6 key
		"F7":   types.VK_F7,       // F7 key
		"F8":   types.VK_F8,       // F8 key
		"F9":   types.VK_F9,       // F9 key
		"F10":  types.VK_F10,      // F10 key
		"F11":  types.VK_F11,      // F11 key
		"F12":  types.VK_F12,      // F12 key
		"F13":  types.VK_F13,      // F13 key
		"F14":  types.VK_F14,      // F14 key
		"F15":  types.VK_F15,      // F15 key
		"F16":  types.VK_F16,      // F16 key
		"F17":  types.VK_F17,      // F17 key
		"F18":  types.VK_F18,      // F18 key
		"F19":  types.VK_F19,      // F19 key
		"F20":  types.VK_F20,      // F20 key
		"F21":  types.VK_F21,      // F21 key
		"F22":  types.VK_F22,      // F22 key
		"F23":  types.VK_F23,      // F23 key
		"F24":  types.VK_F24,      // F24 key

		"numlock": types.VK_NUMLOCK,  // NUM LOCK key
		"scrlock": types.VK_SCROLL,   // SCROLL LOCK key
		"lshift":  types.VK_LSHIFT,   // Left SHIFT key
		"rshift":  types.VK_RSHIFT,   // Right SHIFT key
		"lctrl":   types.VK_LCONTROL, // Left CONTROL key
		"rctrl":   types.VK_RCONTROL, // Right CONTROL key
		"lalt":    types.VK_LMENU,    // Left MENU key
		"ralt":    types.VK_RMENU,    // Right MENU key

		"+": types.VK_OEM_PLUS,   // For any country/region, the '+' key
		",": types.VK_OEM_COMMA,  // For any country/region, the ',' key
		"-": types.VK_OEM_MINUS,  // For any country/region, the '-' key
		".": types.VK_OEM_PERIOD, // For any country/region, the '.' key
		";": types.VK_OEM_1,      // For any country/region, the '.' key
		"/": types.VK_OEM_2,      // For any country/region, the '.' key

		"`":  types.VK_OEM_3, // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '`~' key
		"[":  types.VK_OEM_4, // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '[{' key
		"\\": types.VK_OEM_5, // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '\|' key
		"]":  types.VK_OEM_6, // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key
		"'":  types.VK_OEM_7, // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the 'single-quote/double-quote' key
	}

	keyCodeList = []string{
		"backspace", // BACKSPACE key
		"tab",       // TAB key
		"enter",     // ENTER key
		"pause",     // PAUSE key
		"capslock",  // CAPS LOCK key
		"esc",       // ESC key
		"space",     // SPACEBAR
		"pageup",    // PAGE UP key
		"pagedown",  // PAGE DOWN key
		"end",       // END key
		"home",      // HOME key
		"left",      // LEFT ARROW key
		"up",        // UP ARROW key
		"right",     // RIGHT ARROW key
		"down",      // DOWN ARROW key
		"prtsc",     // PRINT SCREEN key
		"ins",       // INS key
		"del",       // DEL key
		"help",      // HELP key
		"0",         // 0 key
		"1",         // 1 key
		"2",         // 2 key
		"3",         // 3 key
		"4",         // 4 key
		"5",         // 5 key
		"6",         // 6 key
		"7",         // 7 key
		"8",         // 8 key
		"9",         // 9 key
		"a",         // A key
		"b",         // B key
		"c",         // C key
		"d",         // D key
		"e",         // E key
		"f",         // F key
		"g",         // G key
		"h",         // H key
		"i",         // I key
		"j",         // J key
		"k",         // K key
		"l",         // L key
		"m",         // M key
		"n",         // N key
		"o",         // O key
		"p",         // P key
		"q",         // Q key
		"r",         // R key
		"s",         // S key
		"t",         // T key
		"u",         // U key
		"v",         // V key
		"w",         // W key
		"x",         // X key
		"y",         // Y key
		"z",         // Z key
		"lwin",      // Left Windows key (Natural keyboard)
		"rwin",      // Right Windows key (Natural keyboard)
		"num0",      // Numeric keypad 0 key
		"num1",      // Numeric keypad 1 key
		"num2",      // Numeric keypad 2 key
		"num3",      // Numeric keypad 3 key
		"num4",      // Numeric keypad 4 key
		"num5",      // Numeric keypad 5 key
		"num6",      // Numeric keypad 6 key
		"num7",      // Numeric keypad 7 key
		"num8",      // Numeric keypad 8 key
		"num9",      // Numeric keypad 9 key
		"num*",      // Multiply key
		"num+",      // Add key
		"num-",      // Subtract key
		"num.",      // Decimal key
		"num/",      // Divide key
		"F1",        // F1 key
		"F2",        // F2 key
		"F3",        // F3 key
		"F4",        // F4 key
		"F5",        // F5 key
		"F6",        // F6 key
		"F7",        // F7 key
		"F8",        // F8 key
		"F9",        // F9 key
		"F10",       // F10 key
		"F11",       // F11 key
		"F12",       // F12 key
		"F13",       // F13 key
		"F14",       // F14 key
		"F15",       // F15 key
		"F16",       // F16 key
		"F17",       // F17 key
		"F18",       // F18 key
		"F19",       // F19 key
		"F20",       // F20 key
		"F21",       // F21 key
		"F22",       // F22 key
		"F23",       // F23 key
		"F24",       // F24 key
		"numlock",   // NUM LOCK key
		"scrlock",   // SCROLL LOCK key
		"lshift",    // Left SHIFT key
		"rshift",    // Right SHIFT key
		"lctrl",     // Left CONTROL key
		"rctrl",     // Right CONTROL key
		"lalt",      // Left MENU key
		"ralt",      // Right MENU key
		"+",         // For any country/region, the '+' key
		",",         // For any country/region, the ',' key
		"-",         // For any country/region, the '-' key
		".",         // For any country/region, the '.' key
		";",         // For any country/region, the '.' key
		"/",         // For any country/region, the '.' key
		"`",         // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '`~' key
		"[",         // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '[{' key
		"\\",        // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the '\|' key
		"]",         // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key
		"'",         // Used for miscellaneous characters; it can vary by keyboard.  For the US standard keyboard, the 'single-quote/double-quote' key
	}
)

// KeyToVKCode ...
func KeyToVKCode(k string) (types.VKCode, error) {
	if code, ok := keyCodeMap[k]; ok {
		return code, nil
	}
	return 0, fmt.Errorf("invalid key")
}

// KeyList ...
func KeyList() []string {
	return keyCodeList
}
