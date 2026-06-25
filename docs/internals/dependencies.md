# Fugu: l-package dependency graph


```mermaid
graph

Start --> cli
cli --> composer

parser --> composer
project --> composer
preproc --> composer

ast --> parser
lexer --> parser
token --> parser
types --> parser

token --> lexer

casher --> project
runner --> project

library --> casher

library --> preproc
```