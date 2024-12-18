#include "main.h"
#include "ast/Functions/ast.h"
#include "ast/ast.h"
#include "globals.h"

#include <fstream>
#include <iostream>
#include <llvm/Support/raw_ostream.h>
#include <utility>

extern FILE* yyin;
extern int yyparse();
extern BlockAST* res;

using namespace llvm;

std::string externMalloc() {
	std::vector<Type*> argTypes = {i32};
	auto proto = std::make_unique<PrototypeAST>("malloc", std::move(argTypes), pi8);
	return proto->out();
}

std::string externFree() {
	std::vector<Type*> argTypes = {pi8};
	auto proto = std::make_unique<PrototypeAST>("free", std::move(argTypes), i8);
	return proto->out();
}

int main(int argc, char** argv) {
	if(argc > 1) {
		yyin = fopen(argv[1], "r");
	} else {
		throw std::invalid_argument("Did not receive file as second argument!");
	}

	yyparse();

	auto output = externMalloc() + externFree() + res->out();

	std::cout << output << std::endl;
	std::ofstream outputFile;
	outputFile.open(argv[2]);
	outputFile << output;
}
