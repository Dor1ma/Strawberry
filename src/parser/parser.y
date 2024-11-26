%{
#include "C:\Users\ZGHTL\CLionProjects\Strawberry\src\parser\ast.h"
#include <stdio.h>
#include <stdlib.h>

extern ASTNode* create_literal(double value);
extern ASTNode* create_identifier(char* name);
extern ASTNode* create_binary_op(ASTNode* left, int op, ASTNode* right);

void yyerror(const char* s);
int yylex(void);

extern ASTNode* root;
%}

%code requires {
    #include "C:\\Users\\ZGHTL\\CLionProjects\\Strawberry\\src\\parser\\ast.h"
}
%union {
    ASTNode* node;
    double num;
    char* str;
    int boolean;
}


%token OPEN_BRACKET CLOSE_BRACKET
%token OP_MULTIPLY OP_DIVIDE OP_GREATER_EQUAL OP_LESS_EQUAL ASSIGN OP_UNKNOWN
%token BISON_T_STRING
%token OPEN_PAREN CLOSE_PAREN OPEN_BRACE CLOSE_BRACE
%token OP_PLUS OP_MINUS OP_EQUAL OP_NOT_EQUAL OP_GREATER OP_LESS OP_MODULO
%token TYPE_INT TYPE_FLOAT TYPE_BOOL TYPE_STRING TYPE_VOID TYPE_LIST
%token IF ELSE ELIF FOR IN RANGE RETURN
%token ARROW COMMA COLON
%token BERRY SAY RECIPE
%token <num> BISON_T_NUMBER
%token <str> BISON_T_IDENTIFIER
%token <boolean> BISON_T_BOOLEAN



%type <node> expression statement program

%%

program:
    statement { printf("Parsing complete\n"); root = $1; }
;

statement:
    expression ';' { $$ = $1; }
;

expression:
    BISON_T_NUMBER {
        $$ = create_literal($1);
    }
    | BISON_T_IDENTIFIER {
        $$ = create_identifier($1);
    }
    | expression OP_PLUS expression {
        $$ = create_binary_op($1, '+', $3);
    }
    | expression OP_MINUS expression {
        $$ = create_binary_op($1, '-', $3);
    }
    | OPEN_PAREN expression CLOSE_PAREN {
        $$ = $2;
    }
;

%%

// Реализация функций для AST

extern ASTNode* create_literal(double value) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_LITERAL;
    node->data.literal.value = value;
    return node;
}

extern ASTNode* create_identifier(char* name) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_IDENTIFIER;
    node->data.identifier.name = name;
    return node;
}

extern ASTNode* create_binary_op(ASTNode* left, int op, ASTNode* right) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_BINARY_OP;
    node->data.binary_op.left = left;
    node->data.binary_op.op = op;
    node->data.binary_op.right = right;
    return node;
}

void yyerror(const char* s) {
    fprintf(stderr, "Error: %s\n", s);
}