package pystring

type PyString string

func New(s string) PyString {
	return PyString(s)
}

func (s PyString) String() string {
	return string(s)
}
