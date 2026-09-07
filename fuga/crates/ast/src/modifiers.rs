// Copyright (c) 2026 slavkiy

use crate::declaration::Decl;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum Visibility {
    #[default]
    Private,
    Public,
}

#[derive(Debug, Clone, PartialEq, Default)]
pub struct Modifiers {
    pub visibility: Visibility,
    pub mutable: bool,
    pub is_macro: bool,
    pub is_comptim: bool,
    pub directives: Directives,
}

// is_unsafe
#[derive(Debug, Clone, PartialEq, Default)]
pub struct Directives {
    pub values: Vec<Box<Decl>>, // Directive
}
