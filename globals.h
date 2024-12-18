#pragma once

#include "llvm/IR/IRBuilder.h"
#include "llvm/IR/LLVMContext.h"
#include "llvm/IR/Value.h"

using namespace llvm;


extern LLVMContext mContext;
extern IRBuilder<> mBuilder;
extern std::unique_ptr<Module> mModule;
extern std::map<std::string, AllocaInst*> namedVariables;
extern std::map<std::string, Value*> namedArgs;
extern Function* currentFunc;

extern AllocaInst* CreateBlockAlloca(Function* func, std::string name, Type* type);


extern Type* i32;
extern Type* pi32;
extern Type* i8;
extern Type* pi8;
extern Type* dType;
extern Type* pdType;
extern StructType* aType;


extern std::string ARCCurrentFunc;
extern std::map<std::string, std::map<std::string, int>> ARC;
