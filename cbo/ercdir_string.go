// Code generated by "stringer -type=ERCDir"; DO NOT EDIT.

package cbo

import "fmt"

const (
	_ERCDir_name_0 = "Undirected"
	_ERCDir_name_1 = "Bidirectional"
	_ERCDir_name_2 = "Input"
	_ERCDir_name_3 = "Output"
	_ERCDir_name_4 = "MultiOutputSinkFlagNoConnectFlag"
)

var (
	_ERCDir_index_0 = [...]uint8{0, 10}
	_ERCDir_index_1 = [...]uint8{0, 13}
	_ERCDir_index_2 = [...]uint8{0, 5}
	_ERCDir_index_3 = [...]uint8{0, 6}
	_ERCDir_index_4 = [...]uint8{0, 19, 32}
)

func (i ERCDir) String() string {
	switch {
	case i == 0:
		return _ERCDir_name_0
	case i == 66:
		return _ERCDir_name_1
	case i == 73:
		return _ERCDir_name_2
	case i == 79:
		return _ERCDir_name_3
	case 9410 <= i && i <= 9411:
		i -= 9410
		return _ERCDir_name_4[_ERCDir_index_4[i]:_ERCDir_index_4[i+1]]
	default:
		return fmt.Sprintf("ERCDir(%d)", i)
	}
}
