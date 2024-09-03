package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToHumanReadable(t *testing.T) {
	testCases := []struct {
		name string
		arg  int
		want string
	}{
		{
			name: "Should convert to MB",
			arg:  1024,
			want: "1.00 MB",
		},
		{
			name: "Should convert to GB",
			arg:  1048576,
			want: "1.00 GB",
		},
		{
			name: "Should convert to TB",
			arg:  1073741824,
			want: "1.00 TB",
		},
		{
			name: "should convert to PB",
			arg:  1024 * 1024 * 1024 * 1024,
			want: "1.00 PB",
		}, {
			name: "Should fail for negative values",
			arg:  -1024,
			want: "-1.00 KB",
		},
		{
			name: "sould convert to EB",
			arg:  1024 * 1024 * 1024 * 1024 * 1024,
			want: "1.00 EB",
		},
	}

	for _, tst := range testCases {
		t.Run(tst.name, func(t *testing.T) {
			got := KbToHumanReadable(uint(tst.arg))
			assert.Equal(t, tst.want, got)
		})
	}
}
