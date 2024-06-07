package gormhierarchyid

// Parse takes a byte slice and returns a slice of integers representing the hierarchyid
//
// data.  The byte slice is expected to be in the format of a hierarchyid data type from
// SQL Server.  The format is a string of integers separated by slashes.  For example, \1\2\3\
func Parse(data []byte) ([]int, error) {
	var levels []int = make([]int, 0)
	if len(data) == 0 {
		return levels, nil
	}

	return levels, nil
}
