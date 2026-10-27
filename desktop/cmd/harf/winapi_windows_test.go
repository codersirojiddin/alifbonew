//go:build windows && amd64

package main

import (
	"testing"
	"unsafe"
)

// Incorrect native layouts compile but silently break Windows input APIs.
func TestWin64Layouts(t *testing.T){
	for name,c:=range map[string]struct{got,want uintptr}{
		"INPUT":{unsafe.Sizeof(input{}),40},
		"KEYBDINPUT":{unsafe.Sizeof(keyInput{}),24},
		"WNDCLASSEXW":{unsafe.Sizeof(windowClass{}),80},
		"MSG":{unsafe.Sizeof(message{}),48},
		"NOTIFYICONDATAW":{unsafe.Sizeof(notifyIcon{}),976},
		"GUITHREADINFO":{unsafe.Sizeof(guiInfo{}),72},
		"INPUT union":{unsafe.Offsetof(input{}.Data),8},
	}{if c.got!=c.want{t.Errorf("%s=%d, want %d",name,c.got,c.want)}}
}
func TestKeyInput(t *testing.T){
	in:=keyEvent(0,'ö',4)
	k:=(*keyInput)(unsafe.Pointer(&in.Data[0]))
	if in.Type!=1||k.Scan!='ö'||k.Flags!=4||k.Key!=0{t.Fatalf("wrong Unicode input: %+v",k)}
}
