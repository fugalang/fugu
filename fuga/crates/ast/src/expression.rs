// Copyright (c) 2026 slavkiy

use crate::{
    declaration::Parameter,
    literal::Literal,
    modifiers::Modifiers,
    operator::{AssignOp, BinaryOp, PostfixOp, UnaryOp},
    path::Path,
    pattern::Pattern,
    statement::{Block, Stmt},
    types::{Generics, Type},
};

#[derive(Debug, Clone, PartialEq)]
pub enum Expr {
    Literal(Literal),
    Path(Path),
    Unary {
        op: UnaryOp,
        expr: Box<Expr>,
    },
    Binary {
        left: Box<Expr>,
        op: BinaryOp,
        right: Box<Expr>,
    },
    Assign {
        left: Box<Expr>,
        op: AssignOp,
        right: Box<Expr>,
    },
    Postfix {
        expr: Box<Expr>,
        op: PostfixOp,
    },
    Statement {
        statement: Vec<Box<Stmt>>,
        expression: Box<Expr>,
    },
    Call(Call),
    Lambda(Lambda),
    Match(Match),
}

#[derive(Debug, Clone, PartialEq)]
pub struct Call {
    pub callee: Box<Expr>,
    pub args: Vec<Expr>,
    pub generics: Vec<Type>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Lambda {
    pub generics: Generics,
    pub params: Vec<Parameter>,
    pub return_types: Vec<Type>,
    pub body: Block,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Match {
    pub expression: Box<Expr>,
    pub arms: Vec<MatchArm>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct MatchArm {
    pub pattern: Pattern,
    pub body: Block,
}
