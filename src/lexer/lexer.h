#ifndef STRAWBERRY_LEXER_H
#define STRAWBERRY_LEXER_H

int yylex();
extern int yylineno;
void yyerror(const char* s);

#endif //STRAWBERRY_LEXER_H
