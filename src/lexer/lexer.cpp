#include <cstdio>
#include "lexer.h"

void yyerror(const char* s) {
    fprintf(stderr, "Ошибка: %s в строке %d\n", s, yylineno);
}
