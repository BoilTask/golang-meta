package metautf

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var encodings = map[string]encoding.Encoding{
	"UTF-8":        unicode.UTF8,
	"GBK":          simplifiedchinese.GBK,
	"GB18030":      simplifiedchinese.GB18030,
	"Big5":         traditionalchinese.Big5,
	"Shift_JIS":    japanese.ShiftJIS,
	"EUC-JP":       japanese.EUCJP,
	"EUC-KR":       korean.EUCKR,
	"windows-1252": charmap.Windows1252,
}

// 宽容修复 UTF-8（不丢字符）
func fixUTF8(s string) string {
	r := transform.NewReader(strings.NewReader(s), unicode.UTF8.NewDecoder())
	b, _ := io.ReadAll(r)
	return string(b)
}

// 自动检测编码并转 UTF-8
func autoConvertToUTF8(s string) string {
	d := chardet.NewTextDetector()
	r, err := d.DetectBest([]byte(s))
	if err != nil {
		return fixUTF8(s)
	}
	enc, ok := encodings[r.Charset]
	if !ok {
		return fixUTF8(s)
	}

	reader := transform.NewReader(strings.NewReader(s), enc.NewDecoder())
	out, err := io.ReadAll(reader)
	if err != nil {
		return fixUTF8(s)
	}
	return fixUTF8(string(out))
}

// 最后兜底：非法字节替换为 '�'，而不是丢弃
func replaceInvalidUTF8(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError && size == 1 {
			b.WriteRune('�')
			s = s[1:]
		} else {
			b.WriteRune(r)
			s = s[size:]
		}
	}
	return b.String()
}

// removeNullBytes 去掉或替换掉 PostgreSQL 禁止的 0x00
func removeNullBytes(s string) string {
	// 通常使用替换为空字符串：删除
	return strings.ReplaceAll(s, "\x00", "")
	// 如果你希望保留占位：
	// return strings.ReplaceAll(s, "\x00", "�")
}

func SanitizeText(text string) string {
	if utf8.ValidString(text) {
		return removeNullBytes(text)
	}

	converted := autoConvertToUTF8(text)
	if utf8.ValidString(converted) {
		return removeNullBytes(converted)
	}

	return removeNullBytes(replaceInvalidUTF8(converted))
}

func SanitizeTexts(texts []string) (validTexts []string, invalidTexts []string) {
	for _, t := range texts {
		s := SanitizeText(t)
		// removeNullBytes 已经确保不会再因 0x00 报 PostgreSQL 错
		validTexts = append(validTexts, s)
		// 注意：这里 invalidTexts 现在实际不会再出现
	}
	return
}
