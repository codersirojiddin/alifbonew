//go:build windows && amd64

package main

import (
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unsafe"
)

var excluded = map[string]bool{"code.exe":true,"code-insiders.exe":true,"devenv.exe":true,"idea64.exe":true,"pycharm64.exe":true,"webstorm64.exe":true,"goland64.exe":true,"rider64.exe":true,"notepad++.exe":true,"windowsterminal.exe":true,"powershell.exe":true,"pwsh.exe":true,"cmd.exe":true,"mintty.exe":true,"putty.exe":true}

func focusedControl(window uintptr) uintptr {
	thread,_,_:=getWindowThread.Call(window,0)
	info:=guiInfo{Size:uint32(unsafe.Sizeof(guiInfo{}))}
	if ok,_,_:=getGUIThreadInfo.Call(thread,uintptr(unsafe.Pointer(&info)));ok!=0{return info.Focus}
	return 0
}
func (a *application) blocked(window,focus uintptr) bool {
	var pid uint32
	thread,_,_:=getWindowThread.Call(window,uintptr(unsafe.Pointer(&pid)))
	layout,_,_:=getKeyboardLayout.Call(thread)
	if uint16(layout)!=0x0409{return true} // Only translate US physical VK positions.
	if a.processWindow!=window {
		a.processWindow=window;a.processBlocked=true
		process,_,_:=openProcess.Call(0x1000,0,uintptr(pid))
		if process!=0 {
			buffer:=make([]uint16,1024);size:=uint32(len(buffer))
			ok,_,_:=queryProcessName.Call(process,0,uintptr(unsafe.Pointer(&buffer[0])),uintptr(unsafe.Pointer(&size)))
			closeHandle.Call(process)
			if ok!=0{a.processBlocked=excluded[strings.ToLower(filepath.Base(syscall.UTF16ToString(buffer[:size])))]}
		}
	}
	if a.processBlocked{return true}
	if focus!=0 {
		var passwordChar uintptr
		ok,_,_:=sendMessageTimeout.Call(focus,0xD2,0,0,2,20,uintptr(unsafe.Pointer(&passwordChar)))
		if ok==0||passwordChar!=0{return true}
	}
	return false
}

func keyboardProc(code,wp,lp uintptr) uintptr {
	a:=current
	if int32(code)>=0 && a!=nil && a.auto && (wp==0x100||wp==0x104) {
		key:=(*keyboardData)(unsafe.Pointer(lp))
		if key.Flags&0x10==0 { // Ignore injected events; keep only one candidate rune.
			window,_,_:=getForeground.Call();focus:=focusedControl(window)
			if a.previousWindow!=window||a.previousFocus!=focus||time.Since(a.previousTime)>2500*time.Millisecond {a.previous=0}
			a.previousWindow=window;a.previousFocus=focus;a.previousTime=time.Now()
			if down(0x11)||down(0x12)||down(0x5B)||down(0x5C)||a.blocked(window,focus){a.previous=0}else {
				var letter rune
				if key.Key>=65&&key.Key<=90 {
					caps,_,_:=getKeyState.Call(0x14)
					upper:=down(0x10)!=(caps&1!=0)
					letter=rune(key.Key);if !upper{letter=unicode.ToLower(letter)}
				}else if key.Key==0xDE&&!down(0x10){letter='\''}else if key.Key==0xC0&&!down(0x10){letter='`'}
				if key.Key==0x10||key.Key==0xA0||key.Key==0xA1{r,_,_:=nextHook.Call(a.keyboardHook,code,wp,lp);return r}
				previous:=a.previous
				var combined rune
				switch {
				case unicode.ToLower(letter)=='h'&&unicode.ToLower(previous)=='s':combined='ş'
				case unicode.ToLower(letter)=='h'&&unicode.ToLower(previous)=='c':combined='ç'
				case (letter=='\''||letter=='`')&&unicode.ToLower(previous)=='o':combined='ö'
				case (letter=='\''||letter=='`')&&unicode.ToLower(previous)=='g':combined='ğ'
				}
				a.previous=letter
				if combined!=0 {
					if unicode.IsUpper(previous){combined=unicode.ToUpper(combined)}
					a.previous=0
					if sendRune(combined,true){return 1}
				}
			}
		}
	}
	hook:=uintptr(0);if a!=nil{hook=a.keyboardHook}
	r,_,_:=nextHook.Call(hook,code,wp,lp);return r
}
func mouseProc(code,wp,lp uintptr) uintptr {
	if current!=nil&&int32(code)>=0&&(wp==0x201||wp==0x204||wp==0x207){current.previous=0}
	r,_,_:=nextHook.Call(0,code,wp,lp);return r
}
