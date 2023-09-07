package utils

import (
	"reflect"
	"testing"
)

func Test_Map(t *testing.T) {
	tests := []struct {
		name     string
		result   []int
		expected []int
	}{
		{
			name: "# 1",
			result: Map([]int{1, 2, 3, 4}, func(item int, index int) int {
				return item * 2
			}),
			expected: []int{2, 4, 6, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if r := tt.result; !reflect.DeepEqual(r, tt.expected) {
				t.Errorf("Map() = %v, expected %v", r, tt.expected)
			}
		})
	}
}
