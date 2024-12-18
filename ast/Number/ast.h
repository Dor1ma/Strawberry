#pragma once

#include "llvm/IR/Value.h"

#include "../ast.h"

using namespace llvm;

class NumberAST : public AST {
	double val;
	Type* type;

public:
	NumberAST(double val, Type* type)
		: val(val)
		, type(type) {}
	Value* codeGen() override;
	std::string out() override;
};
