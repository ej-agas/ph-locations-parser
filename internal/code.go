package internal

import (
	"errors"
	"strconv"
	"strings"
)

var (
	InvalidPSGCCode = errors.New("invalid PSGC code")
)

type psgc struct {
	Code string
}

func NewPSGC(code string) (*psgc, error) {
	_, err := strconv.Atoi(code)

	if err != nil {
		return &psgc{}, InvalidPSGCCode
	}

	if code == "" || len(code) != 10 {
		return &psgc{}, InvalidPSGCCode
	}

	return &psgc{Code: code}, nil
}

func (p psgc) Region() string {
	code := p.Code[:2]

	code += strings.Repeat("0", 8)

	return code
}

func (p psgc) Province() string {
	code := p.Code[:5]

	code += strings.Repeat("0", 5)

	return code
}

func (p psgc) CityOrMunicipality() string {
	code := p.Code[:7]

	code += strings.Repeat("0", 3)

	return code
}

func (p psgc) Barangay() string {
	return p.Code
}
