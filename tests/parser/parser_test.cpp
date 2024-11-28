#include <iostream>
#include <cassert>
#include <cstring>
#include <cstdio>
#include "ast.h"

// Прототипы функций парсера
extern int yyparse();
extern FILE* yyin;
ASTNode* root;

// Утилита для тестирования
void test_parser(const char* input, const char* test_name) {
    std::cout << "Running test: " << test_name << std::endl;

    FILE* old_stdin = yyin;
    const char* temp_filename = "temp_test_file.txt";
    FILE* temp_file = fopen(temp_filename, "w+");
    assert(temp_file != nullptr);

    fputs(input, temp_file);
    rewind(temp_file);

    yyin = temp_file;
    root = nullptr;

    if (yyparse() == 0) {
        std::cout << "[PASS] " << test_name << std::endl;
    } else {
        std::cerr << "[FAIL] " << test_name << " - Parsing failed!" << std::endl;
    }

    fclose(temp_file);
    remove(temp_filename);
    yyin = old_stdin;

    std::cout << "Finished test: " << test_name << std::endl;
}


// Функции для проверки AST (Пример проверки узлов)
void check_variable_declaration(ASTNode* node, const char* name, int expected_type) {
    assert(node != nullptr);
    assert(node->type == NODE_BINARY_OP);
    assert(strcmp(node->data.binary_op.left->data.identifier.name, name) == 0);
    assert(node->data.binary_op.right->type == expected_type);
}

void run_tests() {
    // Тесты переменных
    test_parser("berry x = 10", "Variable declaration: integer");
    test_parser("berry name = \"Strawberry\"", "Variable declaration: string");
    test_parser("berry l list = [1, 2, 3]", "Variable declaration: list");

    // Тесты if-else
    test_parser(
            "if x > 5 {}",
            "If statement"
    );

    // Тест цикла for
    test_parser(
            "for i in range(1, 10) { say i }",
            "For loop"
    );

    // Тест функции
    test_parser(
            "recipe divide(a: float, b: float) -> result: float { if b == 0 { return 0 } return a / b }",
            "Function declaration"
    );

    // Тест вызова функции
    test_parser("berry result = divide(10.0, 2.0)", "Function call");

    // Тест сложения строк
    test_parser(
            "recipe greet(name: string) -> void { say \"Hello, \" + name }",
            "String concatenation in function"
    );

    std::cout << "All tests completed." << std::endl;
}

int main() {
    run_tests();
    return 0;
}
