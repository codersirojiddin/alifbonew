package main

import "testing"

func TestShortcutURL(t *testing.T) {
	got:=fileURL(`C:\Users\Sir Name\AppData\Local\Harf\Harf.exe`)
	want:="file:///C:/Users/Sir%20Name/AppData/Local/Harf/Harf.exe"
	if got!=want{t.Fatalf("fileURL=%q, want %q",got,want)}
}
func TestCRLF(t *testing.T){if got:=crlf("a\nb\r\nc");got!="a\r\nb\r\nc"{t.Fatal(got)}}
