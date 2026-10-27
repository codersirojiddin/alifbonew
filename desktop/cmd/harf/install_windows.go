//go:build windows && amd64

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

const runKey = `Software\Microsoft\Windows\CurrentVersion\Run`

func startupEnabled() bool {
	var key uintptr
	r,_,_:=regOpen.Call(^uintptr(0x7ffffffe),uintptr(unsafe.Pointer(wide(runKey))),0,1,uintptr(unsafe.Pointer(&key)))
	if r!=0{return false}
	defer regClose.Call(key)
	var size uint32
	r,_,_=regQuery.Call(key,uintptr(unsafe.Pointer(wide("Harf"))),0,0,0,uintptr(unsafe.Pointer(&size)))
	return r==0
}
func setStartup(enabled bool) error {
	var key uintptr
	r,_,_:=regCreate.Call(^uintptr(0x7ffffffe),uintptr(unsafe.Pointer(wide(runKey))),0,0,0,2,0,uintptr(unsafe.Pointer(&key)),0)
	if r!=0{return fmt.Errorf("Startup registry access: %v",syscall.Errno(r))}
	defer regClose.Call(key)
	if !enabled {
		r,_,_=regDelete.Call(key,uintptr(unsafe.Pointer(wide("Harf"))))
		if r==0||r==2{return nil};return fmt.Errorf("Remove startup: %v",syscall.Errno(r))
	}
	executable,err:=os.Executable();if err!=nil{return err}
	value,err:=syscall.UTF16FromString(`"`+executable+`"`);if err!=nil{return err}
	r,_,_=regSet.Call(key,uintptr(unsafe.Pointer(wide("Harf"))),0,1,uintptr(unsafe.Pointer(&value[0])),uintptr(len(value)*2))
	runtime.KeepAlive(value)
	if r!=0{return fmt.Errorf("Set startup: %v",syscall.Errno(r))};return nil
}

func install(executable string) error {
	if dialog("Harf Go yordamchisini o‘rnatasizmi?\n\nDastur joriy Windows foydalanuvchisi uchun o‘rnatiladi.\nAvtomatik almashtirish avval o‘chiq turadi.",0x41)!=1{return nil}
	local,appdata:=os.Getenv("LOCALAPPDATA"),os.Getenv("APPDATA")
	if local==""||appdata==""{return fmt.Errorf("Windows user folders are unavailable")}
	folder:=filepath.Join(local,"Harf")
	if err:=os.MkdirAll(folder,0755);err!=nil{return err}
	data,err:=os.ReadFile(executable);if err!=nil{return err}
	target:=filepath.Join(folder,"Harf.exe")
	if err:=os.WriteFile(target,data,0755);err!=nil{return fmt.Errorf("Eski Harf dasturini tray menyusidan yoping va qayta o‘rnating: %w",err)}
	programs:=filepath.Join(appdata,"Microsoft","Windows","Start Menu","Programs")
	if err:=os.MkdirAll(programs,0755);err!=nil{return err}
	shortcut:="[InternetShortcut]\r\nURL="+fileURL(target)+"\r\nIconFile="+target+"\r\nIconIndex=0\r\n"
	if err:=os.WriteFile(filepath.Join(programs,"Harf.url"),[]byte(shortcut),0644);err!=nil{return err}
	uninstall:=`@echo off
echo Close Harf from its tray menu first.
pause
del "%~dp0Harf.exe"
if exist "%~dp0Harf.exe" (echo Harf is still running. Close it and retry. & pause & exit /b 1)
reg delete HKCU\Software\Microsoft\Windows\CurrentVersion\Run /v Harf /f >nul 2>&1
del "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Harf.url" >nul 2>&1
cd /d "%TEMP%"
rmdir /s /q "%~dp0"
`
	if err:=os.WriteFile(filepath.Join(folder,"Uninstall.cmd"),[]byte(crlf(uninstall)),0644);err!=nil{return err}
	return exec.Command(target).Start()
}
