%{
#include <math.h>
#include <unistd.h>
#include <string>
#include <iostream>
#include "parser.h"

#include "../ast/Number/ast.h"

extern "C" int yywrap() { return 1; }
extern "C" void yyerror (const char *error) { std::cerr << error << std::endl; }

#define SAVE_TOKEN yylval.string = new std::string(yytext, yyleng)
#define TOKEN(t) (yylval.token = t)

%}

DIGIT    [0-9]
ID       [a-z][a-z0-9]*

%%

{DIGIT}+ {
  SAVE_TOKEN; return TINT;
}

{DIGIT}+"."{DIGIT}* {
  SAVE_TOKEN; return TFP;
}

"," {
  SAVE_TOKEN; return TCOMMA;
}


"(" {
  SAVE_TOKEN; return TOPENPAREN;
}

")" {
  SAVE_TOKEN; return TCLOSEPAREN;
}

"->" {
  SAVE_TOKEN; return TCAST;
}

"def" {
  SAVE_TOKEN; return TFUNC;
}

"var" {
  SAVE_TOKEN; return TVAR;
}

"double" {
  SAVE_TOKEN; return TWDOUBLE;
}

"int" {
  SAVE_TOKEN; return TWINT;
}

"array" {
  SAVE_TOKEN; return TWARRAY;
}

"extern" {
  SAVE_TOKEN; return TEXTERN;
}

"end" {
  SAVE_TOKEN; return TEND;
}

"{" {
  SAVE_TOKEN; return TOPENBRACE;
}

"}" {
  SAVE_TOKEN; return TCLOSEBRACE;
}

"[" {
  SAVE_TOKEN; return TOPENBRACKET;
}

"]" {
  SAVE_TOKEN; return TCLOSEBRACKET;
}

"*" {
  SAVE_TOKEN; return TSTAR;
}

"&" {
  SAVE_TOKEN; return TAND;
}

"+"|"-"|TSTAR|"/" {
  SAVE_TOKEN; return TBINOP;
}

"if" {
  SAVE_TOKEN; return TIF;
}

"else" {
  SAVE_TOKEN; return TELSE;
}

"then" {
  SAVE_TOKEN; return TTHEN;
}

"=" {
  SAVE_TOKEN; return TEQUALS;
}

{ID}        SAVE_TOKEN; return TIDENTIFIER;

#.* {/* one-line comment */}

[ \t\n]+

.           printf( "Unrecognized character: %s\n", yytext );

%%
