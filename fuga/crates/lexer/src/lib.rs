// Copyright (c) 2026 slavkiy

use token::{Position, Span, Token, TokenKind, Whitespace};

#[derive(Debug)]
pub struct Lexer<'a> {
    chars: std::str::CharIndices<'a>,
    pos: Position,
}

impl<'a> Lexer<'a> {
    pub fn new(input: &'a str) -> Self {
        Self {
            chars: input.char_indices(),
            pos: Position {
                column: 1,
                line: 1,
                offset: 0,
            },
        }
    }

    pub fn next_token(&mut self) -> Token {
        let start = self.pos;

        let Some(ch) = self.nextch() else {
            return Token::new(
                TokenKind::Eof,
                Span {
                    start: start,
                    end: self.pos,
                },
            );
        };

        let kind = match ch {
            '\n' => TokenKind::NewLine,

            ' ' | '\t' | '\r' => {
                let mut chars = vec![match ch {
                    ' ' => Whitespace::Space,
                    '\r' => Whitespace::Space,
                    '\t' => Whitespace::Tab,
                    _ => unreachable!(),
                }];

                loop {
                    match self.peek() {
                        Some(' ') => {
                            self.nextch();
                            chars.push(Whitespace::Space);
                        }

                        Some('\t') => {
                            self.nextch();
                            chars.push(Whitespace::Tab);
                        }

                        _ => break,
                    }
                }

                TokenKind::Whitespace { chars }
            }

            '/' => match self.peek() {
                Some('/') | Some('*') => self.lex_comment(),

                Some('=') => {
                    self.nextch();
                    TokenKind::SlashEqual
                }

                _ => TokenKind::Slash,
            },

            _ if is_word_char(ch) => self.lex_identifier(ch),
            _ if ch.is_ascii_digit() => self.lex_number(ch),

            '"' => self.lex_string(),
            '`' => self.lex_raw_string(),
            '\'' => self.lex_char(ch),

            '+' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::PlusEqual
                }
                Some('+') => {
                    self.nextch();
                    TokenKind::PlusPlus
                }
                _ => TokenKind::Plus,
            },

