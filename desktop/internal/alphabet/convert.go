// Package alphabet converts Uzbek Latin text without network or file access.
package alphabet

import (
	"strings"
	"unicode"
)

const apostrophes = "'`´ʻʼ‘’"

// Convert preserves whitespace, URLs and email tokens. Reverse expands the new
// letters into the old Latin combinations; original apostrophe variants are lost.
func Convert(text string, reverse bool) string {
	var out strings.Builder
	start := 0
	for pos, r := range text {
		if unicode.IsSpace(r) {
			out.WriteString(token(text[start:pos], reverse))
			out.WriteRune(r)
			start = pos + len(string(r))
		}
	}
	out.WriteString(token(text[start:], reverse))
	return out.String()
}

func token(text string, reverse bool) string {
	if strings.Contains(text, "://") || strings.ContainsRune(text, '@') || strings.HasPrefix(strings.ToLower(text), "www.") {
		return text
	}
	runes := []rune(text)
	var out strings.Builder
	if reverse {
		upper, lower := false, false
		for _, r := range runes {
			upper = upper || unicode.IsUpper(r)
			lower = lower || unicode.IsLower(r)
		}
		for _, r := range runes {
			switch r {
			case 'ö': out.WriteString("o‘")
			case 'Ö': out.WriteString("O‘")
			case 'ğ': out.WriteString("g‘")
			case 'Ğ': out.WriteString("G‘")
			case 'ş': out.WriteString("sh")
			case 'Ş': if upper && !lower { out.WriteString("SH") } else { out.WriteString("Sh") }
			case 'ç': out.WriteString("ch")
			case 'Ç': if upper && !lower { out.WriteString("CH") } else { out.WriteString("Ch") }
			default: out.WriteRune(r)
			}
		}
		return out.String()
	}
	letter := func(i int) bool { return i >= 0 && i < len(runes) && unicode.IsLetter(runes[i]) }
	for i := 0; i < len(runes); i++ {
		current := runes[i]
		var next rune
		if i+1 < len(runes) { next = runes[i+1] }
		low := unicode.ToLower(current)
		var combined rune
		if (low == 'o' || low == 'g') && next != 0 && strings.ContainsRune(apostrophes, next) && (letter(i+2) || low == 'g' && letter(i-1)) {
			if low == 'o' { combined = 'ö' } else { combined = 'ğ' }
		} else if unicode.ToLower(next) == 'h' {
			if low == 's' { combined = 'ş' } else if low == 'c' { combined = 'ç' }
		}
		if combined != 0 {
			if unicode.IsUpper(current) { combined = unicode.ToUpper(combined) }
			out.WriteRune(combined)
			i++
		} else if current != 'ʼ' && strings.ContainsRune(apostrophes, current) && letter(i-1) && letter(i+1) {
			out.WriteRune('ʼ')
		} else { out.WriteRune(current) }
	}
	return out.String()
}
