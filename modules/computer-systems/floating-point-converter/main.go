package main

import (
	"float-to-den/pkg/bin_den"
	"float-to-den/pkg/floatbits"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func usage(ec int) {
	fmt.Println(
		`
$ go run main.go <32-bit float>

Example:
$ go run main.go 11111111100000000000000000000010
Since the exponent is all 1's, and the mantissa is NOT all 0's, this bit pattern represents NaN (Not a Number).`)
	os.Exit(ec)
}

func main() {
	if len(os.Args) < 2 {
		usage(1)
	}
	bits, err := floatbits.New(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		usage(1)
	}

	expConversion := bin_den.Btod(bits.BiasedExp())
	biasedExp := int(expConversion.Calc())
	mantissaConversion := bin_den.BtodFraction(bits.Mantissa())
	mantissa := mantissaConversion.Calc()
	fullMantissa := 1 + mantissa
	fullMantissaWithSign := fullMantissa
	if bits.Sign() == floatbits.SIGN_NEGATIVE {
		fullMantissaWithSign *= -1
	}

	steps := make([]string, 9)

	if bits.Infinity() {
		fmt.Printf(
			"Since the exponent is all 1's, and the mantissa is all 0's, this bit pattern represents INFINITY.\n",
		)
		return
	} else if bits.NaN() {
		fmt.Printf(
			"Since the exponent is all 1's, and the mantissa is NOT all 0's, this bit pattern represents NaN (Not a Number).\n",
		)
		return
	}
	steps[0] = fmt.Sprintf(
		"The first bit of the input '%s' is '%c', which means the sign is: %s.\n",
		bits, bits[0], bits.Sign(),
	)
	steps[1] = fmt.Sprintf(
		"The subsequent 8 bits make up the biased exponent ('%s').", bits.BiasedExp(),
	)
	steps[2] = fmt.Sprintf(
		"Convert the biased exponent to denary:\n\n\t%s\n\nFor a result of %d.\n",
		strings.Join(expConversion.StepStrings(), "\n\t"),
		biasedExp,
	)
	steps[3] = fmt.Sprintf(
		"Since that is the biased exponent, we need to remove the bias (subtract 127):\n=>\t"+
			"Real exponent = %d - 127 = %d\n", biasedExp, biasedExp-127,
	)
	steps[4] = fmt.Sprintf(
		"The final 23 bits make up the mantissa ('%s').The most significat bit is 1/2, then 1/4, 1/8, etc.", bits.Mantissa(),
	)
	steps[5] = fmt.Sprintf(
		"Convert the mantissa to denary:\n\n\t%s\n\nFor a result of %s.\n",
		strings.Join(mantissaConversion.StepStrings(), "\n\t"),
		fmt32(mantissa),
	)
	steps[6] = fmt.Sprintf(
		"Every floating point mantissa begins with an implicit 1:\n=>\t"+
			"Real mantissa = 1 + %s = %s\n", fmt32(mantissa), fmt32(fullMantissa),
	)
	if bits.Sign() == floatbits.SIGN_NEGATIVE {
		steps[7] = fmt.Sprintf(
			"Since the sign was negative, negate the mantissa:\n=>\t"+
				"Realest mantissa = %s\n", fmt32(fullMantissaWithSign),
		)
	}
	steps[8] = fmt.Sprintf(
		"Multiply the mantissa by 2 ^ exponent for the final result:\n\n"+
			"%s\n\t\t= %s * (2 ^ %d)\n\t\t= %s",
		bits,
		fmt32(fullMantissaWithSign), biasedExp-127,
		fmt32((fullMantissaWithSign)*(math.Pow(2, float64(biasedExp-127)))),
	)

	for _, s := range steps {
		fmt.Println(s)
	}
}

func fmt32(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
