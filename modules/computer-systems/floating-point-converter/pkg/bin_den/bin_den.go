package bin_den

import (
	"fmt"
	"math"
	"strconv"
)

type BtodStep struct {
	BitVal rune
	Exp    int
}

func (s BtodStep) String() string {
	return fmt.Sprintf("%c x (2 ^ %d) = %s", s.BitVal, s.Exp, fmt32(s.Calc()))
}

func (s BtodStep) Calc() float64 {
	switch s.BitVal {
	case '0':
		return 0
	case '1':
		return math.Pow(2, float64(s.Exp))
	}
	// invariant
	panic(fmt.Sprintf("bad step: %#v", s))
}

type BtodSteps struct {
	Bits  string
	Steps []BtodStep
}

func (ss BtodSteps) StepStrings() []string {
	ret := make([]string, len(ss.Steps))
	for i, s := range ss.Steps {
		ret[i] = s.String()
	}
	return ret
}

func (ss BtodSteps) Calc() float64 {
	var ret float64 = 0
	for _, s := range ss.Steps {
		ret += s.Calc()
	}
	return ret
}

// Btod converts a strings of bits into a base10 int.
func Btod(bits string) BtodSteps {
	return btod(bits, len(bits)-1)
}

// BtodFraction converts a strings of bits into a base10 fraction.
func BtodFraction(bits string) BtodSteps {
	return btod(bits, -1)
}

func btod(bits string, exp int) BtodSteps {
	ret := BtodSteps{Bits: bits}
	for _, b := range bits {
		ret.Steps = append(ret.Steps, BtodStep{
			BitVal: b,
			Exp:    exp,
		})
		exp--
	}
	return ret
}

func fmt32(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
