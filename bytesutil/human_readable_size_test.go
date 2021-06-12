package bytesutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormats(t *testing.T) {
	tests := map[int64][2]string{
		-1:                        {"", ""},
		0:                         {"0 B", "0 B"},
		1_023:                     {"1023 B", "1.02 KB"},
		2_140:                     {"2.09 KiB", "2.14 KB"},
		2_828_382:                 {"2.70 MiB", "2.83 MB"},
		2_341_234_123_412_341_234: {"2.03 EiB", "2.34 EB"},
	}
	for value, expectedValues := range tests {
		assert.Equal(t, expectedValues[0], BinaryFormat(value))
		assert.Equal(t, expectedValues[1], DecimalFormat(value))
	}
}
