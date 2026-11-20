package writer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
)

// The scalar rules of the `yaml` npm library (stringifyString.js) with the TypeScript SDK writer's options
// (lineWidth 0, so no folding), as the PHP SDK ports them: every SDK quotes the same text the same way.

var plainForbidden = regexp.MustCompile("^[\n\t ,\\[\\]{}#&*!|>'\"%@`]|^[?-]$|^[?-][ \t]|[\n:][ \t]|[ \t]\n|[\n\t ]#|[\n\t :]$")

var otherCoreTypes = []*regexp.Regexp{
	regexp.MustCompile(`^(?:~|[Nn]ull|NULL)?$`),
	regexp.MustCompile(`^(?:[Tt]rue|TRUE|[Ff]alse|FALSE)$`),
	regexp.MustCompile(`^0o[0-7]+$`),
	regexp.MustCompile(`^[-+]?[0-9]+$`),
	regexp.MustCompile(`^0x[0-9a-fA-F]+$`),
	regexp.MustCompile(`^(?:[-+]?\.(?:inf|Inf|INF)|\.nan|\.NaN|\.NAN)$`),
	regexp.MustCompile(`^[-+]?(?:\.[0-9]+|[0-9]+(?:\.[0-9]*)?)[eE][-+]?[0-9]+$`),
	regexp.MustCompile(`^[-+]?(?:\.[0-9]+|[0-9]+\.[0-9]*)$`),
}

var controlCharacter = regexp.MustCompile(`[\x00-\x08\x0b-\x1f\x7f-\x{9f}]`)

var blockEndsWithIndentedLine = regexp.MustCompile("\n[\t ]+$")

var spaceAroundLineBreak = regexp.MustCompile("[ \t]\n|\n[ \t]")

var lineBreaks = regexp.MustCompile("\n+")

const doubleQuotedMinMultiLineLength = 40

// StringToYaml takes the indentation of the scalar's own children (the pair's indentation plus one step).
// Invalid UTF-8 becomes U+FFFD whatever the style, so the file is valid UTF-8.
func StringToYaml(value, indent string, implicitKey bool) string {
	value = string([]rune(value))
	if controlCharacter.MatchString(value) {
		return stringToDoubleQuoted(value, indent, implicitKey)
	}
	return stringToPlain(value, indent, implicitKey)
}

func StringToUTF16Length(value string) int {
	return len(utf16.Encode([]rune(value)))
}

func stringToDoubleQuoted(value, indent string, implicitKey bool) string {
	json := stringToJSON(value)
	multiLine := StringToUTF16Length(json) >= doubleQuotedMinMultiLineLength
	var text strings.Builder
	start := 0
	for i := 0; i < len(json); i++ {
		ch := json[i]
		if ch == ' ' && charAt(json, i+1) == '\\' && charAt(json, i+2) == 'n' {
			text.WriteString(json[start:i] + "\\ ")
			i++
			start = i
			ch = '\\'
		}
		if ch != '\\' {
			continue
		}
		next := charAt(json, i+1)
		if next == 'u' {
			text.WriteString(json[start:i] + unicodeEscapeToYaml(json[i+2:i+6]))
			i += 5
			start = i + 1
			continue
		}
		if next != 'n' || implicitKey || charAt(json, i+2) == '"' || !multiLine {
			i++
			continue
		}
		text.WriteString(json[start:i] + "\n\n")
		for charAt(json, i+2) == '\\' && charAt(json, i+3) == 'n' && charAt(json, i+4) != '"' {
			text.WriteString("\n")
			i += 2
		}
		text.WriteString(indent)
		if charAt(json, i+2) == ' ' {
			text.WriteString("\\")
		}
		i++
		start = i + 1
	}
	if start == 0 {
		return json
	}
	return text.String() + json[start:]
}

