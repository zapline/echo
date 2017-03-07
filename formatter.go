package echo

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/vizee/litebuf"
)

var plainTags = [...]string{
	FatalLevel: " [F]",
	ErrorLevel: " [E]",
	WarnLevel:  " [W]",
	InfoLevel:  " [I]",
	DebugLevel: " [D]",
}

var defaultFormatter = &PlainFormatter{}

type PlainFormatter struct{}

func (*PlainFormatter) Format(buf *litebuf.Buffer, t time.Time, level LogLevel, msg string, fields []Field) {
	TimeFormat(buf.Reserve(19), t)

	buf.WriteString(plainTags[level])

	if msg != "" {
		buf.WriteByte(' ')
		buf.WriteString(msg)
	}

	if len(fields) > 0 {
		buf.WriteByte(' ')

		for i := range fields {
			if i > 0 {
				buf.WriteByte(' ')
			}

			field := &fields[i]

			if field.Key != "" {
				buf.WriteString(field.Key)
				buf.WriteByte('=')
			}

			switch field.Type {
			case TypeNil:
				buf.WriteString("nil")

			case TypeInt, TypeInt32, TypeInt64:
				buf.AppendInt(int64(field.U64), 10)

			case TypeUint, TypeUint32, TypeUint64:
				buf.AppendUint(field.U64, 10)

			case TypeHex:
				buf.WriteString("0x")
				buf.AppendUint(field.U64, 16)

			case TypeBool:
				if field.U64 == 1 {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}

			case TypeFloat32:
				buf.AppendFloat(float64(math.Float32frombits(uint32(field.U64))), 'f', -1, 32)

			case TypeFloat64:
				buf.AppendFloat(math.Float64frombits(field.U64), 'f', -1, 64)

			case TypeString:
				buf.WriteString(field.Str)

			case TypeQuote:
				QuoteString(buf, field.Str, false)

			case TypeError:
				if err := field.Data.(error); err != nil {
					buf.WriteString(err.Error())
				} else {
					buf.WriteString("nil")
				}

			case TypeStringer:
				buf.WriteString(field.Data.(fmt.Stringer).String())

			case TypeOutputer:
				field.Data.(Outputer).Output(buf)

			case TypeInterface:
				fmt.Fprint(buf, field.Data)

			case TypeStack:
				buf.WriteByte('\n')
				buf.Write(field.Data.([]byte))

			default:
				panic(fmt.Sprintf("unknown field type: %d", field.Type))

			}
		}
	}

	buf.WriteByte('\n')
}

type TimeType int

const (
	SimpleTime TimeType = iota
	RFC3339Time
	UnixTimeStamp
	UnixTimeStampNano
)

var jsonTags = [...]string{
	FatalLevel: `"fatal"`,
	ErrorLevel: `"error"`,
	WarnLevel:  `"warn"`,
	InfoLevel:  `"info"`,
	DebugLevel: `"debug"`,
}

type JSONFormatter struct {
	Type          TimeType
	LevelTag      bool // use level tag instead of level number
	EscapeUnicode bool // escape unicode
}

func (f *JSONFormatter) Format(buf *litebuf.Buffer, t time.Time, level LogLevel, msg string, fields []Field) {
	buf.WriteString(`{"time":`)

	switch f.Type {
	case SimpleTime:
		buf.WriteByte('"')
		TimeFormat(buf.Reserve(19)[:], t)
		buf.WriteByte('"')

	case RFC3339Time:
		buf.WriteByte('"')
		t.AppendFormat(buf.Reserve(25)[:0], time.RFC3339)
		buf.WriteByte('"')

	case UnixTimeStamp:
		buf.AppendInt(t.Unix(), 10)

	case UnixTimeStampNano:
		buf.AppendInt(t.UnixNano(), 10)

	default:
		panic(fmt.Sprintf("unknown time type: %d", f.Type))

	}

	buf.WriteString(`,"level":`)
	if f.LevelTag {
		buf.WriteString(jsonTags[level])
	} else {
		buf.AppendInt(int64(level), 10)
	}

	if msg != "" {
		buf.WriteString(`,"msg":`)
		QuoteString(buf, msg, f.EscapeUnicode)
	}

	if len(fields) > 0 {
		buf.WriteString(`,"fields":{`)

		for i := range fields {
			if i > 0 {
				buf.WriteByte(',')
			}

			field := &fields[i]

			// assert(field.Key not empty and field.Key not conflict)
			QuoteString(buf, field.Key, f.EscapeUnicode)

			buf.WriteByte(':')

			switch field.Type {
			case TypeNil:
				buf.WriteString("null")

			case TypeInt, TypeInt32, TypeInt64:
				buf.AppendInt(int64(field.U64), 10)

			case TypeUint, TypeUint32, TypeUint64:
				buf.AppendUint(field.U64, 10)

			case TypeHex:
				buf.WriteString(`"0x`)
				buf.AppendUint(field.U64, 16)
				buf.WriteByte('"')

			case TypeBool:
				if field.U64 == 1 {
					buf.WriteString("true")
				} else {
					buf.WriteString("false")
				}

			case TypeFloat32:
				buf.AppendFloat(float64(math.Float32frombits(uint32(field.U64))), 'f', -1, 32)

			case TypeFloat64:
				buf.AppendFloat(math.Float64frombits(field.U64), 'f', -1, 64)

			case TypeString, TypeQuote:
				QuoteString(buf, field.Str, f.EscapeUnicode)

			case TypeError:
				if err := field.Data.(error); err != nil {
					QuoteString(buf, err.Error(), f.EscapeUnicode)
				} else {
					buf.WriteString("null")
				}

			case TypeStringer:
				QuoteString(buf, field.Data.(fmt.Stringer).String(), f.EscapeUnicode)

			case TypeOutputer:
				field.Data.(Outputer).Output(buf)

			case TypeInterface:
				json.NewEncoder(buf).Encode(field.Data)

			case TypeStack:
				// skip types

			default:
				panic(fmt.Sprintf("unknown field type: %d", field.Type))

			}
		}

		buf.WriteByte('}')
	}

	buf.WriteByte('}')
}
