// Copyright (c) 2026 slavkiy

use ast::{
    Expr, Literal, Path, PathSegment,
    operator::{AssignOp, BinaryOp, PostfixOp, UnaryOp},
};
use token::TokenKind;

use crate::{ParseError, Parser};

impl<'a> Parser<'a> {
    pub(crate) fn parse_assignment(&mut self) -> Result<Expr, ParseError> {
        let left = self.parse_binary(1)?;

        let Some(op) = assignment_operator(self.peek_kind()) else {
            return Ok(left);
        };

        self.bump();
        let right = self.parse_assignment()?;

        Ok(Expr::Assign {
            left: Box::new(left),
            op,
            right: Box::new(right),
        })
    }

    fn parse_binary(&mut self, minimum_precedence: u8) -> Result<Expr, ParseError> {
        let mut left = self.parse_prefix()?;

        while let Some(op) = binary_operator(self.peek_kind()) {
            if op.precedence() < minimum_precedence {
                break;
            }

            self.bump();
            let right = self.parse_binary(op.precedence() + 1)?;
            left = Expr::Binary {
                left: Box::new(left),
                op,
                right: Box::new(right),
            };
        }

        Ok(left)
    }

    fn parse_prefix(&mut self) -> Result<Expr, ParseError> {
        let unary = match self.peek_kind() {
            TokenKind::Minus => Some(UnaryOp::Neg),
            TokenKind::Bang => Some(UnaryOp::Not),
            TokenKind::Tilde => Some(UnaryOp::BitNot),
            _ => None,
        };

        if let Some(op) = unary {
            self.bump();
            return Ok(Expr::Unary {
                op,
                expr: Box::new(self.parse_prefix()?),
            });
        }

        let mut expression = match self.bump().kind {
            TokenKind::Integer { literal } => Expr::Literal(Literal::Integer(
                literal
                    .parse()
                    .map_err(|_| self.error("invalid integer literal"))?,
            )),
            TokenKind::Float { literal } => Expr::Literal(Literal::Float(
                literal
                    .parse()
                    .map_err(|_| self.error("invalid float literal"))?,
            )),
            TokenKind::String { literal } => Expr::Literal(Literal::String(literal)),
            TokenKind::RawString { literal } => Expr::Literal(Literal::RawString(literal)),
            TokenKind::Char { literal } => Expr::Literal(Literal::Char(literal)),
            TokenKind::True => Expr::Literal(Literal::Bool(true)),
            TokenKind::False => Expr::Literal(Literal::Bool(false)),
            TokenKind::Identifier { literal } => Expr::Path(Path {
                segments: vec![PathSegment {
                    separator: None,
                    name: literal,
                }],
            }),
            TokenKind::LeftParen => {
                let expression = self.parse_assignment()?;
                self.expect(&TokenKind::RightParen)?;
                expression
            }
            found => return Err(self.error(format!("expected expression, found {found:?}"))),
        };

        loop {
            let op = match self.peek_kind() {
                TokenKind::PlusPlus => PostfixOp::Increment,
                TokenKind::MinusMinus => PostfixOp::Decrement,
                _ => break,
            };
            self.bump();
            expression = Expr::Postfix {
                expr: Box::new(expression),
                op,
            };
        }

        Ok(expression)
    }
}

fn assignment_operator(kind: &TokenKind) -> Option<AssignOp> {
    Some(match kind {
        TokenKind::Equal => AssignOp::Assign,
        TokenKind::PlusEqual => AssignOp::Add,
        TokenKind::MinusEqual => AssignOp::Sub,
        TokenKind::StarEqual => AssignOp::Mul,
        TokenKind::SlashEqual => AssignOp::Div,
        TokenKind::PercentEqual => AssignOp::Rem,
        TokenKind::CaretEqual => AssignOp::BitXor,
        _ => return None,
    })
}

fn binary_operator(kind: &TokenKind) -> Option<BinaryOp> {
    Some(match kind {
        TokenKind::Plus => BinaryOp::Add,
        TokenKind::Minus => BinaryOp::Sub,
        TokenKind::Star => BinaryOp::Mul,
        TokenKind::Slash => BinaryOp::Div,
        TokenKind::Ampersand => BinaryOp::BitAnd,
        TokenKind::Pipe => BinaryOp::BitOr,
        TokenKind::Caret => BinaryOp::BitXor,
        TokenKind::LessLess => BinaryOp::Shl,
        TokenKind::GreaterGreater => BinaryOp::Shr,
        TokenKind::EqualEqual => BinaryOp::Eq,
        TokenKind::BangEqual => BinaryOp::Ne,
        TokenKind::Less => BinaryOp::Lt,
        TokenKind::LessEqual => BinaryOp::Le,
        TokenKind::Greater => BinaryOp::Gt,
        TokenKind::GreaterEqual => BinaryOp::Ge,
        _ => return None,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use ast::Expr;
    use lexer::Lexer;

    #[test]
    fn parses_assignment_and_shift_precedence() {
        let mut parser = Parser::new(Lexer::new("value += amount << 2"));
        let expression = parser.parse_expression().unwrap();

        assert!(matches!(
            expression,
            Expr::Assign {
                op: AssignOp::Add,
                right,
                ..
            } if matches!(*right, Expr::Binary { op: BinaryOp::Shl, .. })
        ));
    }

    #[test]
    fn parses_unary_and_bitwise_operators() {
        let mut parser = Parser::new(Lexer::new("~value | mask"));
        let expression = parser.parse_expression().unwrap();

        assert!(matches!(
            expression,
            Expr::Binary {
                op: BinaryOp::BitOr,
                left,
                ..
            } if matches!(*left, Expr::Unary { op: UnaryOp::BitNot, .. })
        ));
    }
}
