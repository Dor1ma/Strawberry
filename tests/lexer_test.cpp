#include <iostream>
#include <fstream>
#include "lexer.yy.h"
#include "tokens.h"
#include <string>
#include <filesystem>

std::string tokenToString(int token) {
    switch (token) {
        case BERRY: return "BERRY";
        case SAY: return "SAY";
        case RECIPE: return "RECIPE";
        case INGREDIENT: return "INGREDIENT";
        case IF: return "IF";
        case ELSE: return "ELSE";
        case FOR: return "FOR";
        case IN: return "IN";
        case RANGE: return "RANGE";
        case TYPE_INT: return "TYPE_INT";
        case TYPE_FLOAT: return "TYPE_FLOAT";
        case TYPE_BOOL: return "TYPE_BOOL";
        case TYPE_STRING: return "TYPE_STRING";
        case TYPE_VOID: return "TYPE_VOID";
        case TYPE_LIST: return "TYPE_LIST";
        //case BOOLEAN: return "BOOLEAN";
        //case IDENTIFIER: return "IDENTIFIER";
        //case STRING: return "STRING";
        //case TOKEN_NUMBER: return "TOKEN_NUMBER";
        case OP_PLUS: return "OP_PLUS";
        case OP_MINUS: return "OP_MINUS";
        case OP_MULTIPLY: return "OP_MULTIPLY";
        case OP_DIVIDE: return "OP_DIVIDE";
        case OP_GREATER: return "OP_GREATER";
        case OP_LESS: return "OP_LESS";
        case OP_EQUAL: return "OP_EQUAL";
        case OP_NOT_EQUAL: return "OP_NOT_EQUAL";
        case OP_GREATER_EQUAL: return "OP_GREATER_EQUAL";
        case OP_LESS_EQUAL: return "OP_LESS_EQUAL";
        case ASSIGN: return "ASSIGN";
        case DELIMITER: return "DELIMITER";
        default: return "UNKNOWN_TOKEN";
    }
}

int main() {
    std::cout << "Current path: " << std::filesystem::current_path() << std::endl;
    std::ifstream testFile("../tests/test_input.txt");

    if (!testFile) {
        std::cerr << "Failed to open test_input.txt" << std::endl;
        return 1;
    }

    std::string content((std::istreambuf_iterator<char>(testFile)),
                        std::istreambuf_iterator<char>());
    testFile.close();

    YY_BUFFER_STATE buffer = yy_scan_string(content.c_str());

    int token;
    while ((token = yylex()) != 0) {
        std::cout << "Token: " << tokenToString(token) << " (" << yytext << ")\n";
    }

    yy_delete_buffer(buffer);

    return 0;
}
