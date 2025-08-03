package common

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

var codePageMapping = map[int]encoding.Encoding{
	037:   charmap.CodePage037,   // IBM EBCDIC US-Canada
	437:   charmap.CodePage437,   // OEM United States
	850:   charmap.CodePage850,   // Western European (DOS)
	852:   charmap.CodePage852,   // Central European (DOS)
	855:   charmap.CodePage855,   // OEM Cyrillic (primarily Russian)
	858:   charmap.CodePage858,   // OEM Multilingual Latin 1 + Euro
	860:   charmap.CodePage860,   // Portuguese (DOS)
	862:   charmap.CodePage862,   // Hebrew (DOS)
	863:   charmap.CodePage863,   // French Canadian (DOS)
	865:   charmap.CodePage865,   // Nordic (DOS)
	866:   charmap.CodePage866,   // Russian (DOS)
	936:   simplifiedchinese.GBK, // Chinese (GBK)
	1047:  charmap.CodePage1047,  // IBM EBCDIC Latin 1/Open System
	1140:  charmap.CodePage1140,  // IBM EBCDIC US-Canada with Euro
	1250:  charmap.Windows1250,   // Central European (Windows)
	1251:  charmap.Windows1251,   // Cyrillic (Windows)
	1252:  charmap.Windows1252,   // Western European (Windows)
	1253:  charmap.Windows1253,   // Greek (Windows)
	1254:  charmap.Windows1254,   // Turkish (Windows)
	1255:  charmap.Windows1255,   // Hebrew (Windows)
	1256:  charmap.Windows1256,   // Arabic (Windows)
	1257:  charmap.Windows1257,   // Baltic (Windows)
	1258:  charmap.Windows1258,   // Vietnamese (Windows)
	20866: charmap.KOI8R,         // Russian (KOI8-R)
	21866: charmap.KOI8U,         // Ukrainian (KOI8-U)
	28591: charmap.ISO8859_1,     // Western European (ISO 8859-1)
	28592: charmap.ISO8859_2,     // Central European (ISO 8859-2)
	28593: charmap.ISO8859_3,     // Latin 3 (ISO 8859-3)
	28594: charmap.ISO8859_4,     // Baltic (ISO 8859-4)
	28595: charmap.ISO8859_5,     // Cyrillic (ISO 8859-5)
	28596: charmap.ISO8859_6,     // Arabic (ISO 8859-6)
	28597: charmap.ISO8859_7,     // Greek (ISO 8859-7)
	28598: charmap.ISO8859_8,     // Hebrew (ISO 8859-8)
	28599: charmap.ISO8859_9,     // Turkish (ISO 8859-9)
	28605: charmap.ISO8859_15,    // Latin 9 (ISO 8859-15)
	65001: encoding.Nop,          // Unicode (UTF-8)
}

func ConvertCpToUTF8(input string, codePage int) string {
	enc, exists := codePageMapping[codePage]
	if !exists {
		return input
	}

	reader := transform.NewReader(strings.NewReader(input), enc.NewDecoder())
	utf8Text, err := io.ReadAll(reader)
	if err != nil {
		return input
	}

	return string(utf8Text)
}

func ConvertUTF8toCp(input string, codePage int) string {
	enc, exists := codePageMapping[codePage]
	if !exists {
		return input
	}

	transform.NewWriter(io.Discard, enc.NewEncoder())
	encodedText, err := io.ReadAll(transform.NewReader(strings.NewReader(input), enc.NewEncoder()))
	if err != nil {
		return input
	}

	return string(encodedText)
}
