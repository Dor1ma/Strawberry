%{
#include <stdio.h>
#include <stdlib.h>
#include "tokens.h"
%}

%union {
    bool boolean;
    char* str;
    double num;
}

%token <boolean> BISON_BOOLEAN
%token <str> BISON_STRING BISON_IDENTIFIER
%token <num> BISON_NUMBER

%%

input: /* пустой */
    | input line
    ;

line: BISON_NUMBER { printf("Parsed a number: %f\n", $1); }
    ;

%%

int main() {
    yyparse();
    return 0;
}

int yyerror(const char *s) {
    fprintf(stderr, "%s\n", s);
    return 0;
}
