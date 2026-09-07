// Copyright (c) 2026 slavkiy

use crate::ParseError;
use ast::Program;
use lexer::Lexer;
use token::{Token, TokenKind};

pub struct Parser<'a> {
    lexer: Lexer<'a>,
    current: Token,
    lookahead: Option<Token>,
}

impl<'a> Parser<'a> {
    pub fn new(mut lexer: Lexer<'a>) -> Self {
        let current = lexer.next_token();

        Self {
            lexer,
            current,
            lookahead: None,
        }
    }

    pub fn parse() -> Result<Program, ParseError> {
        Ok(Program::default())
    }

    pub fn peek(&self) -> &Token {
        &self.current
    }

    pub fn peek_kind(&self) -> &TokenKind {
        &self.current.kind
    }

    pub fn peek_next(&mut self) -> &Token {
        if self.lookahead.is_none() {
            self.lookahead = Some(self.lexer.next_token());
        }

        self.lookahead.as_ref().unwrap()
    }

    pub fn parse_expression(&mut self) -> Result<ast::Expr, ParseError> {
        self.parse_assignment()
    }

    pub fn at(&self, kind: &TokenKind) -> bool {
        self.current.kind == *kind
    }

    pub fn bump(&mut self) -> Token {
        let next = self.next_token();

        std::mem::replace(&mut self.current, next)
    }

    fn next_token(&mut self) -> Token {
        loop {
            let token = if let Some(token) = self.lookahead.take() {
                token
            } else {
                self.lexer.next_token()
            };

            match token.kind {
                TokenKind::Whitespace { .. } | TokenKind::Comment { .. } => continue,
                _ => return token,
            }
        }
    }

    pub fn eat(&mut self, kind: &TokenKind) -> Option<Token> {
        if self.at(kind) {
            Some(self.bump())
        } else {
            None
        }
    }

    pub fn expect(&mut self, kind: &TokenKind) -> Result<Token, ParseError> {
        if self.at(kind) {
            Ok(self.bump())
        } else {
            Err(self.error_expected(kind))
        }
    }

    pub fn error(&self, message: impl Into<String>) -> ParseError {
        ParseError::new(self.current.span, message)
    }

    fn error_expected(&self, expected: &TokenKind) -> ParseError {
        self.error(format!(
            "expected {:?}, found {:?}",
            expected, self.current.kind,
        ))
    }
}
