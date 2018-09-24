package helper

import (
	"fmt"
	"math/big"
	"strings"
)

// AddressLength is length af valid ETH address string
const AddressLength = 40

// BigToHex converts big.Int to hexadecimal representation
func BigToHex(bigInt big.Int) string {
	if bigInt.BitLen() == 0 {
		return "0x0"
	}

	return "0x" + strings.ToLower(strings.TrimPrefix(fmt.Sprintf("%x", bigInt.Bytes()), "0"))
}

// HexToBig converts hex string to *big.Int
func HexToBig(hex string) (*big.Int, bool) {
	if hasHexPrefix(hex) {
		hex = hex[2:]
	}
	bigVal := new(big.Int)

	return bigVal.SetString(hex, 16)
}

// IsHexAddress validates that string is valid ETH address
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == AddressLength && IsHex(s)
}

// IsHex validates whether each byte is valid hexadecimal string.
func IsHex(str string) bool {
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// hasHexPrefix validates str begins with '0x' or '0X'.
func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
