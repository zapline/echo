package echo

import (
	"testing"
	"time"

	"github.com/vizee/litebuf"
)

func TestPlainFormat(t *testing.T) {
	buf := litebuf.Buffer{}
	f := PlainFormatter{}
	f.Format(&buf, time.Now(), DebugLevel, "Debug Hello", []Field{String("Name", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.ShowSource = 2
	f.Format(&buf, time.Now(), InfoLevel, "Info Hello", []Field{String("Name", "世\t界")})
	t.Log(buf.String())
}

func TestJSONFormat(t *testing.T) {
	buf := litebuf.Buffer{}
	f := JSONFormatter{}
	f.Format(&buf, time.Now(), DebugLevel, "Debug Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.Format(&buf, time.Now(), InfoLevel, "Info Hello", []Field{String("who", "世\t界")})
	t.Log(buf.String())
	buf.Reset()
	f.EscapeUnicode = true
	f.Format(&buf, time.Now(), InfoLevel, "Info Hello", []Field{String("who", "世界")})
	t.Log(buf.String())
	buf.Reset()
	f.LevelTag = true
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.Type = RFC3339Time
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.Type = UnixTimeStamp
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
	buf.Reset()
	f.Type = UnixTimeStampNano
	f.Format(&buf, time.Now(), InfoLevel, "Hello", []Field{String("who", "World")})
	t.Log(buf.String())
}
