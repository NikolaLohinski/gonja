package pyint

import (
	"math/bits"
	"strconv"
)

type PyInt int

func New(i int) PyInt {
	return PyInt(i)
}

func (pi PyInt) Int() int {
	return int(pi)
}

func (pi PyInt) String() string {
	return strconv.Itoa(int(pi))
}

func (pi PyInt) BitLength() int {
	return bits.Len(uint(pi))
}

func (pi PyInt) BitCount() int {
	return bits.OnesCount(uint(pi))
}

func (pi PyInt) AsIntegerRatio() (int, int) {
	return int(pi), 1
}

func (pi PyInt) IsInteger() bool {
	return true
}

//func (pi PyInt) ToBytes(length=1, byteorder='big', *, signed=False) {
//
//}

//func (pi PyInt) FromBytes(bytes, byteorder='big', *, signed=False) {
//
//}
