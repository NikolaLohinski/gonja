package integration_test

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("strings", func() {
	var (
		identifier = new(string)

		environment   = new(*exec.Environment)
		loader        = new(loaders.Loader)
		configuration = new(*config.Config)

		context = new(*exec.Context)

		returnedResult = new(string)
		returnedErr    = new(error)
	)
	BeforeEach(func() {
		*identifier = "/test"
		*environment = gonja.DefaultEnvironment
		*loader = loaders.MustNewMemoryLoader(nil)
		*configuration = config.New()
	})
	JustBeforeEach(func() {
		var t *exec.Template
		t, *returnedErr = exec.NewTemplate(*identifier, *configuration, *loader, *environment)
		if *returnedErr != nil {
			return
		}
		*returnedResult, *returnedErr = t.ExecuteToString(*context)
	})
	Context("when getting a character by index", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
					[4]:   {{ value[4]   }}
					[-3]:  {{ value[-3]  }}
					[256]: {{ value[256] }}
					[-99]: {{ value[-99] }}
				`),
				})
				(*environment).Context.Set("value", "testing is fun!")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
					[4]:   i
					[-3]:  u
					[256]: 
					[-99]: 
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
	})
	Context("when getting a substring with '[...]' notation", func() {
		Context("default", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					[]:    {{ value[]    }}
					[:]:   {{ value[:]   }}
					[2:]:  {{ value[2:]  }}
					[:3]:  {{ value[:3]  }}
					[:-2]: {{ value[:-2] }}
					[-5:]: {{ value[-5:] }}
				`),
				})
				(*environment).Context.Set("value", "testing is fun!")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
					[]:    
					[:]:   testing is fun!
					[2:]:  sting is fun!
					[:3]:  tes
					[:-2]: testing is fu
					[-5:]:  fun!
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})

		Context("when accessing a raw string literal", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ "testing is fun!"[0:4] }}`,
				})
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("test", *returnedResult)
			})
		})
	})
	Context("when using native python methods", func() {
		var (
			shouldRender = func(template, result string) {
				Context(template, func() {
					BeforeEach(func() {
						*loader = loaders.MustNewMemoryLoader(map[string]string{
							*identifier: template,
						})
					})
					It("should return the expected rendered content", func() {
						By("not returning any error")
						Expect(*returnedErr).To(BeNil())
						By("returning the expected result")
						AssertPrettyDiff(result, *returnedResult)
					})
				})
			}
			shouldFail = func(template, err string) {
				Context(template, func() {
					BeforeEach(func() {
						*loader = loaders.MustNewMemoryLoader(map[string]string{
							*identifier: template,
						})
					})
					It("should return the expected error", func() {
						Expect(*returnedErr).ToNot(BeNil())
						Expect((*returnedErr).Error()).To(MatchRegexp(err))
					})
				})
			}
		)
		Context("upper", func() {
			shouldRender("{{ 'test'.upper() }}", "TEST")
			shouldFail("{{ 'test'.upper('unexpected') }}", "received 1 unexpected positional argument")
			shouldFail("{{ 'test'.upper(unexpected='even more') }}", "received 1 unexpected keyword argument: 'unexpected'")
		})
		Context("startswith", func() {
			shouldRender("{{ 'test123'.startswith('test') }}", "True")
			shouldRender("{{ 'test123'.startswith('foo') }}", "False")
			shouldRender("{{ 'test123'.startswith('123', 4) }}", "True")
			shouldRender("{{ 'test123'.startswith('st', 2, 4) }}", "True")
			shouldRender("{{ 'test123'.startswith('st123', 2, 4) }}", "False")
			shouldRender("{{ 'test123'.startswith('23', -2) }}", "True")
			shouldRender("{{ 'test123'.startswith('23', -2, 56) }}", "True")
			shouldRender("{{ 'test123'.startswith('test', -42, 56) }}", "True")
			shouldRender("{{ 'test123'.startswith('st', 2, -1) }}", "True")
			shouldRender("{{ 'test123'.startswith('st', 2, -54) }}", "False")
			shouldRender("{{ 'test123'.startswith('', 1, 3) }}", "True")
			shouldRender("{{ 'test123'.startswith('', 3, 1) }}", "False")
			shouldRender("{{ 'test123'.startswith(()) }}", "False")
			shouldRender("{{ 'test123'.startswith([]) }}", "False")
			shouldRender("{{ 'test123'.startswith(('test')) }}", "True")
			shouldRender("{{ 'test123'.startswith(['test']) }}", "True")
			shouldRender("{{ 'test123'.startswith(['foo', 'test']) }}", "True")
			shouldFail("{{ 'test123'.startswith(prefix='') }}", "missing required 1st positional argument 'prefix'")
		})
		Context("encode", func() {
			shouldRender("{{ 'test123'.encode() }}", "b'test123'")
			shouldRender("{{ 'test123'.encode(encoding='utf8') }}", "b'test123'")
			shouldRender("{{ 'test123'.encode('latin_1') }}", "b'test123'")
			shouldRender("{{ 'test123'.encode(errors='ignore') }}", "b'test123'")
			shouldRender("{{ 'test123'.encode('iso8859-1', 'ignore') }}", "b'test123'")
			shouldFail("{{ 'test123'.encode('iso8859-1', encoding='utf8') }}", "received 1 unexpected keyword argument: 'encoding'")
		})
		Context("encode", func() {
			shouldRender("{{ ';'.join(['s', 'b', '3']) }}", "s;b;3")
		})
		Context("format", func() {
			shouldRender("{{ 'foo={:,=-10.5G}'.format(77.11121111111112) }}", "foo=,,,,77.111")
		})
		Context("when concatenating strings with the '+' operator", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: `{{ "one" + " " + "two" }}`,
				})
			})
			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				AssertPrettyDiff("one two", *returnedResult)
			})
		})
		Context("https://github.com/NikolaLohinski/gonja/issues/25", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
					{%- set ns = namespace(found=false) -%}
					{%- for v in ["a", "b", "c"] -%}
					      {%- if v.startswith("b") -%}
					        {% set ns.found=true -%}
					      {%- endif -%}
					{%- endfor -%}
					{{ ns.found }}
					{{ name }}
					`),
				})
				(*environment).Context.Set("name", "bob")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := heredoc.Doc(`
				True
				bob
				`)
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
		Context("https://github.com/NikolaLohinski/gonja/issues/45", func() {
			BeforeEach(func() {
				*loader = loaders.MustNewMemoryLoader(map[string]string{
					*identifier: heredoc.Doc(`
						{%- if name.startswith('b') -%}
						hello world
						{%- endif -%}
					`),
				})
				(*configuration).StrictUndefined = true
				(*environment).Context.Set("name", "bob")
			})

			It("should return the expected rendered content", func() {
				By("not returning any error")
				Expect(*returnedErr).To(BeNil())
				By("returning the expected result")
				expected := "hello world"
				AssertPrettyDiff(expected, *returnedResult)
			})
		})
	})
})
