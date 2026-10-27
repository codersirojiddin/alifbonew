package alphabet

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	cases := []struct{ input, want string; reverse bool }{
		{"O‘zbekiston, g‘alaba, shahar, chiroyli.", "Özbekiston, ğalaba, şahar, çiroyli.", false},
		{"O'ZBEK SHAHAR CHIROQ TOG'", "ÖZBEK ŞAHAR ÇIROQ TOĞ", false},
		{"o`rik gʻisht oʼq g’oz o´rmon", "örik ğişt öq ğoz örmon", false},
		{"ma'no san’at Isʼhoq", "maʼno sanʼat Isʼhoq", false},
		{"https://example.com/sh www.shahar.uz user@shahar.uz shahar", "https://example.com/sh www.shahar.uz user@shahar.uz şahar", false},
		{"'o' va 'g'", "'o' va 'g'", false},
		{"", "", false},
		{"ÖĞŞÇ ö ğ ş ç", "ÖĞŞÇ ö ğ ş ç", false},
		{"Shahar — SHAHAR — Chiroq — CHIROQ", "Şahar — ŞAHAR — Çiroq — ÇIROQ", false},
		{"Özbekiston, Şahar, ŞAHAR, ÇIROQ", "O‘zbekiston, Shahar, SHAHAR, CHIROQ", true},
		{"ÖĞŞÇ", "O‘G‘SHCH", true},
		{"sh\r\nch\t o‘rik\u00a0tog‘", "ş\r\nç\t örik\u00a0toğ", false},
		{"Salom 😀 shahar", "Salom 😀 şahar", false},
	}
	for _, c := range cases {
		if got := Convert(c.input, c.reverse); got != c.want { t.Errorf("Convert(%q, %v) = %q, want %q", c.input, c.reverse, got, c.want) }
	}
}

func TestLongText(t *testing.T) {
	input := strings.Repeat("O‘zbekiston shahar chiroq tog‘ ", 10000)
	want := strings.Repeat("Özbekiston şahar çiroq toğ ", 10000)
	if got := Convert(input, false); got != want { t.Fatal("long conversion differs") }
}
