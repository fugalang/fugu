// Copyright (c) 2026 slavkiy

use crate::{
    expression::Expr,
    modifiers::Modifiers,
    path::Path,
    pattern::Pattern,
    statement::Block,
    types::{Generics, Type},
};

#[derive(Debug, Clone, PartialEq)]
pub enum Decl {
    Module { name: Path },
    Package { name: Path, modifiers: Modifiers },
    Import { path: Path, alias: Option<String> },
    Goto { tag: String },
    Variable(Variable),
    Type(TypeDeclaration),
    Struct(StructDeclaration),
    Interface(InterfaceDeclaration),
    Enum(EnumDeclaration),
    Impl(ImplDeclaration),
    Function(Function),
    Unsafe(Block),
    Directive(Directive),
}

// let | const | mut (Patern): Vec<Type> = Vec<Expr>
// let (a, b): (u8, u8) = (1, 1)
#[derive(Debug, Clone, PartialEq)]
pub struct Variable {
    pub pattern: Pattern,
    pub ty: Option<Vec<Type>>,
    pub value: Option<Vec<Expr>>,
    pub modifiers: Modifiers,
}

// type Name Types = expr defualt value
// type UserName str = "def"
#[derive(Debug, Clone, PartialEq)]
pub struct TypeDeclaration {
    pub name: String,
    pub ty: Type,
    pub defualt: Box<Expr>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Field {
    pub name: Option<String>,
    pub ty: Type,
    pub default: Option<Expr>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct StructDeclaration {
    pub name: String,
    pub generics: Generics,
    pub fields: Vec<Field>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct InterfaceDeclaration {
    pub name: String,
    pub generics: Generics,
    pub fields: Vec<Field>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct EnumDeclaration {
    pub name: String,
    pub generics: Generics,
    pub fields: Vec<Field>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct ImplDeclaration {
    pub target: Path,
    pub generics: Generics,
    pub fields: Vec<Field>,
    pub body: Vec<Decl>,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Function {
    pub name: String,
    pub generics: Generics,
    pub params: Vec<Parameter>,
    pub return_types: Vec<Type>,
    pub body: Block,
    pub modifiers: Modifiers,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Parameter {
    pub pattern: Pattern,
    pub ty: Option<Type>,
    pub default: Option<Expr>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct Directive {
    pub name: Path,
    pub args: Vec<DirectiveArg>,
}

#[derive(Debug, Clone, PartialEq)]
pub struct DirectiveArg {
    pub name: Option<Path>,
    pub value: Expr,
}
