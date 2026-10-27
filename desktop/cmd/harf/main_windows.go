//go:build windows && amd64

// Harf uses native Windows controls and the Go standard library. No .NET,
// WebView, browser process, CGO or third-party runtime is required.
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode"
	"unsafe"

	"github.com/codersirojiddin/alifbonew/desktop/internal/alphabet"
)

const (
	trayMessage = 0x8002
	autoID = 101
	startupID = 102
	forwardID = 103
	reverseID = 104
	copyID = 105
	exitID = 106
	showID = 107
)

type insertion struct { Letter rune; Window, Focus uintptr; Started time.Time }
type application struct {
	window, instance, source, output, state, autoBox, startupBox uintptr
	keyboardHook, mouseHook uintptr
	icon notifyIcon
	auto bool
	previous rune
	previousWindow, previousFocus uintptr
	previousTime time.Time
	processWindow uintptr
	processBlocked bool
	queue []insertion
	smoke bool
}

var current *application

func main() {
	runtime.LockOSThread()
	executable, err := os.Executable()
	if err != nil { dialog(err.Error(), 0x10); os.Exit(1) }
	if len(os.Args)>1 && os.Args[1]=="--version" { fmt.Println("Harf 0.2.0 Go"); return }
	smoke := len(os.Args)>1 && os.Args[1]=="--smoke-test"
	if !smoke && strings.HasSuffix(strings.ToLower(executable), "-setup.exe") {
		if err:=install(executable);err!=nil { dialog(err.Error(),0x10);os.Exit(1) };return
	}
	mutex, _, lastError := createMutex.Call(0, 0, uintptr(unsafe.Pointer(wide("Local\\Harf.InputHelper.v1"))))
	if mutex==0 { dialog("Harf mutex could not be created.",0x10);os.Exit(1) }
	defer closeHandle.Call(mutex)
	if lastError==syscall.Errno(183) { dialog("Harf allaqachon ishlayapti. Tray belgisidan oching.",0x40);return }
	a:=&application{smoke:smoke}
	current=a
	if err:=a.run();err!=nil { if smoke { fmt.Fprintln(os.Stderr,err) } else { dialog(err.Error(),0x10) };os.Exit(1) }
}

