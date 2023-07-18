package integration_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	osExec "os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/nikolalohinski/gonja"
	"github.com/nikolalohinski/gonja/exec"
	"github.com/nikolalohinski/gonja/loaders"
	"github.com/yargevad/filepathx"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("miscellaneous templates", func() {
	const (
		fixturesDir  = "testdata/"
		testCasesDir = "testcases/"
	)
	var (
		template    = new(string)
		environment = new(*gonja.Environment)
		context     = new(gonja.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		rand.Seed(42)
		*environment = gonja.NewEnvironment(gonja.NewConfig(), loaders.MustNewFileSystemLoader(fixturesDir))
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = (*environment).FromString(*template)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.Execute(*context)
	})
	Context("nominal", func() {
		files := MustReturn(filepathx.Glob(path.Join(testCasesDir, "**/*.tpl"))).([]string)
		for _, p := range files {
			path := p
			name := strings.TrimSuffix(strings.TrimPrefix(path, testCasesDir), ".tpl")
			Context(strings.Join(strings.Split(name, string(filepath.Separator)), " : "), func() {
				BeforeEach(func() {
					*environment = gonja.NewEnvironment(gonja.NewConfig(), loaders.MustNewFileSystemLoader(filepath.Dir(path)))
					*template = string(MustReturn(ioutil.ReadFile(path)).([]byte))
					*context = Fixtures
				})
				It("should render the expected content", func() {
					By("not returning any error")
					Expect(*returnedErr).To(BeNil())
					By("returning the correct result")
					expected := string(MustReturn(ioutil.ReadFile(path + ".out")).([]byte))
					edits := myers.ComputeEdits("expected", expected, *returnedResult)
					diffs := gotextdiff.ToUnified("expected", "got", expected, edits)
					Expect(diffs.Hunks).To(BeEmpty(), "\n"+fmt.Sprint(diffs))
				})
			})
		}
	})
	Context("miscellaneous templates", func() {
		for _, testCase := range []struct {
			description string
			template    string
			expected    string
			context     gonja.Context
		}{
			{
				"when the template is just data",
				"Hello world!",
				"Hello world!",
				nil,
			},
			{
				"when the template is empty",
				"",
				"",
				nil,
			},
			{
				"when the template is using quoting strings",
				heredoc.Doc(`
				Variables
				{{ "hello" }}
				{{ 'hello' }}
				{{ "hell'o" }}

				Filters
				{{ 'Test'|slice('1:3') }}
				{{ '<div class=\"foo\"><ul class=\"foo\"><li class=\"foo\"><p class=\"foo\">This is a long test.</p></li></ul></div>'|safe }}
				{{ '<a name="link"><p>This </a>is a long test.</p>'|safe }}

				Tags
				{% if 'Text' in map.struct %}Text attribute in map.struct{% endif %}

				Functions
				{{ func_variadic('hello') }}
			`),
				heredoc.Doc(`
				Variables
				hello
				hello
				hell'o

				Filters
				es
				<div class="foo"><ul class="foo"><li class="foo"><p class="foo">This is a long test.</p></li></ul></div>
				<a name="link"><p>This </a>is a long test.</p>

				Tags
				Text attribute in map.struct

				Functions
				hello
			`),
				gonja.Context{
					"map": map[string]interface{}{
						"struct": struct{ Text string }{
							Text: "does not matter",
						},
					},
					"func_variadic": func(msg string, args ...interface{}) string {
						return fmt.Sprintf(msg, args...)
					},
				},
			},
			{
				"when the template is using macros",
				heredoc.Doc(`
				Begin
				{% macro greetings(to, from=name, name2="guest") %}
				Greetings to {{ to }} from {{ from }}. Howdy, {% if name2 == "guest" %}anonymous guest{% else %}{{ name2 }}{% endif %}!
				{% endmacro %}
				{{ greetings('') }}
				{{ greetings(10) }}
				{{ greetings("john") }}
				{{ greetings("john", "michelle") }}
				{{ greetings("john", "michelle", "johann") }}

				{% macro test2(loop, value) %}map[{{ loop.index0 }}] = {{ value }}{% endmacro %}
				{% for item in misc_list %}
				{{ test2(loop, item) }}{% endfor %}

				issue #39 (deactivate auto-escape of macros)
				{% macro html_test(name) %}
				<p>Hello {{ name }}.</p>
				{% endmacro %}
				{{ html_test("Max") }}

				Importing macros
				{% from "macro.helper" import imported_macro, imported_macro as renamed_macro, imported_macro as html_test %}
				{{ imported_macro("User1") }}
				{{ renamed_macro("User2") }}
				{{ html_test("Max") }}

				Chaining macros{% from "macro2.helper" import greeter_macro %}
				{{ greeter_macro() }}
				End
			`),
				heredoc.Doc(`
                Begin
                
                
                Greetings to  from john doe. Howdy, anonymous guest!
                
                
                Greetings to 10 from john doe. Howdy, anonymous guest!
                
                
                Greetings to john from john doe. Howdy, anonymous guest!
                
                
                Greetings to john from michelle. Howdy, anonymous guest!
                
                
                Greetings to john from michelle. Howdy, johann!
                
                
                
                
                map[0] = Hello
                map[1] = 99
                map[2] = 3.14
                map[3] = good
                
                issue #39 (deactivate auto-escape of macros)
                
                
                <p>Hello Max.</p>
                
                
                Importing macros
                
                <p>Hey User1!</p>
                <p>Hey User2!</p>
                <p>Hey Max!</p>
                
                Chaining macros
                
                
                One greeting: <p>Hey Dirk!</p> - <p>Hello mate!</p>
                
                End
			`),
				gonja.Context{
					"name":      "john doe",
					"misc_list": []interface{}{"Hello", 99, 3.14, "good"},
				},
			},
			{
				"when the template is using function call wrappers",
				heredoc.Doc(`
				{{ func_add(func_add(5, 15), number) + 17 }}
				{{ func_add_iface(func_add_iface(5, 15), number) + 17 }}
				{{ func_variadic("hello") }}
				{{ func_variadic("hello, %s", name) }}
				{{ func_variadic("%d + %d %s %d", 5, number, "is", 49) }}
				{{ func_variadic_sum_int() }}
				{{ func_variadic_sum_int(1) }}
				{{ func_variadic_sum_int(1, 19, 185) }}
				{{ func_variadic_sum_int2() }}
				{{ func_variadic_sum_int2(2) }}
				{{ func_variadic_sum_int2(1, 7, 100) }}
				{{ func_with_varargs(1) }}
				{{ func_with_varargs(1, 2) }}
				{{ func_with_varargs(a='a') }}
				{{ func_with_varargs(a='a', b='b') }}
				{{ func_with_varargs(1, 2, 3, a='a', b='b', c='c') }}
			`),
				heredoc.Doc(`
				79
				79
				hello
				hello, john doe
				5 + 42 is 49
				0
				1
				205
				0
				2
				108
				VarArgs(args=[1], kwargs={})
				VarArgs(args=[1, 2], kwargs={})
				VarArgs(args=[], kwargs={a="a"})
				VarArgs(args=[], kwargs={a="a", b="b"})
				VarArgs(args=[1, 2, 3], kwargs={a="a", b="b", c="c"})
			`),
				gonja.Context{
					"name":   "john doe",
					"number": 42,
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
			},
			{
				"when the template is using quoting strings",
				heredoc.Doc(`
                <!DOCTYPE html>
                {# A more complex template using gonja #}
                <html>
                
                <head>
                	<title>My blog page</title>
                </head>
                
                <body>
                	<h1>Blogpost</h1>
                	<div id="content">
                		{{ text|safe }}
                	</div>
                
                	<h1>Comments</h1>
                
                	{% for comment in nested.comments %}
                		<h2>{{ loop.index }}. Comment ({{ loop.revindex}} comment{% if loop.revindex > 1 %}s{% endif %} left)</h2>
                		<p>From: {{ comment.Author.Name }} ({% if comment.Author.Validated %}validated{% else %}not validated{% endif %})</p>
                
                		{% if is_admin(comment.Author) %}
                			<p>This user is an admin!</p>
                		{% else %}
                			<p>This user is not admin!</p>
                		{% endif %}
                
                		<p>Written {{ comment.Date }}</p>
                		<p>{{ comment.Text|striptags }}</p>
                	{% endfor %}
                </body>
                
                </html>
			`),
				heredoc.Doc(`
                <!DOCTYPE html>

                <html>
                
                <head>
                	<title>My blog page</title>
                </head>
                
                <body>
                	<h1>Blogpost</h1>
                	<div id="content">
                		<h2>Hello!</h2><p>Welcome to my new blog page. I'm using gonja which supports {{ variables }} and {% tags %}.</p>
                	</div>
                
                	<h1>Comments</h1>
                
                	
                		<h2>1. Comment (3 comments left)</h2>
                		<p>From: user1 (validated)</p>
                
                		
                			<p>This user is not admin!</p>
                		
                
                		<p>Written 2014-06-10 15:30:15 +0000 UTC</p>
                		<p>"gonja is nice!"</p>
                	
                		<h2>2. Comment (2 comments left)</h2>
                		<p>From: user2 (validated)</p>
                
                		
                			<p>This user is an admin!</p>
                		
                
                		<p>Written 2011-03-21 08:37:56.000000012 +0000 UTC</p>
                		<p>comment2 with unsafe tags in it</p>
                	
                		<h2>3. Comment (1 comment left)</h2>
                		<p>From: user3 (not validated)</p>
                
                		
                			<p>This user is not admin!</p>
                		
                
                		<p>Written 2014-06-10 15:30:15 +0000 UTC</p>
                		<p>hello! there</p>
                	
                </body>
                
                </html>
			`),
				gonja.Context{
					"text": "<h2>Hello!</h2><p>Welcome to my new blog page. I'm using gonja which supports {{ variables }} and {% tags %}.</p>",
					"is_admin": func(u *struct {
						Name      string
						Validated bool
					}) bool {
						for _, a := range adminList {
							if a == u.Name {
								return true
							}
						}
						return false
					},
					"nested": map[string]interface{}{
						"comments": []struct {
							Author *struct {
								Name      string
								Validated bool
							}
							Date time.Time
							Text string
						}{
							{
								Author: &struct {
									Name      string
									Validated bool
								}{
									"user1",
									true,
								},
								Date: time.Date(2014, 06, 10, 15, 30, 15, 0, time.UTC),
								Text: `"gonja is nice!"`,
							},
							{
								Author: &struct {
									Name      string
									Validated bool
								}{
									"user2",
									true,
								},
								Date: time.Date(2011, 03, 21, 8, 37, 56, 12, time.UTC),
								Text: "comment2 with <script>unsafe</script> tags in it",
							},
							{
								Author: &struct {
									Name      string
									Validated bool
								}{
									"user3",
									false,
								},
								Date: time.Date(2014, 06, 10, 15, 30, 15, 0, time.UTC),
								Text: "<b>hello!</b> there",
							},
						},
					},
				},
			},
		} {
			t := testCase
			Context(t.description, func() {
				BeforeEach(func() {
					*template = t.template
					*context = t.context
				})
				It("should return the expected rendered content", func() {
					By("not returning any error")
					Expect(*returnedErr).To(BeNil())
					By("not returning the correct result")
					edits := myers.ComputeEdits("expected", t.expected, *returnedResult)
					diffs := gotextdiff.ToUnified("expected", "got", t.expected, edits)
					Expect(diffs.Hunks).To(BeEmpty(), "\n"+fmt.Sprint(diffs))
				})
			})
			Context("when accessing gonja's version", func() {
				BeforeEach(func() {
					*template = heredoc.Doc(`
					{%- if "v" ~ gonja.version != CI_COMMIT_TAG -%}
					v{{- gonja.version }} != {{ CI_COMMIT_TAG }}
					{%- else -%}
					"v" ~ gonja.version == CI_COMMIT_TAG
					{%- endif %}
				`)
					cmd := osExec.Command("git", "describe", "--tags", "--abbrev=0")
					buf := new(bytes.Buffer)
					cmd.Stdout = buf
					Must(cmd.Run())
					(*environment).Globals.Set("CI_COMMIT_TAG", strings.TrimSpace(buf.String()))
				})
				It("should return the expected rendered content", func() {
					By("not returning any error")
					Expect(*returnedErr).To(BeNil())
					By("not returning the correct result")
					expected := heredoc.Doc(`
					"v" ~ gonja.version == CI_COMMIT_TAG
				`)
					edits := myers.ComputeEdits("expected", expected, *returnedResult)
					diffs := gotextdiff.ToUnified("expected", "got", expected, edits)
					Expect(diffs.Hunks).To(BeEmpty(), "\n"+fmt.Sprint(diffs))
				})
			})
		}
	})
})
