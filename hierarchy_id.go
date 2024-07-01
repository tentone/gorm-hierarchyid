package main

import "strconv"

// HierarchyId is a type to represent a hierarchyid data type from SQL Server
//
// The hierarchyid data type is a series of integers separated by slashes.  For example, \1\2\3\.
type HierarchyId = []int

// Parse takes a byte slice of data stored in SQL Server hierarchyid format and returns a HierarchyId.
//
// SQL server uses a custom binary format for hierarchyid.
func Parse(data []byte) (HierarchyId, error) {
	var levels []int = make([]int, 0)
	if len(data) == 0 {
		return levels, nil
	}

	println(BinaryString(data))
	// TODO <ADD CODE HERE>

	return levels, nil
}

// Receives a byte array and prints as binary (0 and 1) data.
func BinaryString(data []byte) string {
	var str = ""

	for _, b := range data {
		for i := 7; i >= 0; i-- {
			if b&(1<<uint(i)) != 0 {
				str += "1"
			} else {
				str += "0"
			}
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
		result += strconv.Itoa(level) + "/"
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
