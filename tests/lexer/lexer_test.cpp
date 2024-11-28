#include <iostream>
#include <fstream>

#include <string>
#include <filesystem>
#include "parser.tab.h"
#include "lexer.yy.h"

YYSTYPE yylval; // Добавьте это определение

std::string tokenToString(int token) {
    switch (token) {
        case BERRY: return "BERRY";
        case SAY: return "SAY";
        case RECIPE: return "RECIPE";
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
        case BISON_T_BOOLEAN: return "BOOLEAN";
        case BISON_T_IDENTIFIER: return "IDENTIFIER";
        case BISON_T_STRING: return "STRING";
        case BISON_T_NUMBER: return "NUMBER";
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
        case OPEN_PAREN: return "OPEN_PAREN";
        case CLOSE_PAREN: return "CLOSE_PAREN";
        case OPEN_BRACE: return "OPEN_BRACE";
        case CLOSE_BRACE: return "CLOSE_BRACE";
        case OPEN_BRACKET: return "OPEN_BRACKET";
        case CLOSE_BRACKET: return "CLOSE_BRACKET";
        case RETURN: return "RETURN";
        case ARROW: return "ARROW";
        case COMMA: return "COMMA";
        case COLON: return "COLON";
        case OP_MODULO: return "OP_MODULO";
        case ELIF: return "ELIF";
        default: return "UNKNOWN_TOKEN";
    }
}

int main() {
    std::cout << "Current path: " << std::filesystem::current_path() << std::endl;
    std::ifstream testFile("test_input.txt");

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