func (a *application) run() error {
	instance,_,_:=getModuleHandle.Call(0)
	a.instance=instance
	icon,_,_:=loadIcon.Call(0,32512)
	cursor,_,_:=loadCursor.Call(0,32512)
	brush,_,_:=getStockObject.Call(0) // WHITE_BRUSH
	class:=windowClass{Size:uint32(unsafe.Sizeof(windowClass{})),Proc:syscall.NewCallback(windowProc),Instance:instance,Icon:icon,SmallIcon:icon,Cursor:cursor,Background:brush,ClassName:wide("HarfGoWindow")}
	registered,_,err:=registerClass.Call(uintptr(unsafe.Pointer(&class)))
	if registered==0{return fmt.Errorf("RegisterClassEx: %v",err)}
	window,_,err:=createWindow.Call(0,uintptr(unsafe.Pointer(class.ClassName)),uintptr(unsafe.Pointer(wide("Harf — Go | Yangi alifbo"))),0x00CA0000,120,100,730,670,0,0,instance,0)
	if window==0{return fmt.Errorf("CreateWindowEx: %v",err)}
	a.window=window
	defer a.cleanup()
	a.controls()
	a.icon=notifyIcon{Size:uint32(unsafe.Sizeof(notifyIcon{})),Window:window,ID:1,Flags:7,Callback:trayMessage,Icon:icon}
	copy(a.icon.Tip[:],syscall.StringToUTF16("Harf — Go yozish yordamchisi"))
	trayOK,_,trayErr:=shellNotify.Call(0,uintptr(unsafe.Pointer(&a.icon)))
	if trayOK==0&&!a.smoke{return fmt.Errorf("Tray icon: %v",trayErr)}
	for i,key:=range []uintptr{'O','G','S','C'} {
		for _,offset:=range []int{10,20} {
			modifiers:=uintptr(0x4003);if offset==20{modifiers=0x4007}
			ok,_,err:=registerHotKey.Call(window,uintptr(offset+i),modifiers,key)
			if ok==0{return fmt.Errorf("AltGr shortcut band (%c): %v. Boshqa klaviatura yordamchisini yoping.",key,err)}
		}
	}
	if ok,_,err:=registerHotKey.Call(window,30,0x4003,'Y');ok==0{return fmt.Errorf("Ctrl+Alt+Y band: %v",err)}
	a.keyboardHook,_,err=setHook.Call(13,syscall.NewCallback(keyboardProc),instance,0)
	if a.keyboardHook==0{return fmt.Errorf("Keyboard hook: %v",err)}
	a.mouseHook,_,err=setHook.Call(14,syscall.NewCallback(mouseProc),instance,0)
	if a.mouseHook==0{return fmt.Errorf("Mouse hook: %v",err)}
	a.updateStatus()
	if a.smoke {setTimer.Call(window,99,300,0)}else{showWindow.Call(window,5);updateWindow.Call(window)}
	var msg message
	for {
		r,_,err:=getMessage.Call(uintptr(unsafe.Pointer(&msg)),0,0,0)
		if int32(r)<0{return fmt.Errorf("GetMessage: %v",err)}
		if r==0{break}
		translateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		dispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
	return nil
}

func (a *application) control(class,title string,style uintptr,x,y,w,h,id int) uintptr {
	window,_,_:=createWindow.Call(0,uintptr(unsafe.Pointer(wide(class))),uintptr(unsafe.Pointer(wide(title))),0x50000000|style,uintptr(x),uintptr(y),uintptr(w),uintptr(h),a.window,uintptr(id),a.instance,0)
	font,_,_:=getStockObject.Call(17)
	sendMessage.Call(window,0x30,font,1)
	return window
}
func (a *application) controls() {
	a.control("STATIC","Harf.   Yangi alifbo, tanish klaviatura.",0,24,20,660,30,0)
	a.state=a.control("STATIC","",0,24,58,660,34,0)
	a.autoBox=a.control("BUTTON","Avtoalmashtirish (English US klaviaturasi)",3,24,98,660,26,autoID)
	a.startupBox=a.control("BUTTON","Windows bilan ishga tushirish",3,24,132,660,26,startupID)
	if startupEnabled(){sendMessage.Call(a.startupBox,0xF1,1,0)}
	a.control("STATIC","AltGr + O/G/S/C → ö/ğ/ş/ç    |    Shift → katta harf\r\nCtrl + Alt + Y → avtoalmashtirishni almashtirish",0,24,175,660,42,0)
	a.control("STATIC","Matn konvertori",0,24,228,660,24,0)
	a.source=a.control("EDIT","",0x00A110C4,24,257,660,108,201)
	sendMessage.Call(a.source,0xC5,16*1024*1024,0)
	a.control("BUTTON","Eski → Yangi",0x10000,24,378,145,32,forwardID)
	a.control("BUTTON","Yangi → Eski",0x10000,182,378,145,32,reverseID)
	a.output=a.control("EDIT","",0x00A118C4,24,424,660,108,202)
	a.control("BUTTON","Nusxalash",0x10000,24,548,145,32,copyID)
	a.control("BUTTON","Dasturdan chiqish",0x10000,519,548,165,32,exitID)
}
func (a *application) updateStatus() {
	s:= "Avtoalmashtirish O‘CHIQ. AltGr tugmalari ishlaydi."
	if a.auto{s="Avtoalmashtirish YOQILGAN. Maxfiy maydonlarda o‘chirib turing."}
	setText(a.state,s)
	check:=uintptr(0);if a.auto{check=1}
	sendMessage.Call(a.autoBox,0xF1,check,0)
	a.previous=0
}
func (a *application) toggleAuto(){a.auto=!a.auto;a.updateStatus()}
func (a *application) command(id int) {
	switch id {
	case autoID:a.toggleAuto()
	case startupID:
		checked,_,_:=sendMessage.Call(a.startupBox,0xF0,0,0)
		if err:=setStartup(checked==1);err!=nil{dialog(err.Error(),0x10);if startupEnabled(){checked=1}else{checked=0};sendMessage.Call(a.startupBox,0xF1,checked,0)}
	case forwardID:setText(a.output,alphabet.Convert(text(a.source),false))
	case reverseID:setText(a.output,alphabet.Convert(text(a.source),true))
	case copyID:if !clipboardCopy(a.window,text(a.output)){dialog("Nusxalash amalga oshmadi. Matnni belgilang va Ctrl+C bosing.",0x40)}
	case showID:showWindow.Call(a.window,9);setForeground.Call(a.window)
	case exitID:destroyWindow.Call(a.window)
	}
}
func (a *application) trayMenu() {
	menu,_,_:=createPopupMenu.Call()
	defer destroyMenu.Call(menu)
	appendMenu.Call(menu,0,showID,uintptr(unsafe.Pointer(wide("Oynani ochish"))))
	label:="Avtoalmashtirishni yoqish";if a.auto{label="Avtoalmashtirishni o‘chirish"}
	appendMenu.Call(menu,0,autoID,uintptr(unsafe.Pointer(wide(label))))
	appendMenu.Call(menu,0,exitID,uintptr(unsafe.Pointer(wide("Chiqish"))))
	var p point;getCursorPos.Call(uintptr(unsafe.Pointer(&p)))
	setForeground.Call(a.window)
	id,_,_:=trackPopupMenu.Call(menu,0x102,uintptr(p.X),uintptr(p.Y),0,a.window,0)
	if id!=0{a.command(int(id))}
}
func windowProc(window,msg,wp,lp uintptr) uintptr {
	a:=current
	if a!=nil {
		switch msg {
		case 0x111:a.command(int(wp&0xFFFF));return 0 // WM_COMMAND
		case 0x10:showWindow.Call(window,0);return 0 // Close hides to tray.
		case 2:postQuitMessage.Call(0);return 0
		case trayMessage:
			if lp==0x203{a.command(showID)}else if lp==0x205{a.trayMenu()};return 0
		case 0x312:
			id:=int(wp)
			if id==30{a.toggleAuto()}else if id>=10&&id<=13||id>=20&&id<=23 {
				letter:=[]rune{'ö','ğ','ş','ç'}[id%10];if id>=20{letter=unicode.ToUpper(letter)}
				window,_,_:=getForeground.Call();focus:=focusedControl(window)
				if len(a.queue)<16{a.queue=append(a.queue,insertion{letter,window,focus,time.Now()});setTimer.Call(a.window,1,15,0)}
			}
			a.previous=0;return 0
		case 0x113:
			if wp==99{killTimer.Call(window,99);destroyWindow.Call(window)}else if wp==1{a.flushQueue()};return 0
		}
	}
	r,_,_:=defWindowProc.Call(window,msg,wp,lp);return r
}
func (a *application) flushQueue(){
	window,_,_:=getForeground.Call();focus:=focusedControl(window)
	remaining:=a.queue[:0]
	for _,item:=range a.queue {
		if item.Window!=window||item.Focus!=focus||time.Since(item.Started)>1500*time.Millisecond{continue}
		if down(0x11)||down(0x12)||down(0x10){remaining=append(remaining,item);continue}
		if !sendRune(item.Letter,false){setText(a.state,"Bu dasturga harf kiritib bo‘lmadi. Matn konvertoridan foydalaning.")}
	}
	a.queue=remaining
	if len(a.queue)==0{killTimer.Call(a.window,1)}
}
func (a *application) cleanup(){
	if a.keyboardHook!=0{unhook.Call(a.keyboardHook)}
	if a.mouseHook!=0{unhook.Call(a.mouseHook)}
	for id:=10;id<=30;id++{unregisterHotKey.Call(a.window,uintptr(id))}
	shellNotify.Call(2,uintptr(unsafe.Pointer(&a.icon)))
}
