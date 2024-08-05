package main

import (
	"errors"
	"strconv"
	"strings"
)

// HierarchyIdData is a type to represent a hierarchyid data type from SQL Server
//
// The hierarchyid data type is a series of integers separated by slashes.  For example, \1\2\3\.
type HierarchyIdData = []int64

// Check if a hierarchyid is a descendant of another hierarchyid
func IsDescendantOf(child HierarchyIdData, parent HierarchyIdData) bool {
	// If the path to the parent if bigger return false
	if len(child) <= len(parent) {
		return false
	}

	// Check if all levels of the parent are the same as the descendent
	for i := 0; i < len(parent); i++ {
		if child[i] != parent[i] {
			return false
		}
	}

	return true
}

// Create a string representation of the hierarchyid data type
//
// The string representation is a series of integers separated by slashes.  For example, \1\2\3\
func ToString(data HierarchyIdData) string {
	var r string = "/"
	for _, level := range data {
		r += strconv.FormatInt(level, 10) + "/"
	}
	return r
}

// Get all ancestors (parents) of a hierarchyid.
func GetAncestors(data HierarchyIdData) []HierarchyIdData {
	var parents []HierarchyIdData = []HierarchyIdData{}

	for i := 0; i < len(data)-1; i++ {
		var parent = data[0 : i+1]
		parents = append(parents, parent)
	}

	return parents
}

// Get the direct ancestor of a hierarchyid.
func GetAncestor(data HierarchyIdData) HierarchyIdData {
	if len(data) == 0 {
		return []int64{}
	}

	return data[0 : len(data)-1]
}

// Create a hierarchyid data type from a string representation
func FromString(data string) (HierarchyIdData, error) {
	var levels []int64 = []int64{}
	if data == "" {
		return levels, nil
	}

	// Split the string into levels
	var parts = strings.Split(data, "/")
	for _, part := range parts {
		if part == "" {
			continue
		}

		var level, err = strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}

		levels = append(levels, level)
	}

	return levels, nil
}

// Compare two hierarchyid data types
//
// The comparison is done by comparing each level of the hierarchyid.  If the levels are the same, the next level is compared.  If the levels are different, the comparison stops and the result is returned.
func Compare(a HierarchyIdData, b HierarchyIdData) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}

	return 0
}

// Decode takes a byte slice of data stored in SQL Server hierarchyid format and returns a HierarchyId.
//
// SQL server uses a custom binary format for hierarchyid.
func Decode(data []byte) (HierarchyIdData, error) {
	var levels []int64 = []int64{}
	if len(data) == 0 {
		return levels, nil
	}

	// Convert binary data to a string of 0s and 1s
	var bin = binaryString(data)

	for {
		// Find pattern that fits  the binary data
		var pattern, err = testPatterns(bin)
		if err != nil {
			return nil, err
		}

		var value int64
		value, err = decodeValue(pattern.Pattern, bin)
		value += pattern.Min
		if err != nil {
			return nil, err
		}

		// Add value to the list of values
		levels = append(levels, value)

		// Remove already read data from binary string
		bin = bin[len(pattern.Pattern):]
		if bin == "" {
			break
		}
	}

	return levels, nil
}

// Encode a hierarchyid from hierarchyid.
func Encode(levels HierarchyIdData) ([]byte, error) {

	if len(levels) == 0 {
		return []byte{}, nil
	}

	var bin string = ""

	for _, level := range levels {

		// Find pattern that fits the binary data
		var pattern *HierarchyIdPattern = nil
		for i := 0; i < len(Patterns); i++ {
			if Patterns[i].Min <= level && Patterns[i].Max >= level {
				pattern = &Patterns[i]
				break
			}
		}

		if pattern == nil {
			return nil, errors.New("No pattern found for " + strconv.FormatInt(level, 10))
		}

		// Count the number of bits in the pattern
		var bitCount = 0
		for _, c := range pattern.Pattern {
			if c == 'x' {
				bitCount++
			}
		}

		// Convert value to binary
		var binLevel = strconv.FormatInt(level-pattern.Min, 2)
		for len(binLevel) < bitCount {
			binLevel = "0" + binLevel
		}

		// Convert binary to string
		var result = ""

		for i := 0; i < len(pattern.Pattern); i++ {
			var pChar = pattern.Pattern[i]
			if pChar == 'x' {
				if len(binLevel) > 0 {
					result += string(binLevel[0])
					binLevel = binLevel[1:]
				}
			} else if pChar == 'T' {
				result += "1"
			} else {
				result += string(pChar)
			}
		}

		bin += result
	}

	// Convert binary string to byte array
	var binLen = len(bin)
	var binBytes = make([]byte, (binLen+7)/8)
	for i := 0; i < binLen; i++ {
		if bin[i] == '1' {
			// Set the bit in the byte array
			var byteIndex = i / 8
			var bitIndex = i % 8
			binBytes[byteIndex] |= 1 << uint(7-bitIndex)
		}
	}

	return binBytes, nil
}

// Decode values a string representation of the hierarchyid data type for a pattern
func decodeValue(pattern string, bin string) (int64, error) {
	var binValue string = ""

	for i := 0; i < len(pattern); i++ {
		var pChar = pattern[i]
		if pChar == 'x' {
			binValue += string(bin[i])
		}
	}

	value, err := strconv.ParseInt(binValue, 2, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// Test pattern for binary data
//
// Return the pattern that fits the binary data (if any), the length of the pattern and an error.
func testPatterns(bin string) (*HierarchyIdPattern, error) {
	if len(bin) == 0 {
		return nil, errors.New("Binary string is empty")
	}

	if len(bin) < 5 {
		return nil, errors.New("Binary string " + bin + " is too short minimum length is 5")
	}

	// Check which pattern fits the start of the binary string (if any)
	for i := 0; i < len(Patterns); i++ {
		var pattern = Patterns[i].Pattern

		if len(pattern) > len(bin) {
			continue
		}

		// Match each character of the pattern with the binary string
		var patternMatch = false
		for j := 0; j < len(pattern); j++ {
			// Pattern is longer than the binary string
			if j >= len(bin) {
				break
			}

			// Get the pattern and binary characters
			var pChar = pattern[j]
			var bChar = bin[j]

			// If the pattern character is a terminator, stop the comparison (pattern has fully matched)
			if pChar == 'T' && bChar == '1' {
				patternMatch = true
				break
			}

			// If the pattern character is not a fixed value and the binary character is different, the pattern does not match
			if pChar != 'x' && pChar != bChar {
				patternMatch = false
				break
			}
		}

		if patternMatch {
			return &Patterns[i], nil
		}
	}

	return nil, errors.New("No pattern found for " + bin)
}

// Receives a byte array and prints as binary (0 and 1) data.
func binaryString(data []byte) string {
	var str = ""

	// Convert each byte to binary
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			if b&(1<<uint(i)) != 0 {
				str += "1"
			} else {
				str += "0"
			}
		}
	}

	// Remove all trailing zeros
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '0' {
			str = str[0:i]
		} else {
			break
		}
	}

	return str
}
