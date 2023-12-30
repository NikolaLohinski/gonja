package exec_test

import (
	"github.com/nikolalohinski/gonja/v2/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Context("varargs", func() {
	var (
		varargs *exec.VarArgs
	)
	BeforeEach(func() {
		varargs = new(exec.VarArgs)
	})
	Context("first", func() {
		var (
			returnedValue = new(exec.Value)
		)
		JustBeforeEach(func() {
			returnedValue = varargs.First()
		})
		Context("nil if empty", func() {
			It("should return the correct value", func() {
				Expect(returnedValue.IsNil()).To(BeTrue())
			})
		})
		Context("first value", func() {
			BeforeEach(func() {
				varargs = &exec.VarArgs{Args: []*exec.Value{exec.AsValue(42)}}
			})
			It("should return the correct value", func() {
				Expect(returnedValue.Integer()).To(Equal(42))
			})
		})
	})
	Context("GetKwarg", func() {
		var (
			key      = new(string)
			fallback = new(interface{})

			returnedValue = new(exec.Value)
		)
		BeforeEach(func() {
			*key = "key"
			*fallback = "not found"
		})
		JustBeforeEach(func() {
			returnedValue = varargs.GetKwarg(*key, *fallback)
		})
		Context("default if missing", func() {
			It("should return the correct value", func() {
				Expect(returnedValue.String()).To(Equal("not found"))
			})
		})
		Context("value if found", func() {
			BeforeEach(func() {
				varargs = &exec.VarArgs{KwArgs: map[string]*exec.Value{
					"key": exec.AsValue(42),
				}}
			})
			It("should return the correct value", func() {
				Expect(returnedValue.Integer()).To(Equal(42))
			})
		})
	})
	Context("Expect", func() {
		Context("nothing", func() {
			var (
				returnedVarArgs *exec.ReducedVarArgs
			)
			JustBeforeEach(func() {
				returnedVarArgs = varargs.ExpectNothing()
			})
			for _, t := range []struct {
				desc    string
				varargs *exec.VarArgs
				err     string
			}{
				{
					"got nothing",
					&exec.VarArgs{},
					"",
				},
				{
					"got an argument",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42)}},
					`Unexpected argument '42'`,
				},
				{
					"got multiples arguments",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42), exec.AsValue(7)}},
					`Unexpected arguments '42, 7'`,
				},
				{
					"got a keyword argument",
					&exec.VarArgs{KwArgs: map[string]*exec.Value{
						"key": exec.AsValue(42),
					}},
					`Unexpected keyword argument 'key=42'`,
				},
				{
					"got multiple keyword arguments",
					&exec.VarArgs{KwArgs: map[string]*exec.Value{
						"key":   exec.AsValue(42),
						"other": exec.AsValue(7),
					}},
					`Unexpected keyword arguments 'key=42, other=7'`,
				},
				{
					"got one of each",
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key": exec.AsValue(42),
						},
					},
					`Unexpected arguments '42, key=42'`,
				},
			} {
				Context(t.desc, func() {
					BeforeEach(func() {
						varargs = t.varargs
					})
					It("should return an error", func() {
						Expect(returnedVarArgs.IsError()).To(BeTrue(), "should have returned an error")
						Expect(returnedVarArgs.Error()).To(Equal(t.err))
					})
				})
			}
		})
		Context("arguments", func() {
			var (
				returnedVarArgs *exec.ReducedVarArgs
			)

			for _, t := range []struct {
				desc    string
				varargs *exec.VarArgs
				args    int
				err     string
			}{
				{
					"got expected",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42), exec.AsValue(7)}},
					2,
					"",
				},
				{
					"got less arguments",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42)}},
					2,
					`Expected 2 arguments, got 1`,
				},
				{
					"got less arguments (singular)",
					&exec.VarArgs{},
					1,
					`Expected an argument, got 0`,
				},
				{
					"got more arguments",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42), exec.AsValue(7)}},
					1,
					`Unexpected argument '7'`,
				},
				{
					"got a keyword argument",
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key": exec.AsValue(42),
						},
					},
					1,
					`Unexpected keyword argument 'key=42'`,
				},
			} {
				Context(t.desc, func() {
					BeforeEach(func() {
						varargs = t.varargs
					})
					JustBeforeEach(func() {
						returnedVarArgs = varargs.ExpectArgs(t.args)
					})
					It("should return an error", func() {
						Expect(returnedVarArgs.IsError()).To(BeTrue(), "should have returned an error")
						Expect(returnedVarArgs.Error()).To(Equal(t.err))
					})
				})
			}
		})
		Context("keyword arguments", func() {
			var (
				returnedVarArgs *exec.ReducedVarArgs
			)

			for _, t := range []struct {
				desc    string
				varargs *exec.VarArgs
				kwargs  []*exec.KwArg
				err     string
			}{
				{
					"got expected",
					&exec.VarArgs{KwArgs: map[string]*exec.Value{
						"key":   exec.AsValue(42),
						"other": exec.AsValue(7),
					}},
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					"",
				},
				{
					"got unexpected arguments",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42), exec.AsValue(7), exec.AsValue("unexpected")}},
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					`Unexpected argument 'unexpected'`,
				},
				{
					"got an unexpected keyword argument",
					&exec.VarArgs{KwArgs: map[string]*exec.Value{
						"unknown": exec.AsValue(42),
					}},
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					`Unexpected keyword argument 'unknown=42'`,
				},
				{
					"got multiple keyword arguments",
					&exec.VarArgs{KwArgs: map[string]*exec.Value{
						"unknown": exec.AsValue(42),
						"seven":   exec.AsValue(7),
					}},
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					`Unexpected keyword arguments 'seven=7, unknown=42'`,
				},
			} {
				Context(t.desc, func() {
					BeforeEach(func() {
						varargs = t.varargs
					})
					JustBeforeEach(func() {
						returnedVarArgs = varargs.Expect(0, t.kwargs)
					})
					It("should return an error", func() {
						Expect(returnedVarArgs.IsError()).To(BeTrue(), "should have returned an error")
						Expect(returnedVarArgs.Error()).To(Equal(t.err))
					})
				})
			}
		})
		Context("mixed arguments", func() {
			var (
				returnedVarArgs *exec.ReducedVarArgs
			)

			for _, t := range []struct {
				desc     string
				varargs  *exec.VarArgs
				args     int
				kwargs   []*exec.KwArg
				expected *exec.VarArgs
				err      string
			}{
				{
					"got expected",
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue(42),
							"other": exec.AsValue(7),
						},
					},
					1,
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue(42),
							"other": exec.AsValue(7),
						},
					},
					"",
				},
				{
					"fill with default",
					&exec.VarArgs{Args: []*exec.Value{exec.AsValue(42)}},
					1,
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue("default key"),
							"other": exec.AsValue("default other"),
						},
					},
					"",
				},
				{
					"keyword as argument",
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42), exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"other": exec.AsValue(7),
						},
					},
					1,
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue(42),
							"other": exec.AsValue(7),
						},
					},
					"",
				},
				{
					"keyword submitted twice",
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42), exec.AsValue(5)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue(42),
							"other": exec.AsValue(7),
						},
					},
					1,
					[]*exec.KwArg{
						{"key", "default key"},
						{"other", "default other"},
					},
					&exec.VarArgs{
						Args: []*exec.Value{exec.AsValue(42), exec.AsValue(5)},
						KwArgs: map[string]*exec.Value{
							"key":   exec.AsValue(42),
							"other": exec.AsValue(7),
						},
					},
					`Keyword 'key' has been submitted twice`,
				},
			} {
				Context(t.desc, func() {
					BeforeEach(func() {
						varargs = t.varargs
					})
					JustBeforeEach(func() {
						returnedVarArgs = varargs.Expect(t.args, t.kwargs)
					})
					It("should return the expect var args", func() {
						if t.err != "" {
							Expect(returnedVarArgs.IsError()).To(BeTrue(), "should have returned an error")
							Expect(returnedVarArgs.Error()).To(Equal(t.err))
						} else {
							Expect(returnedVarArgs).To(Equal(t.expected))
						}
					})
				})
			}
		})
	})
})
