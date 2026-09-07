// Copyright (c) 2026 slavkiy

use crate::declaration::Decl;

#[derive(Debug, Clone, PartialEq, Default)]
pub struct Program {
    pub declarations: Vec<Decl>,
}
