#include <type_traits>

#include "../../globals.h"
#include "ast.h"

#include "llvm/IR/Constants.h"
#include "llvm/Support/raw_ostream.h"

using namespace llvm;

extern LLVMContext mContext;

Value* NumberAST::codeGen() {
	if(type->isDoubleTy())
		return ConstantFP::get(mContext, APFloat(val));
	return ConstantInt::get(mContext, APInt(32, (uint64_t)val));
}

std::string NumberAST::out() {
	std::string outString;
	llvm::raw_string_ostream rso(outString);
	this->codeGen()->print(rso);

	return outString;
}
