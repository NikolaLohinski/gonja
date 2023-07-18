package integration_test

import (
	"fmt"
	"sort"

	"strings"
	"time"

	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/exec"
)

var adminList = []string{"user2"}

var time1 = time.Date(2014, 06, 10, 15, 30, 15, 0, time.UTC)
var time2 = time.Date(2011, 03, 21, 8, 37, 56, 12, time.UTC)

type post struct {
	Text    string
	Created time.Time
}

type user struct {
	Name      string
	Validated bool
}

func (u *user) String() string {
	return u.Name
}

type comment struct {
	Author *user
	Date   time.Time
	Text   string
}

type person struct {
	FirstName string
	LastName  string
	Gender    string
}

func isAdmin(u *user) bool {
	for _, a := range adminList {
		if a == u.Name {
			return true
		}
	}
	return false
}

func (u *user) IsAdmin() *exec.Value {
	return exec.AsValue(isAdmin(u))
}

func (u *user) IsAdmin2() bool {
	return isAdmin(u)
}

func (p *post) String() string {
	return ":-)"
}

/*
 * End setup sandbox
 */

var Fixtures = gonja.Context{
	"number": 11,
	"simple": map[string]interface{}{
		"number":                   42,
		"name":                     "john doe",
		"included_file":            "INCLUDES.helper",
		"included_file_not_exists": "INCLUDES.helper.not_exists",
		"nil":                      nil,
		"uint":                     uint(8),
		"float":                    float64(3.1415),
		"str":                      "string",
		"chinese_hello_world":      "你好世界",
		"bool_true":                true,
		"bool_false":               false,
		"newline_text": `this is a text
with a new line in it`,
		"long_text": `This is a simple text.

This too, as a paragraph.
Right?

Yep!`,
		"escape_js_test":     `escape sequences \r\n\'\" special chars "?!=$<>`,
		"one_item_list":      []int{99},
		"multiple_item_list": []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55},
		"unsorted_int_list":  []int{192, 581, 22, 1, 249, 9999, 1828591, 8271},
		"fixed_item_list":    [...]int{1, 2, 3, 4},
		"misc_list":          []interface{}{"Hello", 99, 3.14, "good"},
		"escape_text":        "This is \\a Test. \"Yep\". 'Yep'.",
		"xss":                "<script>alert(\"uh oh\");</script>",
		"intmap": map[int]string{
			1: "one",
			5: "five",
			2: "two",
		},
		"strmap": map[string]string{
			"abc": "def",
			"bcd": "efg",
			"zab": "cde",
			"gh":  "kqm",
			"ukq": "qqa",
			"aab": "aba",
		},
		"casedStrmap": map[string]string{
			"a": "a",
			"B": "B",
			"c": "c",
			"D": "D",
			"e": "e",
			"F": "F",
		},
		"func_add": func(a, b int) int {
			return a + b
		},
		"func_add_iface": func(a, b interface{}) interface{} {
			return a.(int) + b.(int)
		},
		"func_variadic": func(msg string, args ...interface{}) string {
			return fmt.Sprintf(msg, args...)
		},
		"func_variadic_sum_int": func(args ...int) int {
			// Create a sum
			s := 0
			for _, i := range args {
				s += i
			}
			return s
		},
		"func_variadic_sum_int2": func(args ...*exec.Value) *exec.Value {
			// Create a sum
			s := 0
			for _, i := range args {
				s += i.Integer()
			}
			return exec.AsValue(s)
		},
		"func_with_varargs": func(params *exec.VarArgs) *exec.Value {
			// arg := params.args[0]
			argsAsStr := []string{}
			for _, arg := range params.Args {
				argsAsStr = append(argsAsStr, arg.String())
			}
			kwargsAsStr := []string{}
			for key, value := range params.KwArgs {
				v := value.String()
				if value.IsString() {
					v = "\"" + v + "\""
				}
				pair := []string{key, v}
				kwargsAsStr = append(kwargsAsStr, strings.Join(pair, "="))
			}
			sort.Strings(kwargsAsStr)
			args := strings.Join(argsAsStr, ", ")
			kwargs := strings.Join(kwargsAsStr, ", ")

			str := fmt.Sprintf("VarArgs(args=[%s], kwargs={%s})", args, kwargs)
			return exec.AsSafeValue(str)
		},
	},
	"complex": map[string]interface{}{
		"user": &user{
			Name:      "john doe",
			Validated: true,
		},
		"is_admin": isAdmin,
		"post": post{
			Text:    "<h2>Hello!</h2><p>Welcome to my new blog page. I'm using gonja which supports {{ variables }} and {% tags %}.</p>",
			Created: time2,
		},
		"comments": []*comment{
			&comment{
				Author: &user{
					Name:      "user1",
					Validated: true,
				},
				Date: time1,
				Text: "\"gonja is nice!\"",
			},
			&comment{
				Author: &user{
					Name:      "user2",
					Validated: true,
				},
				Date: time2,
				Text: "comment2 with <script>unsafe</script> tags in it",
			},
			&comment{
				Author: &user{
					Name:      "user3",
					Validated: false,
				},
				Date: time1,
				Text: "<b>hello!</b> there",
			},
		},
		"comments2": []*comment{
			&comment{
				Author: &user{
					Name:      "user1",
					Validated: true,
				},
				Date: time2,
				Text: "\"gonja is nice!\"",
			},
			&comment{
				Author: &user{
					Name:      "user1",
					Validated: true,
				},
				Date: time1,
				Text: "comment2 with <script>unsafe</script> tags in it",
			},
			&comment{
				Author: &user{
					Name:      "user3",
					Validated: false,
				},
				Date: time1,
				Text: "<b>hello!</b> there",
			},
		},
	},
	"persons": []*person{
		{"John", "Doe", "male"},
		{"Jane", "Doe", "female"},
		{"Akira", "Toriyama", "male"},
		{"Selina", "Kyle", "female"},
		{"Axel", "Haustant", "male"},
	},
	"groupable": []map[string]string{
		{"grouper": "group 1", "value": "value 1-1"},
		{"grouper": "group 2", "value": "value 2-1"},
		{"grouper": "group 3", "value": "value 3-1"},
		{"grouper": "group 1", "value": "value 1-2"},
		{"grouper": "group 2", "value": "value 2-2"},
		{"grouper": "group 3", "value": "value 3-2"},
		{"grouper": "group 1", "value": "value 1-3"},
		{"grouper": "group 2", "value": "value 2-3"},
		{"grouper": "group 3", "value": "value 3-3"},
	},
}
