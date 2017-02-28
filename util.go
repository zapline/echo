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

func TimeFormat(tbuf []byte, t time.Time) {
	tbuf = tbuf[:19]
	year, month, day := t.Date()
	tbuf[0], tbuf[1] = getdigit(byte(year / 100))
	tbuf[2], tbuf[3] = getdigit(byte(year % 100))
	tbuf[4] = '-'
	tbuf[5], tbuf[6] = getdigit(byte(month))
	tbuf[7] = '-'
	tbuf[8], tbuf[9] = getdigit(byte(day))
	tbuf[10] = '/'
	hour, min, sec := t.Clock()
	tbuf[11], tbuf[12] = getdigit(byte(hour))
	tbuf[13] = ':'
	tbuf[14], tbuf[15] = getdigit(byte(min))
	tbuf[16] = ':'
	tbuf[17], tbuf[18] = getdigit(byte(sec))
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
