// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsontmp

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2bbe306fDecodeCourseraWeek3PerfEasyjsontut(in *jlexer.Lexer, out *UserJSONT) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2bbe306fEncodeCourseraWeek3PerfEasyjsontut(out *jwriter.Writer, in UserJSONT) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserJSONT) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2bbe306fEncodeCourseraWeek3PerfEasyjsontut(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserJSONT) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2bbe306fEncodeCourseraWeek3PerfEasyjsontut(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserJSONT) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2bbe306fDecodeCourseraWeek3PerfEasyjsontut(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserJSONT) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2bbe306fDecodeCourseraWeek3PerfEasyjsontut(l, v)
}