#include <iostream>
#include <type_traits>

#include "../../globals.h"
#include "ast.h"

#include "llvm/IR/Constants.h"
#include "llvm/IR/Instructions.h"
#include "llvm/Support/raw_ostream.h"

using namespace llvm;

Value* VariableAST::codeGen() {
	Function* func = currentFunc;

	Value* newValue = value->codeGen();
	auto variableAlloca = namedVariables[name];

	if(variableAlloca) {
		mBuilder.CreateStore(newValue, variableAlloca);
		return newValue;
	}

	if(type == nullptr)
		type = newValue->getType();

	AllocaInst* alloca = CreateBlockAlloca(func, name, type);
	mBuilder.CreateStore(newValue, alloca);
	namedVariables[name] = alloca;

	return alloca;
}

Value* VariableSetAST::codeGen() {
	Value* loadAlloca = mBuilder.CreateLoad(namedVariables[name]);
	Value* newValue = newVal->codeGen();

	mBuilder.CreateStore(newValue, loadAlloca);
	return newValue;
}

Value* PointerSetAST::codeGen() {
	return mBuilder.CreateStore(newVal->codeGen(), ptr->codeGen());
}

Value* VariableGetAST::codeGen() {
	Value* possibleFunction = mModule->getFunction(name);
	if(possibleFunction)
		return possibleFunction;

	Value* alloca = dyn_cast_or_null<Value>(namedVariables[name]);
	if(alloca == nullptr)
		return dyn_cast_or_null<Value>(namedArgs[name]);

	Value* retValue =
		getPointer ? alloca : mBuilder.CreateLoad(alloca);

	return retValue;
}

Value* CastAST::codeGen() {
	Value* V = namedVariables[name];
	if(V == nullptr)
		V = namedArgs[name];

	assert(V != nullptr && "Value to cast is null");

	Value* loadAlloca = mBuilder.CreateLoad(V);
	Instruction* castInst = new BitCastInst(loadAlloca, PointerType::getUnqual(to));
	mBuilder.Insert(castInst);

	AllocaInst* alloca = CreateBlockAlloca(currentFunc, name, castInst->getType());
	mBuilder.CreateStore(castInst, alloca);
	namedVariables[name] = alloca;
	return castInst;
}

Value* ArrayAST::codeGen() {
	std::vector<Value*> arrayValues;

	size_t arraySize = 0;
	for(auto* v : values) {
		Value* newValue = v->codeGen();
		arrayValues.push_back(newValue);

		auto type = newValue->getType();
		auto size = type->getPrimitiveSizeInBits() / 8;

		arraySize += size;
	}

	Function* mallocCall = mModule->getFunction("malloc");
	std::vector<Value*> mallocArgs = {ConstantInt::get(mContext, APInt(32, arraySize))};

	Value* mallocVal = mBuilder.CreateCall(mallocCall, mallocArgs);
	auto* arrayType = PointerType::getUnqual(arrayValues[0]->getType());
	auto* array = mBuilder.CreateBitCast(mallocVal, arrayType);

	auto count = 0;
	for(auto* v: arrayValues) {
		auto typeSize = v->getType()->getPrimitiveSizeInBits() / 8;
		if(typeSize == 0)
			typeSize = 4;

		auto* index =
			ConstantInt::get(mContext, APInt(32, count++));
		auto* element = mBuilder.CreateGEP(array, index);

		mBuilder.CreateStore(v, element);
	}

	return array;
}

Value* ArrayGetAST::codeGen() {
	auto a = namedVariables[arrayName];
	auto* array = mBuilder.CreateLoad(a);

	auto* arrayIndex = index->codeGen();
	assert(arrayIndex->getType()->isIntegerTy() && "Index must be an integer");
	assert(arrayIndex->getType()->getPrimitiveSizeInBits() == 32 && "Warring: non-32bit integer may have unexpected results");

	auto* element = mBuilder.CreateGEP(array, arrayIndex);
	return mBuilder.CreateLoad(element);
}

Value* LoadVariableAST::codeGen() {
	return mBuilder.CreateLoad(ptr->codeGen());
}

std::string CastAST::out() {
	return std::string("Not Implemented");
}

std::string LoadVariableAST::out() {
	return std::string("Not Implemented");
}

std::string VariableAST::out() {
	return std::string("Not Implemented");
}

std::string VariableSetAST::out() {
	return std::string("Not Implemented");
}

std::string VariableGetAST::out() {
	return std::string("Not Implemented");
}

std::string ArrayAST::out() {
	return std::string("Not Implemented");
}

std::string ArrayGetAST::out() {
	return std::string("Not Implemented");
}

std::string PointerSetAST::out() {
	return std::string("Not Implemented");
}