func stringToPlain(value, indent string, implicitKey bool) string {
	multiLine := strings.Contains(value, "\n")
	if implicitKey && multiLine {
		return stringToQuoted(value, indent, implicitKey)
	}
	if plainForbidden.MatchString(value) {
		if implicitKey || !multiLine {
			return stringToQuoted(value, indent, implicitKey)
		}
		return stringToBlock(value, indent, implicitKey)
	}
	if multiLine {
		return stringToBlock(value, indent, implicitKey)
	}
	for _, pattern := range otherCoreTypes {
		if pattern.MatchString(value) {
			return stringToQuoted(value, indent, implicitKey)
		}
	}
	return value
}

func stringToQuoted(value, indent string, implicitKey bool) string {
	if strings.Contains(value, `"`) && !strings.Contains(value, "'") {
		return stringToSingleQuoted(value, indent, implicitKey)
	}
	return stringToDoubleQuoted(value, indent, implicitKey)
}

func stringToSingleQuoted(value, indent string, implicitKey bool) string {
	if (implicitKey && strings.Contains(value, "\n")) || spaceAroundLineBreak.MatchString(value) {
		return stringToDoubleQuoted(value, indent, implicitKey)
	}
	return "'" + lineBreaksToIndented(strings.ReplaceAll(value, "'", "''"), "\n"+indent) + "'"
}

func stringToBlock(value, indent string, implicitKey bool) string {
	if blockEndsWithIndentedLine.MatchString(value) {
		return stringToQuoted(value, indent, implicitKey)
	}
	endStart := len(value)
	for endStart > 0 && strings.ContainsRune("\n\t ", rune(value[endStart-1])) {
		endStart--
	}
	end := value[endStart:]
	endLineBreak := strings.Index(end, "\n")
	chomp := ""
	if endLineBreak == -1 {
		chomp = "-"
	} else if value == end || endLineBreak != len(end)-1 {
		chomp = "+"
	}
	if end != "" {
		value = value[:endStart]
		end = strings.TrimSuffix(end, "\n")
		end = innerLineBreaksToIndented(end, indent)
	}
	startsWithSpace := false
	startLineBreak := -1
	for index := 0; index < len(value) && strings.ContainsRune(" \n", rune(value[index])); index++ {
		if value[index] == ' ' {
			startsWithSpace = true
			continue
		}
		startLineBreak = index
	}
	start := value[:startLineBreak+1]
	if start != "" {
		value = value[len(start):]
		start = lineBreaksToIndented(start, indent)
	}
	header := chomp
	if startsWithSpace {
		header = "2" + chomp
	}
	return "|" + header + "\n" + indent + start + lineBreaksToIndented(value, indent) + end
}

func lineBreaksToIndented(value, indent string) string {
	return lineBreaks.ReplaceAllStringFunc(value, func(breaks string) string {
		return breaks + indent
	})
}

func innerLineBreaksToIndented(value, indent string) string {
	var text strings.Builder
	for index := 0; index < len(value); index++ {
		text.WriteByte(value[index])
		if value[index] == '\n' && index+1 < len(value) && value[index+1] != '\n' {
			text.WriteString(indent)
		}
	}
	return text.String()
}

func stringToJSON(value string) string {
	var text strings.Builder
	text.WriteByte('"')
	for _, r := range value {
		switch r {
		case '"':
			text.WriteString(`\"`)
		case '\\':
			text.WriteString(`\\`)
		case '\b':
			text.WriteString(`\b`)
		case '\f':
			text.WriteString(`\f`)
		case '\n':
			text.WriteString(`\n`)
		case '\r':
			text.WriteString(`\r`)
		case '\t':
			text.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(&text, `\u%04x`, r)
				continue
			}
			text.WriteRune(r)
		}
	}
	text.WriteByte('"')
	return text.String()
}

func unicodeEscapeToYaml(code string) string {
	escapes := map[string]string{
		"0000": `\0`,
		"0007": `\a`,
		"000b": `\v`,
		"001b": `\e`,
	}
	if escape, found := escapes[code]; found {
		return escape
	}
	return `\x` + code[2:]
}

func charAt(value string, index int) byte {
	if index >= len(value) {
		return 0
	}
	return value[index]
}
