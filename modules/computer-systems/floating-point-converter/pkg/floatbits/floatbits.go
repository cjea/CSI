package floatbits

import (
	"fmt"
	"strings"
)

type FloatBits string

func New(s string) (FloatBits, error) {
	ret := FloatBits(s)
	if err := ret.Validate(); err != nil {
		return "", err
	}
	return ret, nil
}

func (f FloatBits) Validate() error {
	for pos, bit := range f {
		if bit == '0' || bit == '1' {
			continue
		}
		return fmt.Errorf(
			"char at position %d ('%c') is not a valid bit",
			pos, bit,
		)
	}
	if l := len(f); l != 32 {
		return fmt.Errorf("received %d bits, need 32\n", l)
	}
	return nil
}

const (
	SIGN_NEGATIVE = "NEGATIVE"
	SIGN_POSITIVE = "POSITIVE"
)

func (f FloatBits) Sign() string {
	if f[0] == '0' {
		return SIGN_POSITIVE
	} else {
		return SIGN_NEGATIVE
	}
}

func (f FloatBits) BiasedExp() string {
	return string(f[1:9])
}

func (f FloatBits) Mantissa() string {
	return string(f[9:])
}

// Special indicates that the float is either infinity (mantissa all 0's) or NaN
// (mantissa not all 0's).
func (f FloatBits) Special() bool {
	return f.BiasedExp() == "11111111"
}

func (f FloatBits) Infinity() bool {
	return f.Special() && f.Mantissa() == strings.Repeat("0", 23)
}

func (f FloatBits) NaN() bool {
	return f.Special() && !f.Infinity()
}
