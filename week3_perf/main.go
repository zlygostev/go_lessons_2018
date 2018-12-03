package main

import "fmt"

func main() {
	data := []byte(`{"browsers":["Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36","LG-LX550 AU-MIC-LX550/2.0 MMP/2.0 Profile/MIDP-2.0 Configuration/CLDC-1.1","Mozilla/5.0 (Android; Linux armv7l; rv:10.0.1) Gecko/20100101 Firefox/10.0.1 Fennec/10.0.1","Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; MATBJS; rv:11.0) like Gecko"],"company":"Flashpoint","country":"Dominican Republic","email":"JonathanMorris@Muxo.edu","job":"Programmer Analyst #{N}","name":"Sharon Crawford","phone":"176-88-49"}`)
	//data := []byte(`{email:"abc@as.ru", name:"Hey", browsers:[Moz FF]}`)
	//data := []byte(`{"Browsers":["Moz","FF"],"Email":"abc@as.ru","Name":"Hey"}`)
	v := UserJSONT{}
	v.UnmarshalJSON(data)
	fmt.Println(v.Browsers, "name: ", v.Name, "e-mail:", v.Email)
	v = UserJSONT{}
	v.Email = "abc@as.ru"
	v.Name = "Hey"
	v.Browsers = Browsers{"Moz", "FF"}
	fmt.Println(v)
	res, err := v.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	// data = []byte(`[{ "range":213,"Val":4}, {"str":3445,"Val":2}]`)
	// v1 := ItemList{}
	// v1.UnmarshalJSON(data)
	// fmt.Println("ItemList:", len(v1), v1[0].Val, v1[1].Val)
	// res, err1 := v1.MarshalJSON()
	// fmt.Println(err1)
	// if err1 != nil {
	// 	panic(err1)
	// }

	fmt.Println(string(res))

}
