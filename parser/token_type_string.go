// Code generated by "stringer -type TokenType -output token_type_string.go"; DO NOT EDIT.

package parser

import "fmt"

const _TokenType_name = "TokenNilTokenWhitespaceTokenBangTokenPercentTokenBitwiseAndTokenOParenTokenCParenTokenStarTokenPlusTokenCommaTokenMinusTokenDotTokenSlashTokenColonTokenSemicolonTokenLessThanTokenAssignTokenGreaterThanTokenQuestionTokenCommentTokenIdentTokenNumberLitTokenStringLitTokenOBrackTokenCBrackTokenCaretTokenOBraceTokenBitwiseOrTokenCBraceTokenBitwiseNotTokenDashDashTokenDotDotTokenAndTokenOrTokenEqualTokenNotEqualTokenLessThanEqTokenGreaterThanEqTokenEOFTokenInvalidTokenBadUTF8"

var _TokenType_map = map[TokenType]string{
	0:      _TokenType_name[0:8],
	32:     _TokenType_name[8:23],
	33:     _TokenType_name[23:32],
	37:     _TokenType_name[32:44],
	38:     _TokenType_name[44:59],
	40:     _TokenType_name[59:70],
	41:     _TokenType_name[70:81],
	42:     _TokenType_name[81:90],
	43:     _TokenType_name[90:99],
	44:     _TokenType_name[99:109],
	45:     _TokenType_name[109:119],
	46:     _TokenType_name[119:127],
	47:     _TokenType_name[127:137],
	58:     _TokenType_name[137:147],
	59:     _TokenType_name[147:161],
	60:     _TokenType_name[161:174],
	61:     _TokenType_name[174:185],
	62:     _TokenType_name[185:201],
	63:     _TokenType_name[201:214],
	67:     _TokenType_name[214:226],
	73:     _TokenType_name[226:236],
	78:     _TokenType_name[236:250],
	83:     _TokenType_name[250:264],
	91:     _TokenType_name[264:275],
	93:     _TokenType_name[275:286],
	94:     _TokenType_name[286:296],
	123:    _TokenType_name[296:307],
	124:    _TokenType_name[307:321],
	125:    _TokenType_name[321:332],
	126:    _TokenType_name[332:347],
	8212:   _TokenType_name[347:360],
	8230:   _TokenType_name[360:371],
	8743:   _TokenType_name[371:379],
	8744:   _TokenType_name[379:386],
	8788:   _TokenType_name[386:396],
	8800:   _TokenType_name[396:409],
	8804:   _TokenType_name[409:424],
	8805:   _TokenType_name[424:442],
	9220:   _TokenType_name[442:450],
	65533:  _TokenType_name[450:462],
	128169: _TokenType_name[462:474],
}

func (i TokenType) String() string {
	if str, ok := _TokenType_map[i]; ok {
		return str
	}
	return fmt.Sprintf("TokenType(%d)", i)
}
