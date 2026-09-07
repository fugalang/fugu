// Copyright (c) 2026 slavkiy

use crate::{declaration::Variable, expression::Expr, pattern::Pattern};

#[derive(Debug, Clone, PartialEq, Default)]
pub struct Block {
    pub statements: Vec<Stmt>,
}

#[derive(Debug, Clone, PartialEq)]
pub enum Stmt {
    Variable(Variable),
    Expression(Expr),

    Goto {
        tag: String,
    },

    If(IfStatement),
    Switch {
        expression: Expr,
        cases: Vec<SwitchCase>,
        default: Option<Block>,
    },
    For {
        condition: Option<Expr>,
        body: Block,
    },
    Continue,
    Break,
    Defer(DeferStatement),
}

#[derive(Debug, Clone, PartialEq)]
pub struct SwitchCase {
    pub pattern: Option<Pattern>,
    pub expression: Option<Expr>,
    pub body: Block,
}

#[derive(Debug, Clone, PartialEq)]
pub struct IfStatement {
    // if (condition) {}
    // if (let var := expr; var > 0) {}
    pub condition: Expr,
    pub then_branch: Block,
    pub else_branch: Option<Box<Stmt>>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct DeferStatement {
    pub statement: Option<Box<Stmt>>,
    pub expression: Option<Expr>,
}
