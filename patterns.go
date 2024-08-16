package hierarchyid

// Represents a possible pattern for hierarchyid values.
//
// The structure used to store the values, changes based on the value size. These patterns codify the possible structures based on value size.
type HierarchyIdPattern struct {
	// Minimum value
	Min int64

	// Maximum value
	Max int64

	// Pattern structure (0 and 1 are fixed}, x corresponds to the value and T is the terminator.
	Pattern string
}

// List of possible patterns for hierarchyid values.
//
// Sorted from the largest to the smallest.
var Patterns = []HierarchyIdPattern{
	{-281479271682120, -4294971465, "000100xxxxxxxxxxxxxx0xxxxxxxxxxxxxxxxxxxxx0xxxxxx0xxx0x1xxxT"},
	{4294972496, 281479271683151, "111111xxxxxxxxxxxxxx0xxxxxxxxxxxxxxxxxxxxx0xxxxxx0xxx0x1xxxT"},
	{-4294971464, -4169, "000101xxxxxxxxxxxxxxxxxxx0xxxxxx0xxx0x1xxxT"},
	{5200, 4294972495, "111110xxxxxxxxxxxxxxxxxxx0xxxxxx0xxx0x1xxxT"},
	{-4168, -73, "000110xxxxx0xxx0x1xxxT"},
	{1104, 5199, "11110xxxxx0xxx0x1xxxT"},
	{80, 1103, "1110xxx0xxx0x1xxxT"},
	{-72, -9, "0010xx0x1xxxT"},
	{16, 79, "110xx0x1xxxT"},
	{-8, -1, "00111xxxT"},
	{8, 15, "101xxxT"},
	{4, 7, "100xxT"},
	{0, 3, "01xxT"},
}
