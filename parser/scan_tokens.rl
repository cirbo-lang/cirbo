package parser

import (
	"github.com/cirbo-lang/cirbo/source"
)

// This file is generated from scan_tokens.rl. DO NOT EDIT.
%%{
  # (except you are actually in scan_tokens.rl here, so edit away!)

  machine cirbotok;
  write data;
}%%

func scanTokens(data []byte, filename string, start source.Pos, mode scanMode) []Token {
    f := &tokenAccum{
        Filename: filename,
        Bytes:    data,
        Pos:      start,
    }

    %%{
        include UnicodeDerived "unicode_derived.rl";

        UTF8Cont = 0x80 .. 0xBF;
        AnyUTF8 = (
            0x00..0x7F |
            0xC0..0xDF . UTF8Cont |
            0xE0..0xEF . UTF8Cont . UTF8Cont |
            0xF0..0xF7 . UTF8Cont . UTF8Cont . UTF8Cont
        );
        BrokenUTF8 = any - AnyUTF8;

        NumberLitContinue = (digit|'.'|('e'|'E') ('+'|'-')? digit);
        NumberLit = digit ("" | (NumberLitContinue - '.') | (NumberLitContinue* (NumberLitContinue - '.')));
        StringLit = '"' (AnyUTF8 - ('"' | '\\' | '\r' | '\n') | '\\' AnyUTF8)+ '"';
        Ident = (('+' | '-') digit+ 'V' digit*) | ('~'? ID_Start ('~'? ID_Continue)*) | ("`" (AnyUTF8 - "`")+ "`");

        # Symbols that just represent themselves are handled as a single rule.
        SelfToken = "{" | "}" | "[" | "]" | "(" | ")" | "." | "," | "*" | "/" | "+" | "-" | '%' | "=" | "<" | ">" | "!" | "?" | ":" | "&" | "|" | "^" | ";";

        OPointyPointy = "<<";
        CPointyPointy = ">>";

        Equal = "==";
        NotEqual = "!=";
        GreaterThanEqual = ">=";
        LessThanEqual = "<=";
        LogicalAnd = "&&";
        LogicalOr = "||";

        DashDash = "--";
        DotDot = "..";
        BarDashDash = "|--";
        DashDashBar = "--|";

        Newline = '\r' ? '\n';
        EndOfLine = Newline;

        Comment = (
            ("//" (any - EndOfLine)* EndOfLine) |
            ("/*" any* "*/")
        );

        Whitespace = (' ' | '\r' | '\n' | '\t')+;

        main := |*
            Whitespace       => { token(TokenWhitespace) };
            NumberLit        => { token(TokenNumberLit) };
            StringLit        => { token(TokenStringLit) };
            Ident            => { token(TokenIdent) };

            Comment          => { token(TokenComment) };
            Newline          => { token(TokenNewline) };

            Equal            => { token(TokenEqual); };
            NotEqual         => { token(TokenNotEqual); };
            GreaterThanEqual => { token(TokenGreaterThanEq); };
            LessThanEqual    => { token(TokenLessThanEq); };
            LogicalAnd       => { token(TokenAnd); };
            LogicalOr        => { token(TokenOr); };
            BarDashDash      => { token(TokenBarDashDash); };
            DashDashBar      => { token(TokenDashDashBar); };
            DashDash         => { token(TokenDashDash); };
            DotDot           => { token(TokenDotDot); };
            OPointyPointy    => { token(TokenOPoint); };
            CPointyPointy    => { token(TokenCPoint); };
            SelfToken        => { selfToken() };

            BrokenUTF8       => { token(TokenBadUTF8) };
            AnyUTF8          => { token(TokenInvalid) };
        *|;

    }%%

    // Ragel state
	p := 0  // "Pointer" into data
	pe := len(data) // End-of-data "pointer"
    ts := 0
    te := 0
    act := 0
    eof := pe
    cs := cirbotok_en_main

    %%{
        prepush {
            stack = append(stack, 0);
        }
        postpop {
            stack = stack[:len(stack)-1];
        }
    }%%

    // Make Go compiler happy
    _ = ts
    _ = te
    _ = act
    _ = eof

    token := func (ty TokenType) {
        f.emitToken(ty, ts, te)
    }
    selfToken := func () {
        b := data[ts:te]
        if len(b) != 1 {
            // should never happen
            panic("selfToken only works for single-character tokens")
        }
        f.emitToken(TokenType(b[0]), ts, te)
    }

    %%{
        write init nocs;
        write exec;
    }%%

    // If we fall out here without being in a final state then we've
    // encountered something that the scanner can't match, which we'll
    // deal with as an invalid.
    if cs < cirbotok_first_final {
        f.emitToken(TokenInvalid, p, len(data))
    }

    // We always emit a synthetic EOF token at the end, since it gives the
    // parser position information for an "unexpected EOF" diagnostic.
    f.emitToken(TokenEOF, len(data), len(data))

    return f.Tokens
}
