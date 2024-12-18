#include <type_traits>

#include "../../globals.h"
#include "ast.h"

#include "llvm/IR/Constants.h"
#include "llvm/Support/raw_ostream.h"

using namespace llvm;

Value* IfAST::codeGen() {
	Value* conditionV = condition->codeGen();
	assert(conditionV != nullptr);

	conditionV = mBuilder.CreateSIToFP(conditionV, dType);
	conditionV = mBuilder.CreateFCmpONE(conditionV, ConstantFP::get(mContext, APFloat(0.0)), "ifcond");

	auto* func = currentFunc;
	auto* bodyBlock = BasicBlock::Create(mContext, "then", func);
	auto* elseBlock = BasicBlock::Create(mContext, "else");
	auto* endBlock = BasicBlock::Create(mContext, "end");

	mBuilder.CreateCondBr(conditionV, bodyBlock, elseStatements != nullptr ? elseBlock : endBlock);

	mBuilder.SetInsertPoint(bodyBlock);

	Value* lastLineIf = ConstantInt::get(mContext, APInt(8, 0));
	for(auto* line : ifStatments->statements) {
		lastLineIf = line->codeGen();
		delete line;
	}

	ifStatments->statements.clear();

	mBuilder.CreateBr(endBlock);
	bodyBlock = mBuilder.GetInsertBlock();

	if(elseStatements == nullptr) {
		elseStatements = new BlockAST();
	}

	func->getBasicBlockList().push_back(elseBlock);
	mBuilder.SetInsertPoint(elseBlock);

	Value* lastLineElse =
		ConstantInt::get(mContext, APInt(8, 0));
	for(auto* line : elseStatements->statements) {
		lastLineElse = line->codeGen();
		delete line;
	}

	elseStatements->statements.clear();

	mBuilder.CreateBr(endBlock);
	elseBlock = mBuilder.GetInsertBlock();

	func->getBasicBlockList().push_back(endBlock);
	mBuilder.SetInsertPoint(endBlock);

	assert(lastLineElse->getType() == lastLineIf->getType() && "Both return types must be equal");

	auto* ternaryNode = mBuilder.CreatePHI(lastLineElse->getType(), 2, "iftmp");
	ternaryNode->addIncoming(lastLineIf, bodyBlock);
	ternaryNode->addIncoming(lastLineElse, elseBlock);

	return ternaryNode;
}

std::string IfAST::out() {
	return std::string("Not Implemented");
}