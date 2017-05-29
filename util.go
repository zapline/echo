package echo

import (
	"sync"
	"time"
	"unicode/utf8"

	"github.com/vizee/litebuf"
)

const (
	hexdigits = "0123456789abcdef"

	noescchr = '0'
	esctable = `00000000btn0fr00000000000000000000"000000000000/00000000000000000000000000000000000000000000\00000000000000000000000000000000000`
)

var bufpool = sync.Pool{
	New: func() interface{} {
		return &litebuf.Buffer{}
	},
}

func getdigit(n byte) (byte, byte) {
	return n/10 + '0', n%10 + '0'
}

func TimeFormat(buf []byte, t time.Time) {
	buf = buf[:19]
	year, month, day := t.Date()
	buf[0], buf[1] = getdigit(byte(year / 100))
	buf[2], buf[3] = getdigit(byte(year % 100))
	buf[4] = '-'
	buf[5], buf[6] = getdigit(byte(month))
	buf[7] = '-'
	buf[8], buf[9] = getdigit(byte(day))
	buf[10] = '/'
	hour, min, sec := t.Clock()
	buf[11], buf[12] = getdigit(byte(hour))
	buf[13] = ':'
	buf[14], buf[15] = getdigit(byte(min))
	buf[16] = ':'
	buf[17], buf[18] = getdigit(byte(sec))
}

func QuoteString(buf *litebuf.Buffer, s string, unicode bool) {
	buf.WriteByte('"')
	p := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < utf8.RuneSelf {
			if esctable[c] != noescchr {
				if p < i {
					buf.WriteString(s[p:i])
				}
				buf.WriteByte('\\')
				buf.WriteByte(esctable[c])
				p = i + 1
			}
		} else if unicode {
			if p < i {
				buf.WriteString(s[p:i])
			}
			r, n := utf8.DecodeRuneInString(s[i:])
			h := [6]byte{
				'\\',
				'u',
				hexdigits[(r>>12)&0xf],
				hexdigits[(r>>8)&0xf],
				hexdigits[(r>>4)&0xf],
				hexdigits[(r)&0xf],
			}
			buf.Write(h[:])
			i += n - 1
			p = i + 1
		}
	}
	if p < len(s) {
		buf.WriteString(s[p:])
	}
	buf.WriteByte('"')
}
