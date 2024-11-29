%{
#include "ast.h"
#include "ast.c"
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
    #include "ast.h"
}
%union {
    ASTNode* node;
    double num;
    char* str;
    int boolean;
}


%token OPEN_BRACKET CLOSE_BRACKET
%token OP_MULTIPLY OP_DIVIDE OP_GREATER_EQUAL OP_LESS_EQUAL ASSIGN OP_UNKNOWN
%token <str> BISON_T_STRING
%token OPEN_PAREN CLOSE_PAREN OPEN_BRACE CLOSE_BRACE
%token OP_PLUS OP_MINUS OP_EQUAL OP_NOT_EQUAL OP_GREATER OP_LESS OP_MODULO
%token TYPE_INT TYPE_FLOAT TYPE_BOOL TYPE_STRING TYPE_VOID TYPE_LIST
%token IF ELSE ELIF FOR IN RANGE RETURN
%token ARROW COMMA COLON
%token BERRY SAY RECIPE
%token <num> BISON_T_NUMBER
%token <str> BISON_T_IDENTIFIER
%token <boolean> BISON_T_BOOLEAN


%type <node> program statement_list statement variable_declaration function_declaration control_statement expression argument_list parameter_list parameter else_branch


%nonassoc IFX
%nonassoc ELSE


%%

program:
    statement_list { root = $1; }
;

statement_list:
    statement
    | statement_list statement { $$ = append_statements($1, $2); }
;

statement:
    expression { $$ = $1; }
    | variable_declaration { $$ = $1; }
    | function_declaration { $$ = $1; }
    | control_statement { $$ = $1; }
;

variable_declaration:
    BERRY BISON_T_IDENTIFIER ASSIGN expression {
        $$ = create_variable_declaration($2, $4);
    }
    | BERRY BISON_T_IDENTIFIER TYPE_INT ASSIGN expression {
        $$ = create_typed_variable_declaration($2, "int", $5);
    }
    | BERRY BISON_T_IDENTIFIER TYPE_LIST ASSIGN expression {
        $$ = create_typed_variable_declaration($2, "list", $5);
    }
;

function_declaration:
    RECIPE BISON_T_IDENTIFIER OPEN_PAREN parameter_list CLOSE_PAREN ARROW TYPE_FLOAT OPEN_BRACE statement_list CLOSE_BRACE {
        $$ = create_function_declaration($2, $4, "float", $<node>8);
    }
;

parameter_list:
    /* empty */ { $$ = NULL; }
    | parameter
    | parameter_list COMMA parameter { $$ = append_parameters($1, $3); }
;

parameter:
    BISON_T_IDENTIFIER COLON TYPE_FLOAT { $$ = create_parameter($1, "float"); }
;

control_statement:
    IF expression COLON statement_list else_branch {
        $$ = create_if_statement($2, $4, $5);
    }
;

else_branch:
    /* empty */ { $$ = NULL; }
    | ELSE COLON statement_list {
        $$ = $3;  // Прямая ссылка на тело else
    }
    | ELIF expression COLON statement_list else_branch {
        $$ = create_if_statement($2, $4, $5); // Рекурсивное создание elif
    }
;



expression:
    BISON_T_NUMBER { $$ = create_literal($1); }
    | BISON_T_IDENTIFIER { $$ = create_identifier($1); }
    | expression OP_PLUS expression { $$ = create_binary_op($1, '+', $3); }
    | expression OP_MINUS expression { $$ = create_binary_op($1, '-', $3); }
    | expression OP_MULTIPLY expression { $$ = create_binary_op($1, '*', $3); }
    | expression OP_DIVIDE expression { $$ = create_binary_op($1, '/', $3); }
    | expression OP_MODULO expression { $$ = create_binary_op($1, '%', $3); }
    | OPEN_PAREN expression CLOSE_PAREN { $$ = $2; }
    | BISON_T_IDENTIFIER OPEN_PAREN argument_list CLOSE_PAREN {
        $$ = create_function_call($1, $3);
    }
    | BISON_T_STRING { $$ = create_literal_string($1); }
    | OPEN_BRACKET argument_list CLOSE_BRACKET { $$ = create_list($2); }
;

argument_list:
    /* empty */ { $$ = NULL; }
    | expression { $$ = create_argument_list($1); }
    | argument_list COMMA expression { $$ = append_arguments($1, $3); }
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