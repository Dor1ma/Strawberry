#pragma once

#include "llvm/IR/Value.h"

#include "../ast.h"

using namespace llvm;

class VariableAST : public AST {
	std::string name;
	AST* value;
	Type* type;

public:
	VariableAST(std::string name, AST* value, Type* type)
		: name(name)
		, value(value)
		, type(type) {}
	Value* codeGen() override;
	std::string out() override;
};

class VariableSetAST : public AST {
	std::string name;
	AST* newVal;

public:
	VariableSetAST(std::string name, AST* newVal)
		: name(name)
		, newVal(newVal) {}
	Value* codeGen() override;
	std::string out() override;
};

class PointerSetAST : public AST {
	AST* ptr;
	AST* newVal;

public:
	PointerSetAST(AST* ptr, AST* newVal)
		: ptr(ptr)
		, newVal(newVal) {}
	Value* codeGen() override;
	std::string out() override;
};

class VariableGetAST : public AST {
	std::string name;
	bool getPointer;

public:
	VariableGetAST(std::string name, bool getPointer = false)
		: name(name)
		, getPointer(getPointer) {}
	Value* codeGen() override;
	std::string out() override;
};

class CastAST : public AST {
	std::string name;
	Type* to;

public:
	CastAST(std::string name, Type* to)
		: name(name)
		, to(to) {}
	Value* codeGen() override;
	std::string out() override;
};

class ArrayAST : public AST {
	std::vector<AST*> values;

public:
	ArrayAST(std::vector<AST*> values)
		: values(values) {}
	Value* codeGen() override;
	std::string out() override;
};

class ArrayGetAST : public AST {
	std::string arrayName;
	AST* index;

public:
	ArrayGetAST(std::string arrayName, AST* index)
		: arrayName(arrayName)
		, index(index) {}
	Value* codeGen() override;
	std::string out() override;
};

class LoadVariableAST : public AST {
	AST* ptr;

public:
	LoadVariableAST(AST* ptr)
		: ptr(ptr) {}
	Value* codeGen() override;
	std::string out() override;
};
