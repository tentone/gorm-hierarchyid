package main

import (
	"errors"
	"fmt"
	"strconv"
)

// HierarchyId is a type to represent a hierarchyid data type from SQL Server
//
// The hierarchyid data type is a series of integers separated by slashes.  For example, \1\2\3\.
type HierarchyId = []int64

// Parse takes a byte slice of data stored in SQL Server hierarchyid format and returns a HierarchyId.
//
// SQL server uses a custom binary format for hierarchyid.
func Parse(data []byte) (HierarchyId, error) {
	var levels []int64 = []int64{}
	if len(data) == 0 {
		return levels, nil
	}

	// Convert binary data to a string of 0s and 1s
	var bin = BinaryString(data)

	fmt.Println(" - Trying to parse data ", bin)

	for {
		// Find pattern that fits  the binary data
		var pattern, err = TestPatterns(bin)
		if err != nil {
			return nil, err
		}

		fmt.Println("    - Found pattern ", pattern)

		var value int64
		value, err = DecodeValue(pattern.Pattern, bin)
		value += pattern.Min
		if err != nil {
			return nil, err
		}

		fmt.Println("    - Decoded value ", value)

		// Add value to the list of values
		levels = append(levels, value)

		// Remove already read data from binary string
		bin = bin[0 : len(bin)-len(pattern.Pattern)]
		if bin == "" {
			break
		}

		fmt.Println("    - Remaining data to analyse ", bin)
	}

	// Reverse list of levels
	for i, j := 0, len(levels)-1; i < j; i, j = i+1, j-1 {
		levels[i], levels[j] = levels[j], levels[i]
	}

	return levels, nil
}

// DecodeValue a string representation of the hierarchyid data type for a pattern
func DecodeValue(pattern string, bin string) (int64, error) {
	var binValue string = ""

	for i := 0; i < len(pattern); i++ {
		var pChar = pattern[len(pattern)-i-1]

		if pChar == 'x' {
			binValue = string(bin[len(bin)-i-1]) + binValue
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
func TestPatterns(bin string) (*HierarchyIdPattern, error) {
	if len(bin) == 0 {
		return nil, errors.New("Binary string is empty")
	}

	if len(bin) < 5 {
		return nil, errors.New("Binary string is too short minimum length is 5")
	}

	// Check wich pattern fits the start of the binary string (if any)
	for i := 0; i < len(Patterns); i++ {
		var pattern = Patterns[i].Pattern

		// Match each character of the pattern with the binary string
		var patternMatch = false

		for j := 0; j < len(pattern); j++ {
			// Pattern is longer than the binary string
			var bIndex = j
			if bIndex >= len(bin) {
				break
			}

			// Get the pattern and binary characters
			var pChar = pattern[j]
			var bChar = bin[bIndex]

			// If the pattern character is a terminator, stop the comparison (pattern has fully matched)
			if pChar == 'T' {
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
func BinaryString(data []byte) string {
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

// Create a string representation of the hierarchyid data type
//
// The string representation is a series of integers separated by slashes.  For example, \1\2\3\
func ToString(data HierarchyId) string {
	var result string = "/"
	for _, level := range data {
		result += strconv.FormatInt(level, 10) + "/"
	}
	return result
}

// Compare two hierarchyid data types
//
// The comparison is done by comparing each level of the hierarchyid.  If the levels are the same, the next level is compared.  If the levels are different, the comparison stops and the result is returned.
func Compare(a HierarchyId, b HierarchyId) int {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}

	return 0
}
