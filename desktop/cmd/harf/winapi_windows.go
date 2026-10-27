//go:build windows && amd64

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

// Resolve system DLLs using absolute paths; no DLLs or runtimes are bundled.
func systemDLL(name string) *syscall.LazyDLL {
	root := os.Getenv("SystemRoot")
	if root == "" { panic("SystemRoot is not set") }
	return syscall.NewLazyDLL(filepath.Join(root, "System32", name))
}

var (
	user32 = systemDLL("user32.dll")
	kernel32 = systemDLL("kernel32.dll")
	shell32 = systemDLL("shell32.dll")
	advapi32 = systemDLL("advapi32.dll")
	gdi32 = systemDLL("gdi32.dll")
	messageBox = user32.NewProc("MessageBoxW")
	defWindowProc = user32.NewProc("DefWindowProcW")
	registerClass = user32.NewProc("RegisterClassExW")
	createWindow = user32.NewProc("CreateWindowExW")
	showWindow = user32.NewProc("ShowWindow")
	updateWindow = user32.NewProc("UpdateWindow")
	getMessage = user32.NewProc("GetMessageW")
	translateMessage = user32.NewProc("TranslateMessage")
	dispatchMessage = user32.NewProc("DispatchMessageW")
	postQuitMessage = user32.NewProc("PostQuitMessage")
	getModuleHandle = kernel32.NewProc("GetModuleHandleW")
	loadCursor = user32.NewProc("LoadCursorW")
	loadIcon = user32.NewProc("LoadIconW")
	getStockObject = gdi32.NewProc("GetStockObject")
	sendMessage = user32.NewProc("SendMessageW")
	sendMessageTimeout = user32.NewProc("SendMessageTimeoutW")
	setWindowText = user32.NewProc("SetWindowTextW")
	getWindowText = user32.NewProc("GetWindowTextW")
	getWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	setForeground = user32.NewProc("SetForegroundWindow")
	getForeground = user32.NewProc("GetForegroundWindow")
	getWindowThread = user32.NewProc("GetWindowThreadProcessId")
	getKeyboardLayout = user32.NewProc("GetKeyboardLayout")
	getGUIThreadInfo = user32.NewProc("GetGUIThreadInfo")
	registerHotKey = user32.NewProc("RegisterHotKey")
	unregisterHotKey = user32.NewProc("UnregisterHotKey")
	setHook = user32.NewProc("SetWindowsHookExW")
	unhook = user32.NewProc("UnhookWindowsHookEx")
	nextHook = user32.NewProc("CallNextHookEx")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	getKeyState = user32.NewProc("GetKeyState")
	sendInput = user32.NewProc("SendInput")
	setTimer = user32.NewProc("SetTimer")
	killTimer = user32.NewProc("KillTimer")
	destroyWindow = user32.NewProc("DestroyWindow")
	shellNotify = shell32.NewProc("Shell_NotifyIconW")
	createPopupMenu = user32.NewProc("CreatePopupMenu")
	appendMenu = user32.NewProc("AppendMenuW")
	trackPopupMenu = user32.NewProc("TrackPopupMenu")
	destroyMenu = user32.NewProc("DestroyMenu")
	getCursorPos = user32.NewProc("GetCursorPos")
	openProcess = kernel32.NewProc("OpenProcess")
	queryProcessName = kernel32.NewProc("QueryFullProcessImageNameW")
	closeHandle = kernel32.NewProc("CloseHandle")
	createMutex = kernel32.NewProc("CreateMutexW")
	openClipboard = user32.NewProc("OpenClipboard")
	closeClipboard = user32.NewProc("CloseClipboard")
	emptyClipboard = user32.NewProc("EmptyClipboard")
	setClipboard = user32.NewProc("SetClipboardData")
	globalAlloc = kernel32.NewProc("GlobalAlloc")
	globalLock = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	globalFree = kernel32.NewProc("GlobalFree")
	regCreate = advapi32.NewProc("RegCreateKeyExW")
	regOpen = advapi32.NewProc("RegOpenKeyExW")
	regClose = advapi32.NewProc("RegCloseKey")
	regQuery = advapi32.NewProc("RegQueryValueExW")
	regSet = advapi32.NewProc("RegSetValueExW")
	regDelete = advapi32.NewProc("RegDeleteValueW")
)

