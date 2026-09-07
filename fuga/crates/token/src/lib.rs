// Copyright (c) 2026 slavkiy

#[derive(Debug, Clone, PartialEq)]
pub enum TokenKind {
    Eof,
    NewLine,                                      // \n
    Whitespace { chars: Vec<Whitespace> },        //
    Comment { literal: String, multiline: bool }, // // comment

    Identifier { literal: String }, // var

    Float { literal: String },     // 1.4
    Integer { literal: String },   // 123
    Complex { literal: String },   // 1.3i
    String { literal: String },    // "string"
    RawString { literal: String }, // `string`
    Char { literal: char },        // 'с'

    True,  // true
    False, // false

    Module,    // module
    Package,   // pkg
    Import,    // imp
    Pub,       // pub
    Fn,        // fn
    Return,    // return
    Struct,    // struct
    Interface, // interface
    Type,      // type
    Impl,      // impl
    Let,       // let
    Mut,       // mut
    Const,     // const
    If,        // if
    Else,      // else
    Break,     // break
    Switch,    // switch
    Case,      // case
    Enum,      // enum
    Match,     // match
    For,       // for
    Continue,  // continue
    Goto,      // goto
    Defer,     // defer
    Unsafe,    // unsafe

    Plus,    // +
    Minus,   // -
    Star,    // *
    Slash,   // /
    Percent, // %
    Caret,   // ^

    EqualEqual,   // ==
    BangEqual,    // !=
    Less,         // <
    LessEqual,    // <=
    Greater,      // >
    GreaterEqual, // >=

    ColonEqual,   // :=
    Equal,        // =
    PlusEqual,    // +=
    MinusEqual,   // -=
    StarEqual,    // *=
    SlashEqual,   // /=
    PercentEqual, // %=
    CaretEqual,   // ^=

    Ampersand, // &
    Pipe,      // |
    Tilde,     // ~

    PlusPlus,   // ++
    MinusMinus, // --

    LessLess,       // <<
    GreaterGreater, // >>

    MinusGreater, // ->
    EqualGreater, // =>
    DotDot,       // ..
    DotDotDot,    // ...

    HashBracket, // #[

    LeftParen,    // (
    RightParen,   // )
    LeftBrace,    // {
    RightBrace,   // }
    LeftBracket,  // [
    RightBracket, // ]

    Bang,     // !
    Question, // ?

    Comma,      // ,
    Dot,        // .
    Colon,      // :
    ColonColon, // ::
    Semicolon,  // ;

    Illegal { literal: String },
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Whitespace {
    Space,
    Tab,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub struct Position {
    pub offset: usize,
    pub column: usize,
    pub line: usize,
}

#[derive(Debug, Clone, Copy)]
pub struct Span {
    pub start: Position,
    pub end: Position,
}

#[derive(Debug, Clone)]
pub struct Token {
    pub kind: TokenKind,
    pub span: Span,
}

impl Token {
    pub fn new(kind: TokenKind, span: Span) -> Self {
        Self { kind, span }
    }
}
