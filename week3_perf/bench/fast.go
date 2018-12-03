package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

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

func easyjson6afd136DecodeCourseraWeek3Perf(in *jlexer.Lexer, out *UserJSONT) {
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
		case "browsers":
			(out.Browsers).UnmarshalEasyJSON(in)
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
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
func easyjson6afd136EncodeCourseraWeek3Perf(out *jwriter.Writer, in UserJSONT) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"browsers\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Browsers).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserJSONT) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6afd136EncodeCourseraWeek3Perf(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserJSONT) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6afd136EncodeCourseraWeek3Perf(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserJSONT) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6afd136DecodeCourseraWeek3Perf(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserJSONT) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6afd136DecodeCourseraWeek3Perf(l, v)
}
func easyjson6afd136DecodeCourseraWeek3Perf1(in *jlexer.Lexer, out *Browsers) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Browsers, 0, 4)
			} else {
				*out = Browsers{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 string
			v1 = string(in.String())
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6afd136EncodeCourseraWeek3Perf1(out *jwriter.Writer, in Browsers) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			out.String(string(v3))
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Browsers) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6afd136EncodeCourseraWeek3Perf1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Browsers) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6afd136EncodeCourseraWeek3Perf1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Browsers) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6afd136DecodeCourseraWeek3Perf1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Browsers) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6afd136DecodeCourseraWeek3Perf1(l, v)
}

func getMD5(buf []byte, hash *[16]uint8) {
	*hash = md5.Sum(buf)
}

// Было оптимальную этой функции
func NotFastSearch(out io.Writer) {
	SlowSearch(out)
}

//UserT for JSon parsing
type UserT struct {
	id             int
	isAndroid      bool
	isMSIE         bool
	isAndroidFound bool
	isMSIEFound    bool
	EMail          string //email
	Name           string //name
	Hash           [16]uint8
	Browsers       []string //browsers
	isBrwsFound    bool
	//browser   string
}

func (usr *UserT) reset() {
	// usr.isAndroidFound = false
	// usr.isMSIEFound = false
	// usr.EMail = usr.EMail[0:0]
	// usr.Name = usr.Name[0:0]
	usr.Browsers = append([]string(nil), usr.Browsers[:0]...)
}

//easyjson:json
type Browsers []string

//easyjson:json
type UserJSONT struct {
	Browsers Browsers `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

func (v *UserJSONT) reset() {
	// v.Browsers = append([]string(nil), v.Browsers[:0]...)
	// v.Email = v.Email[:0]
	// v.Name = v.Name[:0]
}

type resultsBuilderT struct {
	foundUsers bytes.Buffer
	//uniqBrowsersMD5Set map[[16]uint8]struct{}
	uniqBrowsersMD5Slice [][16]uint8
	seenBrowsers         []string
}

// func (builder *resultsBuilderT) create(buf *[]byte) {
// 	uniqBrowsersMD5
// }

func (builder *resultsBuilderT) build(buf *[]byte, usr *UserT, usrJSON *UserJSONT) {
	//res := strings.Contains("Android")
	err := usrJSON.UnmarshalJSON(*buf)
	if err != nil {
		panic(err)
	}
	usr.isAndroidFound = false
	usr.isMSIEFound = false
	for _, browser := range usrJSON.Browsers {
		usr.isAndroid = strings.Contains(browser, "Android")
		if usr.isAndroid {
			usr.isAndroidFound = true
		}
		usr.isMSIE = strings.Contains(browser, "MSIE")
		if usr.isMSIE {
			usr.isMSIEFound = true
		}
		if !(usr.isAndroid || usr.isMSIE) {
			continue
		}

		// getMD5([]byte(browser), &usr.Hash)
		// usr.isBrwsFound = false
		// for md5 := builder.uniqBrowsersMD5Slice{
		// 	usr.isBrwsFound =
		// }
		// if _, exist := builder.uniqBrowsersMD5Set[usr.Hash]; !exist {
		// 	//			fmt.Printf("[%d] New browser %v <%s>\n", uid, hash, browser)
		// 	builder.uniqBrowsersMD5Set[usr.Hash] = struct{}{}
		// }
		usr.isBrwsFound = false
		for _, item := range builder.seenBrowsers {
			if item == browser {
				usr.isBrwsFound = true
			}
		}
		if !usr.isBrwsFound {
			builder.seenBrowsers = append(builder.seenBrowsers, browser)
		}

	}

	if usr.isAndroidFound && usr.isMSIEFound {
		usr.EMail = strings.Replace(usrJSON.Email, "@", " [at] ", -1)
		builder.foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", usr.id, usrJSON.Name, usr.EMail))
	}
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	resultsBuilder := resultsBuilderT{}
	//resultsBuilder.uniqBrowsersMD5Set = make(map[[16]uint8]struct{})
	resultsBuilder.seenBrowsers = make([]string, 0, 115)

	var line []byte
	usr := UserT{}
	usrJson := UserJSONT{}
	usr.id = 0
	for {
		line, err = reader.ReadBytes('\n')
		if err != nil {
			break
		}

		usr.reset()
		//usrJson.reset()
		resultsBuilder.build(&line, &usr, &usrJson)
		// usrJson.Browsers = []string{`"SonyEricssonK810i/R1KG Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1"`,
		// 	`"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 920) like Gecko"`}
		// usrJson.Email = "jMatthews@Edgewire.mil"
		// usrJson.Name = "Donna Richardson"

		//fmt.Println("Serialized:", string(buf))

		usr.id++
	}
	fmt.Fprintf(out, "found users:\n%s\n", resultsBuilder.foundUsers.String())
	//fmt.Fprintln(out, "Total unique browsers", len(resultsBuilder.uniqBrowsersMD5Set))
	fmt.Fprintln(out, "Total unique browsers", len(resultsBuilder.seenBrowsers))
	if err != io.EOF {
		fmt.Printf(" > Failed!: %v\n", err)
	}
}

// func main() {
// 	out := os.Stdout
// 	FastSearch(out)

// }
