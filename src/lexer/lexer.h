#ifndef STRAWBERRY_LEXER_H
#define STRAWBERRY_LEXER_H

int yylex();
extern int yylineno;
void yyerror(const char* s);

typedef union {
    int num;
    bool boolean;
    char* str;
} YYSTYPE;

extern YYSTYPE yylval;

#endif //STRAWBERRY_LEXER_H
