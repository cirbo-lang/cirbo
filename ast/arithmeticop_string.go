// Code generated by "stringer -type=ArithmeticOp"; DO NOT EDIT.

package ast

import "fmt"

const _ArithmeticOp_name = "ArithmeticOpNilAddSubtractLessThanEqualGreaterThanExponentModuloNotNegateMultiplyDivideConcatAndOrNotEqualLessThanOrEqualGreaterThanOrEqual"

var _ArithmeticOp_map = map[ArithmeticOp]string{
	0:    _ArithmeticOp_name[0:15],
	43:   _ArithmeticOp_name[15:18],
	45:   _ArithmeticOp_name[18:26],
	60:   _ArithmeticOp_name[26:34],
	61:   _ArithmeticOp_name[34:39],
	62:   _ArithmeticOp_name[39:50],
	94:   _ArithmeticOp_name[50:58],
	109:  _ArithmeticOp_name[58:64],
	172:  _ArithmeticOp_name[64:67],
	177:  _ArithmeticOp_name[67:73],
	215:  _ArithmeticOp_name[73:81],
	247:  _ArithmeticOp_name[81:87],
	8230: _ArithmeticOp_name[87:93],
	8743: _ArithmeticOp_name[93:96],
	8744: _ArithmeticOp_name[96:98],
	8800: _ArithmeticOp_name[98:106],
	8804: _ArithmeticOp_name[106:121],
	8805: _ArithmeticOp_name[121:139],
}

func (i ArithmeticOp) String() string {
	if str, ok := _ArithmeticOp_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ArithmeticOp(%d)", i)
}
