package main

import (
	"encoding/hex"
	"slices"
	"strings"
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

		var result = binaryString(input)

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
		{"01011", "01xxT"},
		{"100001", "100xxT"},
		{"001110101", "00111xxxT"},
	}

	for _, d := range data {
		result, err := testPatterns(d.input)
		if err != nil {
			t.Errorf("Error testing %v: %v", d.input, err)
		}

		if result.Pattern != d.output {
			t.Errorf("Expected %v to return %v, got %v", d.input, d.output, result)
		}
	}
}

type TestData struct {
	output []int64
	input  string
}

var TestEncodeDecodeData []TestData = []TestData{
	{nil, ""},
	{[]int64{-73}, "1BEEFC"},
	{[]int64{-72}, "2088"},
	{[]int64{-64}, "2188"},
	{[]int64{-56}, "2488"},
	{[]int64{-48}, "2588"},
	{[]int64{-40}, "2888"},
	{[]int64{-32}, "2988"},
	{[]int64{-24}, "2C88"},
	{[]int64{-16}, "2D88"},
	{[]int64{-10}, "2DE8"},
	{[]int64{-9}, "2DF8"},
	{[]int64{-8}, "3880"},
	{[]int64{-7}, "3980"},
	{[]int64{-6}, "3A80"},
	{[]int64{-5}, "3B80"},
	{[]int64{-4}, "3C80"},
	{[]int64{-3}, "3D80"},
	{[]int64{-2}, "3E80"},
	{[]int64{-1}, "3F80"},
	{[]int64{0}, "48"},
	{[]int64{1}, "58"},
	{[]int64{2}, "68"},
	{[]int64{3}, "78"},
	{[]int64{4}, "84"},
	{[]int64{5}, "8C"},
	{[]int64{6}, "94"},
	{[]int64{7}, "9C"},
	{[]int64{8}, "A2"},
	{[]int64{9}, "A6"},
	{[]int64{10}, "AA"},
	{[]int64{11}, "AE"},
	{[]int64{12}, "B2"},
	{[]int64{13}, "B6"},
	{[]int64{14}, "BA"},
	{[]int64{15}, "BE"},
	{[]int64{16}, "C110"},
	{[]int64{17}, "C130"},
	{[]int64{18}, "C150"},
	{[]int64{19}, "C170"},
	{[]int64{20}, "C190"},
	{[]int64{21}, "C1B0"},
	{[]int64{22}, "C1D0"},
	{[]int64{23}, "C1F0"},
	{[]int64{24}, "C310"},
	{[]int64{32}, "C910"},
	{[]int64{40}, "CB10"},
	{[]int64{48}, "D110"},
	{[]int64{56}, "D310"},
	{[]int64{64}, "D910"},
	{[]int64{72}, "DB10"},
	{[]int64{80}, "E00440"},
	{[]int64{88}, "E00C40"},
	{[]int64{96}, "E02440"},
	{[]int64{128}, "E06440"},
	{[]int64{136}, "E06C40"},
	{[]int64{192}, "E0E440"},
	{[]int64{320}, "E2E440"},
	{[]int64{576}, "E6E440"},
	{[]int64{1088}, "EEE440"},
	{[]int64{1104}, "F00088"},
	{[]int64{2128}, "F20088"},
	{[]int64{3152}, "F40088"},
	{[]int64{4176}, "F60088"},
	{[]int64{5200}, "F80000000220"},
	{[]int64{3, 1}, "7AC0"},
	{[]int64{1, 1}, "5AC0"},
	{[]int64{2, 1}, "6AC0"},
	{[]int64{2, 1, 1}, "6AD6"},
	{[]int64{1, 1, 2}, "5ADA"},
	{[]int64{1, 1, 3}, "5ADE"},
	{[]int64{1, 1, 1}, "5AD6"},
	{[]int64{1, 1, 4}, "5AE1"},
	{[]int64{1, -1, 4}, "59FE10"},
	{[]int64{1, 1, 1, 1, 1, 1, 2, 1, 1, 2}, "5AD6B5ADAB5B40"},
	{[]int64{1, 2, 754}, "5B7A9150"},
	{[]int64{1, 1, 1, 1}, "5AD6B0"},
	{[]int64{2, 1, 1, 3}, "6AD6F0"},
	{[]int64{2, 1, 1, 1}, "6AD6B0"},
	{[]int64{2, 1, 1, 2}, "6AD6D0"},
	{[]int64{1, 1, 1, 1, 1}, "5AD6B580"},
}

func TestDecode(t *testing.T) {
	for _, d := range TestEncodeDecodeData {
		input, err := hex.DecodeString(d.input)
		if err != nil {
			t.Errorf("Error decoding %v: %v", d.input, err)
		}

		result, err := Decode(input)
		if err != nil {
			t.Errorf("Error parsing %v: %v", d.input, err)
		}

		if !slices.Equal(result, d.output) {
			t.Errorf("Expected 0x%v to return %v, got %v", d.input, d.output, result)
		}
	}
}

func TestEncode(t *testing.T) {
	for _, d := range TestEncodeDecodeData {

		result, err := Encode(d.output)
		if err != nil {
			t.Errorf("Error parsing %v: %v", d.input, err)
		}

		encoded := hex.EncodeToString(result)

		if strings.ToUpper(encoded) != strings.ToUpper(d.input) {
			t.Errorf("Expected %v to return %v, got %v", d.output, d.input, encoded)
		}
	}
}
