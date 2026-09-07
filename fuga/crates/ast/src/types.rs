// Copyright (c) 2026 slavkiy

use crate::path::Path;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Type {
    Unit,
    Primitive(PrimitiveType),
    Named { path: Path, generics: Vec<Type> },
    Generic(String),
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum PrimitiveType {
    Unit,

    Bool,
    Char,
    I8,
    I16,
    I32,
    I64,
    I128,
    U8,
    U16,
    U32,
    U64,
    U128,
    F32,
    F64,
    C32,
    C64,
    Str,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct GenericParam {
    pub name: String,
    pub bound: Option<Type>,
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub struct Generics {
    pub params: Vec<GenericParam>,
}
