package psgc

import (
	"errors"
	"strconv"
)

var (
	InvalidPSGCCode = errors.New("invalid PSGC code")
)

type CodeResolver struct {
	Code string
}

func (c CodeResolver) Resolve() (string, error) {
	_, err := strconv.Atoi(c.Code)

	if err != nil {
		return "", InvalidPSGCCode
	}

	if c.Code == "" || len(c.Code) != 10 {
		return "", InvalidPSGCCode
	}

	count := 0

	for i := len(c.Code) - 1; i >= 0; i-- {
		if c.Code[i] == '0' {
			count++
			continue
		}

		break // Stop counting when a non-zero digit is encountered
	}

	switch count {
	case 8:
		return "Region", nil
	case 5:
		return "Province", nil
	case 3:
		return "City/Municipality/Sub Municipality", nil
	case 0:
		return "Barangay", nil
	default:
		return "", InvalidPSGCCode
	}
}
