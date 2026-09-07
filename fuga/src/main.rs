// Copyright (c) 2026 slavkiy

fn main() {
    let mut lexer = lexer::Lexer::new("let x := 42");

    loop {
        let token = lexer.next_token();
        println!("{token:?}");

        if token.kind == token::TokenKind::Eof {
            break;
        }
    }
}
