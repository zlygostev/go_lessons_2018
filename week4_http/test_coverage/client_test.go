package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getInt(req string) (val int) {
	val, err := strconv.Atoi(req)
	if err != nil {
		panic(err)
	}
	return val
}

const (
	OrderByName = iota
	OrderByID
	OrderByAge
	ErrorOrderField
)

// type XMLContent struct {
// 	Version string   `xml:"version,attr"`
// 	Users   XMLUsers `xml:"root"`
// }

type XMLUsers struct {
	XMLName  xml.Name
	XMLUsers []XMLUser `xml:"row"`
}

type XMLUser struct {
	XMLName xml.Name `xml:"row"`
	ID      int      `xml:"id"`
	//GUID string `xml:"guid"`
	// IsAcive bool `xml:"isActive"`
	// Balabce string `xml:"balance"`
	// Picture string `xml:"picture"`
	Age int `xml:"age"`
	// EyeColor string `xml:"eyeColor"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	// Company string `xml:"company"`
	// Email string `xml:"email"`
	// Phone string `xml:"Phone"`
	// Address string `xml:"Address"`
	About string `xml:"about"`
	// <registered>2017-02-05T06:23:27 -03:00</registered>
	// <favoriteFruit>apple</favoriteFruit>
}

func PreprocessingParams(w http.ResponseWriter, r *http.Request, req *SearchRequest) (int, error) {
	accessToken := r.Header.Get("AccessToken")
	if accessToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		jsonErrString, _ := json.Marshal(SearchErrorResponse{Error: "Unauthorized"})
		w.Write(jsonErrString)
		return 0, errors.New("StatusUnauthorized")
	}

	req.Limit = getInt(r.URL.Query().Get("limit"))

	req.Offset = getInt(r.URL.Query().Get("offset"))
	if req.Offset < 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return 0, errors.New("StatusInternalServerError")
	}
	// look out through Name and About fields
	req.Query = r.URL.Query().Get("query")
	if req.Query == "StatusInternalServerError" {
		w.WriteHeader(http.StatusInternalServerError)
		return 0, errors.New("StatusInternalServerError")
	}
	//Id, Age and Name
	req.OrderField = r.URL.Query().Get("order_field")
	order := OrderByName
	switch req.OrderField {
	case "":
		fallthrough
	case "Name":
		order = OrderByName
	case "Id":
		order = OrderByID
	case "Age":
		order = OrderByAge
	default:
		order = ErrorOrderField
		// Unsupported order type
		if req.Query == "badJSON" {
			w.WriteHeader(http.StatusBadRequest)
			notJSON := ":}}}"
			w.Write([]byte(notJSON))
		} else if req.Query == "unknownError" {
			w.WriteHeader(http.StatusBadRequest)
			jsonErrString, _ := json.Marshal(SearchErrorResponse{Error: "NotWaitingMe?!"})
			w.Write(jsonErrString)
		} else if req.Query == "badJSONBodyInSuccess" {
			notJSON := ">:}}}"
			w.Write([]byte(notJSON))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			jsonErrString, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
			w.Write(jsonErrString)
		}
		return 0, errors.New("Bad OrderField")
	}
	req.OrderBy = getInt(r.URL.Query().Get("order_by"))
	switch req.OrderBy {
	case OrderByAsc:
		fallthrough
	case OrderByAsIs:
		fallthrough
	case OrderByDesc:

	default:
		jsonErrString, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
		w.Write(jsonErrString)
		return 0, errors.New("Bad OrderBy")
	}
	return order, nil
}

// By is the type of a "less" function that defines the ordering of its elements arguments.
type By func(p1, p2 *XMLUser) bool

