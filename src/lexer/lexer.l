%{
#include "parser.tab.h"
#include <string.h>

int get_keyword_token(const char* lexeme);
int get_type_token(const char* lexeme);
int get_operator_token(const char* lexeme);

#ifdef __cplusplus
extern "C" {
    int yywrap();
}
#endif
%}

/* Регулярные выражения для токенов */
KEYWORD     (berry|say|recipe|if|else|elif|for|in|range|return)
TYPE        (int|float|bool|string|void|list)
BOOLEAN     (true|false)
ID          [a-zA-Z_][a-zA-Z0-9_]*
STRING      \"(\\.|[^\"])*\"
NUMBER      [0-9]+(\.[0-9]+)?
OPERATOR    ([+\-*/%=><]|==|!=|>=|<=)
DELIM       [,:\(\)\[\]\{\};]
ARROW       ->

%%

"->"            { return ARROW; }
"+"             { return OP_PLUS; }
"-"             { return OP_MINUS; }
"*"             { return OP_MULTIPLY; }
"/"             { return OP_DIVIDE; }
"%"             { return OP_MODULO; }
">="            { return OP_GREATER_EQUAL; }
"<="            { return OP_LESS_EQUAL; }
"=="            { return OP_EQUAL; }
"!="            { return OP_NOT_EQUAL; }
">"             { return OP_GREATER; }
"<"             { return OP_LESS; }
"="             { return ASSIGN; }
"["             { return OPEN_BRACKET; }
"]"             { return CLOSE_BRACKET; }
"("             { return OPEN_PAREN; }
")"             { return CLOSE_PAREN; }
"{"             { return OPEN_BRACE; }
"}"             { return CLOSE_BRACE; }
","             { return COMMA; }
":"             { return COLON; }
"berry"          { return BERRY; }
"int"            { return TYPE_INT; }
"list"           { return TYPE_LIST; }
"if"             { return IF; }
"say"            { return SAY; }


{KEYWORD}       { return get_keyword_token(yytext); }
{TYPE}          { return get_type_token(yytext); }
{BOOLEAN}       { yylval.boolean = (strcmp(yytext, "true") == 0); return BISON_T_BOOLEAN; }
{ID}            { yylval.str = strdup(yytext); return BISON_T_IDENTIFIER; }
{STRING}        { yylval.str = strdup(yytext); return BISON_T_STRING; }
{NUMBER}        { yylval.num = atof(yytext); return BISON_T_NUMBER; }

[ \t\n]+        ; // Пропуск пробельных символов
"//".*          ; // Однострочные комментарии

%%

int get_keyword_token(const char* lexeme) {
    if (strcmp(lexeme, "berry") == 0) return BERRY;
    if (strcmp(lexeme, "say") == 0) return SAY;
    if (strcmp(lexeme, "recipe") == 0) return RECIPE;
    if (strcmp(lexeme, "if") == 0) return IF;
    if (strcmp(lexeme, "else") == 0) return ELSE;
    if (strcmp(lexeme, "elif") == 0) return ELIF;
    if (strcmp(lexeme, "for") == 0) return FOR;
    if (strcmp(lexeme, "in") == 0) return IN;
    if (strcmp(lexeme, "range") == 0) return RANGE;
    if (strcmp(lexeme, "return") == 0) return RETURN;
    return BISON_T_IDENTIFIER;
}

int get_type_token(const char* lexeme) {
    if (strcmp(lexeme, "int") == 0) return TYPE_INT;
    if (strcmp(lexeme, "float") == 0) return TYPE_FLOAT;
    if (strcmp(lexeme, "bool") == 0) return TYPE_BOOL;
    if (strcmp(lexeme, "string") == 0) return TYPE_STRING;
    if (strcmp(lexeme, "void") == 0) return TYPE_VOID;
    if (strcmp(lexeme, "list") == 0) return TYPE_LIST;
    return BISON_T_IDENTIFIER;
}

int get_operator_token(const char* lexeme) {
    if (strcmp(lexeme, "+") == 0) return OP_PLUS;
    if (strcmp(lexeme, "-") == 0) return OP_MINUS;
    if (strcmp(lexeme, "*") == 0) return OP_MULTIPLY;
    if (strcmp(lexeme, "/") == 0) return OP_DIVIDE;
    if (strcmp(lexeme, "%") == 0) return OP_MODULO;
    if (strcmp(lexeme, ">") == 0) return OP_GREATER;
    if (strcmp(lexeme, "<") == 0) return OP_LESS;
    if (strcmp(lexeme, "==") == 0) return OP_EQUAL;
    if (strcmp(lexeme, "!=") == 0) return OP_NOT_EQUAL;
    if (strcmp(lexeme, ">=") == 0) return OP_GREATER_EQUAL;
    if (strcmp(lexeme, "<=") == 0) return OP_LESS_EQUAL;
    if (strcmp(lexeme, "=") == 0) return ASSIGN;
    return OP_UNKNOWN;
}

int yywrap() {
    return 1;
}
