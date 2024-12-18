%code requires { 
#include "../ast/Number/ast.h"
#include "../ast/BinOp/ast.h"
#include "../ast/Functions/ast.h"
#include "../ast/Variables/ast.h"
#include "../ast/ControlFlow/ast.h"
#include "../ast/ast.h"
#include "../globals.h"

#include "llvm/IR/Constants.h"
}

%{

#include <iostream>
#include <cctype>
#include <cstring>
#include <vector>
#include <stack>
#include <string>
#include <utility>

#include "../ast/ast.h"

#define YYERROR_VERBOSE

void yyerror (const char *error) {
    std::cerr << error << std::endl;
    std::cerr << 'test';
    std::string *e;
    e->size();
}

extern int  yylex ();

BlockAST *res;

%}

%union {
    std::string *string;
    NumberAST *number;
    BlockAST *base;
    BinOpAST *binOp;
    AST *any;
    FunctionAST *func;
    std::vector<std::pair<std::string, llvm::Type *>> *args;
    std::vector<Type *> *typeArgs;
    std::vector<AST *> *valueArgs;
    llvm::Type *type;
}

%token <string> TIDENTIFIER TINT TBINOP TOPENPAREN TCLOSEPAREN TOPENBRACE TCLOSEBRACE TFUNC TVAR TCAST TWDOUBLE TWINT TEQUALS TEXTERN TAND TOPENBRACKET TCLOSEBRACKET TSTAR TEND TCOLON TCOMMA TWARRAY TIF TELSE TTHEN TFP
%token <args> arg

%type <any> constExpr func var cast extern call array varArray ifStatement
%type <base> program statments optionalElseStatement
%type <args> funcArgs optionalArgParens
%type <typeArgs> typeArgs optionalTypeArgParens
%type <type> type
%type <valueArgs> valueArgs arrayValues
%type <type> returnType varType

%start program

%%
    program   :  { res = new BlockAST(); }
    | statments { res = $1; }
          ;
    statments : constExpr { $$ = new BlockAST(); $$->statements.push_back($<any>1); }
    | statments constExpr { $1->statements.push_back($<any>2); }
          ;
    constExpr : TINT { $$ = new NumberAST(std::stoi(*$1), i32); }
    | TFP { $$ = new NumberAST(std::stoi(*$1), dType); }
    | func { $$ = $<any>1; }
    | var { $$ = $<any>1; }
    | cast { $$ = $<any>1; }
    | extern { $$ = $<any>1; }
    | call { $$ = $<any>1; }
    | TIDENTIFIER { $$ = new VariableGetAST(*$<string>1); }
    | array { $$ = $<any>1; }
    | constExpr TBINOP constExpr { $$ = new BinOpAST($1, $3, $2->c_str()[0]); }
    | constExpr TSTAR constExpr { $$ = new BinOpAST($1, $3, $2->c_str()[0]); }
    | constExpr TSTAR { $$ = new LoadVariableAST($<any>1); }
    | TIDENTIFIER TAND { $$ = new VariableGetAST(*$<string>1, true); }
    | TOPENPAREN constExpr TCLOSEPAREN { $$ = $2; }
    | varArray { $$ = $<any>1; }
    | ifStatement { $$ = $1; }
    | { $$ = 0; }
          ;
    var         : TIDENTIFIER TEQUALS constExpr varType {
        $$ = new VariableAST(*$1, $<any>3, $4);
    }
          ;
    varType     : TOPENPAREN type TCLOSEPAREN {
        $$ = $2;
    }
    | {
        $$ = 0;
    }
    varArray    : TOPENBRACKET arrayValues TCLOSEBRACKET {
        $$ = new ArrayAST(*$<valueArgs>2);
    }
          ;
    arrayValues : constExpr {
        auto *tmpValues = new std::vector<AST *>();
        if ($1) tmpValues->push_back($1);
        $$ = tmpValues;
    }
    | arrayValues TCOMMA constExpr {
        $1->push_back($<any>3);
    }
          ;

    array       : TIDENTIFIER TOPENBRACKET constExpr TCLOSEBRACKET {
        $$ = new ArrayGetAST(*$<string>1, $<any>3);
    }
          ;
    call        : TIDENTIFIER TOPENPAREN valueArgs TCLOSEPAREN {
        $$ = new CallAST(*$<string>1, *$<valueArgs>3);
    }
          ;
    func        : TFUNC TIDENTIFIER optionalArgParens returnType statments TEND {
        $$ = new FunctionAST(*$<string>2, *$<args>3, $<base>5, $4);
    }
          ;
    returnType  : TCAST type {
        $$ = $2;
    }
    | {
        $$ = i8;
    }
          ;
    optionalArgParens : TOPENPAREN funcArgs TCLOSEPAREN {
        $$ = $2;
    }
    | funcArgs {
        $$ = $1;
    }
    	  ;
    optionalTypeArgParens : TOPENPAREN typeArgs TCLOSEPAREN {
    	$$ = $2;
    }
    | typeArgs {
	$$ = $1;
    }
          ;
    funcArgs    : {
        auto tmpArgs = std::vector<std::pair<std::string, llvm::Type *>>();
        $$ = &tmpArgs;
    }
    | funcArgs seperator type TIDENTIFIER {
        $1->push_back(std::make_pair(*$<string>4, $<type>3));
    }
          ;
    typeArgs    : {
        auto tmpArgs = std::vector<Type *>();
        $$ = &tmpArgs;
    }
    | typeArgs seperator type {
        $1->push_back($<type>3);
    }
          ;
    valueArgs   : constExpr {
        auto *tmpArgs = new std::vector<AST *>();
        if ($1) tmpArgs->push_back($1);
        $$ = tmpArgs;
    }
    | valueArgs TCOMMA constExpr {
        $1->push_back($<any>3);
    }
          ;
    seperator   : TCOMMA { }
    | { }
          ;
    cast        : TIDENTIFIER TCAST type { $$ = new CastAST(*$<string>1, $<type>3); }
          ;
    type        : TWDOUBLE { $$ = dType; }
    | TWINT { $$ = i32; }
    | TWARRAY { $$ = pi8; }
    | type TSTAR { $$ = PointerType::getUnqual($<type>1); }
    | TFUNC optionalTypeArgParens returnType {
    	auto* tmpArgs = $2;
    	auto* FTH = new FunctionTypeHelper($3, *tmpArgs);
    	$$ = PointerType::getUnqual(FTH->codeGen());
     }
          ;
    extern      : TEXTERN TIDENTIFIER TOPENPAREN typeArgs TCLOSEPAREN returnType {
        $$ = new PrototypeAST(*$<string>2, *$<typeArgs>4, $<type>6);
    }
          ;
    ifStatement	: TIF TOPENPAREN constExpr TCLOSEPAREN statments optionalElseStatement {
	$$ = new IfAST($3, $5, $6);
    }
    optionalElseStatement : TELSE statments TEND {
	$$ = $<base>2;
    }
    | TEND {
    	$$ = nullptr;
    }

%%