            '-' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::MinusEqual
                }
                Some('>') => {
                    self.nextch();
                    TokenKind::MinusGreater
                }
                Some('-') => {
                    self.nextch();
                    TokenKind::MinusMinus
                }
                _ => TokenKind::Minus,
            },
            '*' => self.match_next('=', TokenKind::StarEqual, TokenKind::Star),
            '%' => self.match_next('=', TokenKind::PercentEqual, TokenKind::Percent),
            '^' => self.match_next('=', TokenKind::CaretEqual, TokenKind::Caret),

            '=' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::EqualEqual
                }
                Some('>') => {
                    self.nextch();
                    TokenKind::EqualGreater
                }
                _ => TokenKind::Equal,
            },

            '!' => self.match_next('=', TokenKind::BangEqual, TokenKind::Bang),

            '<' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::LessEqual
                }
                Some('<') => {
                    self.nextch();
                    TokenKind::LessLess
                }
                _ => TokenKind::Less,
            },

            '>' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::GreaterEqual
                }
                Some('>') => {
                    self.nextch();
                    TokenKind::GreaterGreater
                }
                _ => TokenKind::Greater,
            },

            ':' => match self.peek() {
                Some('=') => {
                    self.nextch();
                    TokenKind::ColonEqual
                }
                Some(':') => {
                    self.nextch();
                    TokenKind::ColonColon
                }
                _ => TokenKind::Colon,
            },

            '&' => TokenKind::Ampersand,
            '|' => TokenKind::Pipe,
            '~' => TokenKind::Tilde,
            '.' => match self.peek() {
                Some('.') => {
                    self.nextch();
                    if self.peek() == Some('.') {
                        self.nextch();
                        TokenKind::DotDotDot
                    } else {
                        TokenKind::DotDot
                    }
                }
                _ => TokenKind::Dot,
            },

            '#' => match self.peek() {
                Some('[') => {
                    self.nextch();
                    TokenKind::HashBracket
                }
                _ => TokenKind::Illegal {
                    literal: "#".to_string(),
                },
            },
            '(' => TokenKind::LeftParen,
            ')' => TokenKind::RightParen,
            '{' => TokenKind::LeftBrace,
            '}' => TokenKind::RightBrace,
            '[' => TokenKind::LeftBracket,
            ']' => TokenKind::RightBracket,
            '?' => TokenKind::Question,
            ',' => TokenKind::Comma,
            ';' => TokenKind::Semicolon,
            _ => TokenKind::Illegal {
                literal: ch.to_string(),
            },
        };

        Token::new(
            kind,
            Span {
                start,
                end: self.pos,
            },
        )
    }

    fn lex_identifier(&mut self, first: char) -> TokenKind {
        let mut value = String::new();
        value.push(first);

        while let Some(ch) = self.peek() {
            if is_word_char(ch) {
                value.push(self.nextch().unwrap());
            } else {
                break;
            }
        }

        match value.as_str() {
            "module" => TokenKind::Module,
            "pkg" => TokenKind::Package,
            "imp" => TokenKind::Import,

            "pub" => TokenKind::Pub,
            "let" => TokenKind::Let,
            "mut" => TokenKind::Mut,
            "const" => TokenKind::Const,

            "fn" => TokenKind::Fn,
            "return" => TokenKind::Return,

            "struct" => TokenKind::Struct,
            "interface" => TokenKind::Interface,
            "type" => TokenKind::Type,

            "if" => TokenKind::If,
            "else" => TokenKind::Else,
            "break" => TokenKind::Break,

            "switch" => TokenKind::Switch,
            "case" => TokenKind::Case,

            "enum" => TokenKind::Enum,
            "match" => TokenKind::Match,

            "for" => TokenKind::For,
            "continue" => TokenKind::Continue,

            "goto" => TokenKind::Goto,

            "defer" => TokenKind::Defer,
            "unsafe" => TokenKind::Unsafe,

            "true" => TokenKind::True,
            "false" => TokenKind::False,

            _ => TokenKind::Identifier { literal: value },
        }
    }

    fn lex_number(&mut self, first: char) -> TokenKind {
        let mut value = String::new();
        value.push(first);

        // part
        while let Some(ch) = self.peek() {
            if ch.is_ascii_digit() {
                value.push(self.nextch().unwrap());
            } else {
                break;
            }
        }

        // float
        if self.peek() == Some('.') {
            let mut chars = self.chars.clone();
            chars.next();

            if chars.next().map(|(_, ch)| ch) != Some('.') {
                value.push(self.nextch().unwrap());

                while let Some(ch) = self.peek() {
                    if ch.is_ascii_digit() {
                        value.push(self.nextch().unwrap());
                    } else {
                        break;
                    }
                }

                if self.peek() == Some('i') {
                    value.push(self.nextch().unwrap());

                    return TokenKind::Complex { literal: value };
                }

                return TokenKind::Float { literal: value };
            }
        }

        // complex integer: 123i
        if self.peek() == Some('i') {
            value.push(self.nextch().unwrap());

            return TokenKind::Complex { literal: value };
        }

        TokenKind::Integer { literal: value }
    }

    fn lex_string(&mut self) -> TokenKind {
        let mut value = String::new();

        loop {
            match self.peek() {
                Some('\\') => {
                    self.nextch();

                    let ch = match self.nextch() {
                        Some('n') => '\n',
                        Some('t') => '\t',
                        Some('r') => '\r',
                        Some('\\') => '\\',
                        Some('\'') => '\'',
                        Some('"') => '"',
                        Some(ch) => panic!("Unknown escape sequence: \\{}", ch),
                        None => panic!("Unterminated string"),
                    };

                    value.push(ch);
                }

                Some('"') => {
                    self.nextch();
                    return TokenKind::String { literal: value };
                }

                Some('\n') => {
                    panic!("Unterminated string")
                }

                Some(ch) => {
                    value.push(ch);
                    self.nextch();
                }

                None => {
                    panic!("Unterminated string");
                }
            }
        }
    }

    fn lex_raw_string(&mut self) -> TokenKind {
        let mut value = String::new();

        loop {
            let r = self.peek();
            match r {
                Some('`') => {
                    self.nextch();
                    return TokenKind::RawString { literal: value };
                }

                Some(ch) => {
                    value.push(ch);
                    self.nextch();
                }

                None => {
                    // Todo error;
                    panic!("Unterminated raw string");
                }
            }
        }
    }

    fn lex_char(&mut self, first: char) -> TokenKind {
        if first != '\'' {
            panic!("Expected char literal");
        }

        let ch = match self.nextch() {
            Some('\\') => match self.nextch() {
                Some('n') => '\n',
                Some('t') => '\t',
                Some('r') => '\r',
                Some('\\') => '\\',
                Some('\'') => '\'',
                Some('"') => '"',
                Some(ch) => panic!("Unknown escape sequence: \\{}", ch),
                None => panic!("Unterminated char literal"),
            },

            Some('\'') => {
                // ''
                panic!("Empty char literal");
            }

            Some(ch) => ch,

            None => {
                panic!("Unterminated char literal");
            }
        };

        match self.nextch() {
            Some('\'') => TokenKind::Char { literal: ch },

            Some(ch) => panic!("Expected closing ', found '{}'", ch),

            None => panic!("Unterminated char literal"),
        }
    }

    fn lex_comment(&mut self) -> TokenKind {
        match self.peek() {
            // // comment
            Some('/') => {
                self.nextch();

                let mut literal = String::new();

                while let Some(ch) = self.peek() {
                    if ch == '\n' {
                        break;
                    }

                    literal.push(self.nextch().unwrap());
                }

                TokenKind::Comment {
                    literal,
                    multiline: false,
                }
            }

            // /* comment */
            Some('*') => {
                self.nextch();

                let mut literal = String::new();

                loop {
                    let Some(ch) = self.nextch() else {
                        panic!("unterminated comment")
                    };

                    if ch == '*' && self.peek() == Some('/') {
                        self.nextch();
                        break;
                    }

                    literal.push(ch);
                }

                TokenKind::Comment {
                    literal,
                    multiline: true,
                }
            }

            _ => TokenKind::Slash,
        }
    }

    fn peek(&self) -> Option<char> {
        self.chars.clone().next().map(|(_, ch)| ch)
    }

    fn nextch(&mut self) -> Option<char> {
        let (_, ch) = self.chars.next()?;

        self.pos.offset += ch.len_utf8();

        if ch == '\n' {
            self.pos.line += 1;
            self.pos.column = 1;
        } else {
            self.pos.column += 1;
        }

        Some(ch)
    }

    fn match_next(
        &mut self,
        expected: char,
        matched: TokenKind,
        unmatched: TokenKind,
    ) -> TokenKind {
        if self.peek() == Some(expected) {
            self.nextch();
            matched
        } else {
            unmatched
        }
    }
}

fn is_word_char(ch: char) -> bool {
    ch.is_ascii_alphanumeric() || ch == '_'
}
