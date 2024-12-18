#pragma once

#include "llvm/IR/Value.h"

#include "../ast.h"

using namespace llvm;

class BinOpAST : public AST {
	AST* L;
	AST* R;
	int op;

public:
	BinOpAST(AST* L, AST* R, int op)
		: L(std::move(L))
		, R(std::move(R))
		, op(op) {}
	Value* codeGen() override;
	std::string out() override;
};
