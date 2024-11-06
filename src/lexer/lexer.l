%{
#include "tokens.h"
#include "lexer.h"
#include <string>

int get_keyword_token(const char* lexeme);
int get_type_token(const char* lexeme);
int get_operator_token(const char* lexeme);

extern "C" {
    int yywrap();
}

int yywrap() {
    return 1;
}

extern YYSTYPE yylval;
%}


/* Регулярные выражения для токенов */
KEYWORD     (berry|say|recipe|ingredient|if|else|for|in|range)
TYPE        (int|float|bool|string|void|list)
BOOLEAN     (true|false)
ID          [a-zA-Z_][a-zA-Z0-9_]*
STRING      \"(\\.|[^\"])*\"
NUMBER      [0-9]+(\.[0-9]+)?
OPERATOR    ([+\-*/=><]|==|!=|>=|<=)
ASSIGN      =
DELIM       (,|:|\(|\)|\[|\]|\{|\}|;)

%%

{KEYWORD}       { return get_keyword_token(yytext); }
{TYPE}          { return get_type_token(yytext); }
{BOOLEAN}       { yylval.boolean = (strcmp(yytext, "true") == 0); return BOOLEAN; }

{ID}            { yylval.str = strdup(yytext); return IDENTIFIER; }
{STRING}        { yylval.str = strdup(yytext); return STRING; }
{NUMBER}        { yylval.num = atof(yytext); return NUMBER; }
{OPERATOR}      { return get_operator_token(yytext); }
{ASSIGN}        { return ASSIGN; }
{DELIM}         { return yytext[0]; }

[ \t\n]+        ;   // Пропуск пробельных символов
"//".*          ;   // Однострочные комментарии

%%

int get_keyword_token(const char* lexeme) {
    if (strcmp(lexeme, "berry") == 0) return BERRY;
    if (strcmp(lexeme, "say") == 0) return SAY;
    if (strcmp(lexeme, "recipe") == 0) return RECIPE;
    if (strcmp(lexeme, "ingredient") == 0) return INGREDIENT;
    if (strcmp(lexeme, "if") == 0) return IF;
    if (strcmp(lexeme, "else") == 0) return ELSE;
    if (strcmp(lexeme, "for") == 0) return FOR;
    if (strcmp(lexeme, "in") == 0) return IN;
    if (strcmp(lexeme, "range") == 0) return RANGE;
    return IDENTIFIER;
}

int get_type_token(const char* lexeme) {
    if (strcmp(lexeme, "int") == 0) return TYPE_INT;
    if (strcmp(lexeme, "float") == 0) return TYPE_FLOAT;
    if (strcmp(lexeme, "bool") == 0) return TYPE_BOOL;
    if (strcmp(lexeme, "string") == 0) return TYPE_STRING;
    if (strcmp(lexeme, "void") == 0) return TYPE_VOID;
    if (strcmp(lexeme, "list") == 0) return TYPE_LIST;
    return IDENTIFIER;
}

int get_operator_token(const char* lexeme) {
    if (strcmp(lexeme, "+") == 0) return OP_PLUS;
    if (strcmp(lexeme, "-") == 0) return OP_MINUS;
    if (strcmp(lexeme, "*") == 0) return OP_MULTIPLY;
    if (strcmp(lexeme, "/") == 0) return OP_DIVIDE;
    if (strcmp(lexeme, ">") == 0) return OP_GREATER;
    if (strcmp(lexeme, "<") == 0) return OP_LESS;
    if (strcmp(lexeme, "==") == 0) return OP_EQUAL;
    if (strcmp(lexeme, "!=") == 0) return OP_NOT_EQUAL;
    if (strcmp(lexeme, ">=") == 0) return OP_GREATER_EQUAL;
    if (strcmp(lexeme, "<=") == 0) return OP_LESS_EQUAL;
    return OP_UNKNOWN;
}
