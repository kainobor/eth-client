package helper

import (
	"fmt"
	"math/big"
	"strings"
)

// BigToHex converts big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}

	return "0x" + strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0"))
}

// HexToBig converts hex string to *big.Int
func HexToBig(hex string) (*big.Int, bool) {
	hex = strings.TrimPrefix(hex, "0x")
	bigVal := new(big.Int)

	return bigVal.SetString(hex, 16)
}
