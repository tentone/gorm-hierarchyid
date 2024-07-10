package main

import (
	"encoding/hex"
	"slices"
	"testing"
)

func TestBinaryString(t *testing.T) {
	type TestData struct {
		input  string
		output string
	}

	var data []TestData = []TestData{
		{"48", "01001"},
		{"58", "01011"},
		{"68", "01101"},
		{"78", "01111"},
		{"84", "100001"},
		{"8C", "100011"},
		{"94", "100101"},
		{"9C", "100111"},
		{"A2", "1010001"},
		{"A6", "1010011"},
		{"AA", "1010101"},
		{"AE", "1010111"},
		{"B2", "1011001"},
		{"B6", "1011011"},
		{"BA", "1011101"},
		{"BE", "1011111"},
		{"C110", "110000010001"},
		{"C130", "110000010011"},
		{"C150", "110000010101"},
		{"C170", "110000010111"},
		{"C190", "110000011001"},
		{"C1B0", "110000011011"},
		{"C1D0", "110000011101"},
		{"C1F0", "110000011111"},
		{"C310", "110000110001"},
		{"C910", "110010010001"},
		{"CB10", "110010110001"},
		{"D110", "110100010001"},
		{"D310", "110100110001"},
		{"D910", "110110010001"},
		{"DB10", "110110110001"},
		{"E00440", "111000000000010001"},
		{"E00C40", "111000000000110001"},
		{"E02440", "111000000010010001"},
		{"E06440", "111000000110010001"},
		{"E06C40", "111000000110110001"},
		{"E0E440", "111000001110010001"},
		{"E2E440", "111000101110010001"},
		{"E6E440", "111001101110010001"},
		{"EEE440", "111011101110010001"},
		{"F00088", "111100000000000010001"},
		{"F20088", "111100100000000010001"},
		{"F40088", "111101000000000010001"},
		{"F60088", "111101100000000010001"},
		{"F80000000220", "1111100000000000000000000000000000000010001"},
	}

	for _, d := range data {
		input, err := hex.DecodeString(d.input)
		if err != nil {
			t.Errorf("Error decoding %v: %v", d.input, err)
		}

		var result = BinaryString(input)

		if result != d.output {
			t.Errorf("Expected 0x%v to return %v, got %v", d.input, d.output, result)
		}
	}

}

func TestTestPatterns(t *testing.T) {
	type TestData struct {
		input  string
		output string
	}

	var data []TestData = []TestData{
		{"01001", "01xxT"},
		{"01010", "01xxT"},
		{"100001", "100xxT"},
		{"001110101", "00111xxxT"},
	}

	for _, d := range data {
		result, err := TestPatterns(d.input)
		if err != nil {
			t.Errorf("Error testing %v: %v", d.input, err)
		}

		if result.Pattern != d.output {
			t.Errorf("Expected %v to return %v, got %v", d.input, d.output, result)
		}
	}
}

func TestParse(t *testing.T) {
	type TestData struct {
		output []int
		input  string
	}

	var data []TestData = []TestData{
		{nil, ""},
		{[]int{-73}, "1BEEFC"},
		{[]int{-72}, "2088"},
		{[]int{-64}, "2188"},
		{[]int{-56}, "2488"},
		{[]int{-48}, "2588"},
		{[]int{-40}, "2888"},
		{[]int{-32}, "2988"},
		{[]int{-24}, "2C88"},
		{[]int{-16}, "2D88"},
		{[]int{-10}, "2DE8"},
		{[]int{-9}, "2DF8"},
		{[]int{-8}, "3880"},
		{[]int{-7}, "3980"},
		{[]int{-6}, "3A80"},
		{[]int{-5}, "3B80"},
		{[]int{-4}, "3C80"},
		{[]int{-3}, "3D80"},
		{[]int{-2}, "3E80"},
		{[]int{-1}, "3F80"},
		{[]int{0}, "48"},
		{[]int{1}, "58"},
		{[]int{2}, "68"},
		{[]int{3}, "78"},
		{[]int{4}, "84"},
		{[]int{5}, "8C"},
		{[]int{6}, "94"},
		{[]int{7}, "9C"},
		{[]int{8}, "A2"},
		{[]int{9}, "A6"},
		{[]int{10}, "AA"},
		{[]int{11}, "AE"},
		{[]int{12}, "B2"},
		{[]int{13}, "B6"},
		{[]int{14}, "BA"},
		{[]int{15}, "BE"},
		{[]int{16}, "C110"},
		{[]int{17}, "C130"},
		{[]int{18}, "C150"},
		{[]int{19}, "C170"},
		{[]int{20}, "C190"},
		{[]int{21}, "C1B0"},
		{[]int{22}, "C1D0"},
		{[]int{23}, "C1F0"},
		{[]int{24}, "C310"},
		{[]int{32}, "C910"},
		{[]int{40}, "CB10"},
		{[]int{48}, "D110"},
		{[]int{56}, "D310"},
		{[]int{64}, "D910"},
		{[]int{72}, "DB10"},
		{[]int{80}, "E00440"},
		{[]int{88}, "E00C40"},
		{[]int{96}, "E02440"},
		{[]int{128}, "E06440"},
		{[]int{136}, "E06C40"},
		{[]int{192}, "E0E440"},
		{[]int{320}, "E2E440"},
		{[]int{576}, "E6E440"},
		{[]int{1088}, "EEE440"},
		{[]int{1104}, "F00088"},
		{[]int{2128}, "F20088"},
		{[]int{3152}, "F40088"},
		{[]int{4176}, "F60088"},
		{[]int{5200}, "F80000000220"},
		{[]int{3, 1}, "7AC0"},
		{[]int{1, 1}, "5AC0"},
		{[]int{2, 1}, "6AC0"},
		{[]int{2, 1, 1}, "6AD6"},
		{[]int{1, 1, 2}, "5ADA"},
		{[]int{1, 1, 3}, "5ADE"},
		{[]int{1, 1, 1}, "5AD6"},
		// {[]int{1, 1, 4}, "5AE1"},
		// {[]int{1, 2, 754}, "5B7A9150"},
		// {[]int{1, 1, 1, 1}, "5AD6B0"},
		// {[]int{2, 1, 1, 3}, "6AD6F0"},
		// {[]int{2, 1, 1, 1}, "6AD6B0"},
		// {[]int{2, 1, 1, 2}, "6AD6D0"},
		// {[]int{1, 1, 1, 1, 1}, "5AD6B580"},
	}

	for _, d := range data {
		input, err := hex.DecodeString(d.input)
		if err != nil {
			t.Errorf("Error decoding %v: %v", d.input, err)
		}

		result, err := Parse(input)
		if err != nil {
			t.Errorf("Error parsing %v: %v", d.input, err)
		}

		if !slices.Equal(result, d.output) {
			t.Errorf("Expected 0x%v to return %v, got %v", d.input, d.output, result)
		}
	}
}
