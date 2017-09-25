// Code generated by "stringer -type=ERCOutputType"; DO NOT EDIT.

package cbo

import "fmt"

const (
	_ERCOutputType_name_0 = "NoOutput"
	_ERCOutputType_name_1 = "UnknownOutput"
	_ERCOutputType_name_2 = "OpenCollector"
	_ERCOutputType_name_3 = "OpenEmitter"
	_ERCOutputType_name_4 = "PushPull"
	_ERCOutputType_name_5 = "Tristate"
)

var (
	_ERCOutputType_index_0 = [...]uint8{0, 8}
	_ERCOutputType_index_1 = [...]uint8{0, 13}
	_ERCOutputType_index_2 = [...]uint8{0, 13}
	_ERCOutputType_index_3 = [...]uint8{0, 11}
	_ERCOutputType_index_4 = [...]uint8{0, 8}
	_ERCOutputType_index_5 = [...]uint8{0, 8}
)

func (i ERCOutputType) String() string {
	switch {
	case i == 0:
		return _ERCOutputType_name_0
	case i == 63:
		return _ERCOutputType_name_1
	case i == 67:
		return _ERCOutputType_name_2
	case i == 69:
		return _ERCOutputType_name_3
	case i == 80:
		return _ERCOutputType_name_4
	case i == 177:
		return _ERCOutputType_name_5
	default:
		return fmt.Sprintf("ERCOutputType(%d)", i)
	}
}