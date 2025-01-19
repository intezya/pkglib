package pkglib

type math struct{}

func (math) Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

var Math math
