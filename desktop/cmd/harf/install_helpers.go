package main

import (
	"net/url"
	"strings"
)

func fileURL(path string) string {
	return (&url.URL{Scheme:"file",Path:"/"+strings.ReplaceAll(path,`\`,"/")}).String()
}
func crlf(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s,"\r\n","\n"),"\n","\r\n") }