type point struct { X, Y int32 }
type rect struct { Left, Top, Right, Bottom int32 }
type windowClass struct {
	Size, Style uint32
	Proc uintptr
	ClassExtra, WindowExtra int32
	Instance, Icon, Cursor, Background uintptr
	MenuName, ClassName *uint16
	SmallIcon uintptr
}
type message struct {
	Window uintptr
	ID uint32
	WParam, LParam uintptr
	Time uint32
	Point point
	Private uint32
}
type keyboardData struct { Key, Scan, Flags, Time uint32; Extra uintptr }
type guiInfo struct { Size, Flags uint32; Active, Focus, Capture, MenuOwner, MoveSize, Caret uintptr; CaretRect rect }
type notifyIcon struct {
	Size uint32
	Window uintptr
	ID, Flags, Callback uint32
	Icon uintptr
	Tip [128]uint16
	State, StateMask uint32
	Info [256]uint16
	Timeout uint32
	Title [64]uint16
	InfoFlags uint32
	GUID [16]byte
	BalloonIcon uintptr
}
// INPUT includes a 32-byte union at offset 8 on Windows amd64.
type input struct { Type uint32; Padding uint32; Data [32]byte }
type keyInput struct { Key, Scan uint16; Flags, Time uint32; Extra uintptr }

func wide(s string) *uint16 { p, err := syscall.UTF16PtrFromString(s); if err != nil { panic(err) }; return p }
func dialog(text string, flags uintptr) uintptr { r, _, _ := messageBox.Call(0, uintptr(unsafe.Pointer(wide(text))), uintptr(unsafe.Pointer(wide("Harf — Go"))), flags); return r }
func setText(window uintptr, text string) { setWindowText.Call(window, uintptr(unsafe.Pointer(wide(text)))) }
func text(window uintptr) string {
	n, _, _ := getWindowTextLength.Call(window)
	if n > 16*1024*1024 { return "" }
	buffer := make([]uint16, n+1)
	getWindowText.Call(window, uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	return syscall.UTF16ToString(buffer)
}
func down(key uintptr) bool { r, _, _ := getAsyncKeyState.Call(key); return uint16(r)&0x8000 != 0 }
func keyEvent(key, scan uint16, flags uint32) input {
	in := input{Type: 1}
	*(*keyInput)(unsafe.Pointer(&in.Data[0])) = keyInput{Key: key, Scan: scan, Flags: flags}
	return in
}
func sendRune(r rune, erase bool) bool {
	keys := make([]input, 0, 6)
	if erase { keys = append(keys, keyEvent(8, 0, 0), keyEvent(8, 0, 2)) }
	for _, u := range utf16.Encode([]rune{r}) { keys = append(keys, keyEvent(0, u, 4), keyEvent(0, u, 6)) }
	n, _, _ := sendInput.Call(uintptr(len(keys)), uintptr(unsafe.Pointer(&keys[0])), unsafe.Sizeof(input{}))
	runtime.KeepAlive(keys)
	return n == uintptr(len(keys))
}

func clipboardCopy(owner uintptr, s string) bool {
	if s == "" { return false }
	ok, _, _ := openClipboard.Call(owner)
	if ok == 0 { return false }
	defer closeClipboard.Call()
	units := utf16.Encode([]rune(s + "\x00"))
	handle, _, _ := globalAlloc.Call(2, uintptr(len(units)*2))
	if handle == 0 { return false }
	buffer, _, _ := globalLock.Call(handle)
	if buffer == 0 { globalFree.Call(handle); return false }
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(buffer)), len(units)), units)
	globalUnlock.Call(handle)
	emptyClipboard.Call()
	ok, _, _ = setClipboard.Call(13, handle)
	if ok == 0 { globalFree.Call(handle); return false }
	return true // Windows owns handle after SetClipboardData succeeds.
}
