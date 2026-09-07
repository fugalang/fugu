// Copyright (c) 2026 slavkiy

use crate::{expression::Expr, literal::Literal, modifiers::Modifiers, path::Path};

#[derive(Debug, Clone, PartialEq)]
pub enum Pattern {
    // _
    Wildcard,
    Name(String),
    Literal(Literal),
    // Some(value) | Result::Ok(value)
    Constructor {
        path: Path,
        args: Vec<Pattern>,
    },
    Tuple(Vec<Pattern>), // (a, b, c)
    // { name: pattern, ... }
    Record {
        fields: Vec<PatternField>,
    },
    // start..end | start..=end
    Range {
        start: Option<Box<Expr>>,
        end: Option<Box<Expr>>,
        inclusive: bool,
    },
    Expr(Box<Expr>),
}

#[derive(Debug, Clone, PartialEq)]
pub struct PatternField {
    pub name: String,
    pub pattern: Pattern, // pattern::Expr
    pub modifiers: Modifiers,
}