// XMLUsersSorter joins a By function and a slice of FileInfo to be sorted.
type XMLUsersSorter struct {
	users []XMLUser
	by    func(p1, p2 *XMLUser) bool // Closure used in the Less method.
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(users []XMLUser) {
	dt := &XMLUsersSorter{
		users: users,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(dt)
}

// Len is part of sort.Interface.
func (s *XMLUsersSorter) Len() int {
	return len(s.users)
}

// Swap is part of sort.Interface.
func (s *XMLUsersSorter) Swap(i, j int) {
	s.users[i], s.users[j] = s.users[j], s.users[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *XMLUsersSorter) Less(i, j int) bool {
	return s.by(&s.users[i], &s.users[j])
}

// func sortByName(filesInfo []os.FileInfo) (err error) {
// 	name := func(p1, p2 *os.FileInfo) bool {
// 		return (*p1).Name() < (*p2).Name()
// 	}
// 	// Sort the files by the various criteria.
// 	By(name).Sort(filesInfo)

// 	return err
// }
func getOrdFunction(orderBy int, orderField int) By {
	switch orderField {
	case OrderByID:
		return func(p1, p2 *XMLUser) bool {
			if orderBy == OrderByDesc {
				return (*p1).ID > (*p2).ID
			}
			return (*p1).ID < (*p2).ID
		}

	case OrderByAge:
		return func(p1, p2 *XMLUser) bool {
			if orderBy == OrderByDesc {
				return (*p1).Age > (*p2).Age
			}
			return (*p1).Age < (*p2).Age
		}
	case OrderByName:
		fallthrough
	default:
		return func(p1, p2 *XMLUser) bool {
			if orderBy == OrderByDesc {
				return (*p1).FirstName+(*p1).LastName > (*p2).FirstName+(*p2).LastName
			}
			return (*p1).FirstName+(*p1).LastName < (*p2).FirstName+(*p2).LastName
		}
	}
}

func MakeSort(users []XMLUser, orderBy int, orderField int) {
	if orderBy == OrderByAsIs {
		return
	}
	ordFunc := getOrdFunction(orderBy, orderField)
	By(ordFunc).Sort(users)
}

func GetFilteredUsers(users []XMLUser, query string, offset int, limit int) (result []User) {

	for _, usr := range users {
		if strings.Contains(usr.FirstName+usr.LastName, query) || strings.Contains(usr.About, query) {
			result =
				append(result,
					User{
						Id:     usr.ID,
						Name:   usr.FirstName + usr.LastName,
						Age:    usr.Age,
						About:  usr.About,
						Gender: usr.Gender,
					})
		}
	}
	if len(result) > offset {
		result = result[offset:]
		if len(result) > limit {
			result = result[:limit]
		}
	}
	return result
}

// код писать тут
func SearchServer(w http.ResponseWriter, r *http.Request) {
	req := SearchRequest{}

	orderField, err := PreprocessingParams(w, r, &req)
	if err != nil {
		return
	}
	// Read data
	xmlData, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Parse users
	xmlContent := new(XMLUsers)
	err = xml.Unmarshal(xmlData, &xmlContent)

	if err != nil {
		//fmt.Println("Can't parse xml")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Reordering
	MakeSort(xmlContent.XMLUsers, req.OrderBy, orderField)
	//query
	resUsers := GetFilteredUsers(xmlContent.XMLUsers, req.Query, req.Offset, req.Limit)
	//Send results
	usersJSON, err := json.Marshal(resUsers)
	if err != nil {
		//fmt.Println("Can't serialize json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(usersJSON)
}

func SlowServer(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond * 1010)
	//w.WriteHeader(http.StatusRequestTimeout)
}

type TestCase struct {
	Req     SearchRequest
	Result  *RequestResult
	IsError bool
}

type RequestResult struct {
	Status   int
	Err      string
	Response *SearchResponse
}

func CheckCase(t *testing.T, idx int, item TestCase, resp *SearchResponse, err error) string {
	if err != nil && !item.IsError {
		return fmt.Sprintf("[%d] unexpected error: %#v - %#v", idx, err, resp)
	}

	if err != nil && item.Result != nil {
		//fmt.Printf("[%d] To check '%s'. Expected '%#v'\n", idx, err, item.Result)
		errStr := err.Error()
		//fmt.Printf("[%d] errStr %#v\n", idx, errStr)
		if len(errStr) > len(item.Result.Err) {
			errStr = errStr[:len(item.Result.Err)]
		}
		if item.IsError && item.Result.Err != errStr {
			return fmt.Sprintf("[%d] unexpected error: '%s'. Expected '%s'", idx, errStr, item.Result.Err)
		}
	}

	if err == nil && item.IsError {
		return fmt.Sprintf("[%d] expected error %s, got nil", idx, item.Result.Err)
	}
	if resp != nil && !reflect.DeepEqual(item.Result.Response, resp) {
		if !(resp.NextPage == item.Result.Response.NextPage &&
			len(resp.Users) == len(item.Result.Response.Users)) {
			return fmt.Sprintf("[%d] wrong result, expected Next:%t, Count:%d, got Next:%t, Count:%d",
				idx, item.Result.Response.NextPage, len(item.Result.Response.Users), resp.NextPage, len(resp.Users))

		}
	}
	return ""
}

func TestFindUsersEAccess(t *testing.T) {
	case1 := TestCase{
		Req: SearchRequest{},
		Result: &RequestResult{
			Err: "Bad AccessToken",
		},
		IsError: true,
	}
	caseNum := 0

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		AccessToken: "",
		URL:         ts.URL,
	}
	result, err := c.FindUsers(case1.Req)
	errStr := CheckCase(t, caseNum, case1, result, err)
	if errStr != "" {
		t.Error(errStr)
	}
	ts.Close()
}

func TestFindUsersBadURL(t *testing.T) {

	cases := []TestCase{
		TestCase{
			Req: SearchRequest{},
			Result: &RequestResult{
				Err: "unknown error",
			},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		AccessToken: "sad",
		URL:         "",
	}
	for idx, item := range cases {
		result, err := c.FindUsers(item.Req)
		errStr := CheckCase(t, idx, item, result, err)
		if errStr != "" {
			t.Error(errStr)
		}

	}
	ts.Close()
}

func TestFindUsersTimeout(t *testing.T) {

	cases := []TestCase{
		TestCase{
			Req: SearchRequest{},
			Result: &RequestResult{
				Err: "timeout for",
			},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SlowServer))
	c := &SearchClient{
		AccessToken: "asd",
		URL:         ts.URL,
	}
	for idx, item := range cases {
		result, err := c.FindUsers(item.Req)
		errStr := CheckCase(t, idx, item, result, err)
		if errStr != "" {
			t.Error(errStr)
		}
	}
	ts.Close()
}

func TestFindUsersLimits(t *testing.T) {

	cases := []TestCase{
		TestCase{
			Req: SearchRequest{
				Limit: -1,
			},
			Result: &RequestResult{
				Err: "limit must be > 0",
			},
			IsError: true,
		},
		TestCase{
			Req: SearchRequest{
				Limit: 100,
				Query: "",
			},
			Result: &RequestResult{
				Err: "",
				Response: &SearchResponse{
					Users:    make([]User, 25),
					NextPage: true,
				},
			},
			IsError: false,
		},
		TestCase{
			Req: SearchRequest{
				Limit: 10,
				Query: "Boyd",
			},
			Result: &RequestResult{
				Err: "",
				Response: &SearchResponse{
					Users:    make([]User, 1),
					NextPage: false,
				},
			},
			IsError: false,
		},
		TestCase{
			Req: SearchRequest{
				Limit: 10,
				Query: "",
			},
			Result: &RequestResult{
				Err: "",
				Response: &SearchResponse{
					Users:    make([]User, 10),
					NextPage: true,
				},
			},
			IsError: false,
		},
		TestCase{
			Req: SearchRequest{
				Query: "StatusInternalServerError",
			},
			Result: &RequestResult{
				Err: "SearchServer fatal error",
			},
			IsError: true,
		},
		TestCase{
			Req: SearchRequest{
				Offset: -1,
			},
			Result: &RequestResult{
				Err: "offset must be > 0",
			},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		AccessToken: "asd",
		URL:         ts.URL,
	}
	for idx, item := range cases {
		result, err := c.FindUsers(item.Req)
		errStr := CheckCase(t, idx, item, result, err)
		if errStr != "" {
			t.Error(errStr)
		}
	}
	ts.Close()
}

func TestFindUsersOrderFields(t *testing.T) {

	cases := []TestCase{
		TestCase{
			Req: SearchRequest{
				OrderField: "Unexpected",
				Query:      "badJSON",
			},
			Result: &RequestResult{
				Err: "cant unpack error json",
			},
			IsError: true,
		},
		TestCase{
			Req: SearchRequest{
				OrderField: "Unexpected",
			},
			Result: &RequestResult{
				Err: "OrderFeld ",
			},
			IsError: true,
		},
		TestCase{
			Req: SearchRequest{
				OrderField: "Unexpected",
				Query:      "unknownError",
			},
			Result: &RequestResult{
				Err: "unknown bad reque",
			},
			IsError: true,
		},
		TestCase{
			Req: SearchRequest{
				OrderField: "Unexpected",
				Query:      "badJSONBodyInSuccess",
			},
			Result: &RequestResult{
				Err: "cant unpack result json:",
			},
			IsError: true,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{
		AccessToken: "asd",
		URL:         ts.URL,
	}
	for idx, item := range cases {
		result, err := c.FindUsers(item.Req)
		errStr := CheckCase(t, idx, item, result, err)
		if errStr != "" {
			t.Error(errStr)
		}
	}
	ts.Close()
}

// func TestXMLParse(t *testing.T) {

// 	xmlData, err := ioutil.ReadFile("dataset.xml")
// 	if err != nil {
// 		t.Error(http.StatusInternalServerError)
// 		return
// 	}
// 	// Parse users
// 	xmlContent := new(XMLUsers)
// 	err = xml.Unmarshal(xmlData, &xmlContent)
// 	if err != nil {
// 		t.Error("Can't parse xml")
// 		return
// 	}
// 	if len(xmlContent.XMLUsers) == 0 {
// 		fmt.Println("XML content: ", xmlContent)
// 		t.Error("Parsed struct is empty")
// 		return

// 	}
// 	fmt.Println("Xml Content: ", len(xmlContent.XMLUsers),
// 		", FN: ", xmlContent.XMLUsers[0].FirstName,
// 		", LN: ", xmlContent.XMLUsers[0].LastName,
// 		", ID: ", xmlContent.XMLUsers[0].ID,
// 		", About: ", xmlContent.XMLUsers[0].About,
// 		", Age: ", xmlContent.XMLUsers[0].Age,
// 		", Gender: ", xmlContent.XMLUsers[0].Gender)

// }
