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
	s = TrimHexPrefix(s)

	return len(s) == AddressLength && isHex(s)
}

// IsHex validates whether each byte is valid hexadecimal string.
func IsHexString(s string) bool {
	return isHex(TrimHexPrefix(s))
}

func isHex(s string) bool {
	for _, c := range []byte(s) {
		if !isHexCharacter(c) {
			return false
		}
	}

	return true
}

// TrimHexPrefix trims prefix like '0x' or '0X'.
func TrimHexPrefix(s string) string {
	if hasHexPrefix(s) {
		return s[2:]
	}

	return s
}

// hasHexPrefix validates str begins with '0x' or '0X'.
func hasHexPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isHexCharacter returns bool of c being a valid hexadecimal.
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
