package jsontmp

// package main

//easyjson:json
type ItemList []Item

//easyjson:json
type Item struct {
	Val int
}

// func TestUnpack() {
// 	data := []byte(`[{"Val":1}, {"Val":2}]`)
// 	v := ItemList{}

// 	v.UnmarshalJSON(data)

// 	fmt.Println(v)
// }

// func main() {
// 	usrJSON := UserJSONT{}
// 	usrJSON.browsers = []string{`"SonyEricssonK810i/R1KG Browser/NetFront/3.3 Profile/MIDP-2.0 Configuration/CLDC-1.1"`,
// 		`"Mozilla/5.0 (Windows Phone 8.1; ARM; Trident/7.0; Touch; rv:11.0; IEMobile/11.0; NOKIA; Lumia 920) like Gecko"`}
// 	usrJSON.company = "Jabbercube"
// 	usrJSON.country = "Indonesia"
// 	usrJSON.email = "jMatthews@Edgewire.mil"
// 	usrJSON.name = "Donna Richardson"

// 	buf, err := (usrJSON).MarshalJSON()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Serialized:'%s'\n", string(buf))
// 	usrJSON.reset()
// 	jsn := `{"browsers":["Mozilla/5.0 (X11; Linux i686; rv:12.0) Gecko/20120502 Firefox/12.0 SeaMonkey/2.9.1","EmailWolf 1.00","Mozilla/5.0 (Windows NT 6.2; ARM; Trident/7.0; Touch; rv:11.0; WPDesktop; NOKIA; Lumia 635) like Gecko","Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.3; Trident/7.0; .NET4.0E; .NET4.0C)"],"company":"Mybuzz","country":"Tajikistan","email":"yWilliamson@Muxo.name","job":"Junior Executive","name":"Anna Reid","phone":"945-73-11"}`

// 	err = (&usrJSON).UnmarshalJSON([]byte(jsn))
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Parsed result:'%v'\n", usrJSON)

// 	//out := os.Stdout
// 	//FastSearch(out)
// }
