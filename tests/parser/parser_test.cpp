#include <iostream>
#include <cassert>
#include <cstring>
#include <cstdio>
#include <unordered_set>
#include "ast.h"

// Прототипы функций парсера
extern int yyparse();
extern FILE* yyin;
ASTNode* root;

struct NodePointerHash {
    std::size_t operator()(const ASTNode* node) const {
        return reinterpret_cast<std::size_t>(node);
    }
};

struct NodePointerEqual {
    bool operator()(const ASTNode* lhs, const ASTNode* rhs) const {
        return lhs == rhs;
    }
};

// Функция для печати дерева AST с учетом зацикливания
void print_ast(ASTNode* root, std::unordered_set<const ASTNode*, NodePointerHash, NodePointerEqual>& visited) {
    if (root == NULL) {
        printf("NULL node\n");
        return;
    }

    // Если узел уже посещен, пропускаем его
    if (visited.find(root) != visited.end()) {
        printf("Cycle detected, skipping node\n");
        return;
    }

    // Помечаем узел как посещенный
    visited.insert(root);

    // Печать текущего узла в зависимости от его типа
    switch (root->type) {
        case NODE_IF:
            printf("IfStatement: \n");

            // Печать condition, если оно существует
            if (root->data.if_statement.condition != NULL) {
                printf("  Condition: ");
                print_ast(root->data.if_statement.condition, visited);
            } else {
                printf("  Condition: NULL\n");
            }

            // Печать body
            if (root->data.if_statement.body != NULL) {
                printf("\n  Body: ");
                print_ast(root->data.if_statement.body, visited);
            } else {
                printf("\n  Body: NULL\n");
            }

            // Печать else_branch, если оно существует
            if (root->data.if_statement.else_branch) {
                printf("\n  Else branch: ");
                print_ast(root->data.if_statement.else_branch, visited);
            } else {
                printf("  Else branch: NULL\n");
            }
            break;

        case NODE_LITERAL:
            printf("Literal: %f\n", root->data.literal.value);
            break;

        case NODE_IDENTIFIER:
            printf("Identifier: %s\n", root->data.identifier.name);
            break;

        case NODE_BINARY_OP:
            printf("Binary Operation: \n");
            if (root->data.binary_op.left != NULL) {
                print_ast(root->data.binary_op.left, visited);
            } else {
                printf("  Left operand: NULL\n");
            }

            printf(" Operator: %d\n", root->data.binary_op.op);

            if (root->data.binary_op.right != NULL) {
                print_ast(root->data.binary_op.right, visited);
            } else {
                printf("  Right operand: NULL\n");
            }
            break;

            // Печать других типов узлов, если они есть
        default:
            printf("Unknown node type\n");
            break;
    }

    // Рекурсивный вызов для следующего узла (если есть)
    if (root->next != NULL) {
        print_ast(root->next, visited);
    }
}

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

    std::unordered_set<const ASTNode*, NodePointerHash, NodePointerEqual> visited;
    print_ast(root, visited);

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
    //test_parser("berry name = \"Strawberry\"", "Variable declaration: string");
    //test_parser("berry l list = [1, 2, 3]", "Variable declaration: list");

    // Тесты if-else
    /*test_parser(
            "if x > 5 {}",
            "If statement"
    );*/

    // Тест цикла for
    /*test_parser(
            "for i in range(1, 10) { say i }",
            "For loop"
    );*/

    // Тест функции
   /* test_parser(
            "recipe divide(a: float, b: float) -> result: float { if b == 0 { return 0 } return a / b }",
            "Function declaration"
    );*/

    // Тест вызова функции
    //test_parser("berry result = divide(10.0, 2.0)", "Function call");

    // Тест сложения строк
    /*test_parser(
            "recipe greet(name: string) -> void { say \"Hello, \" + name }",
            "String concatenation in function"
    );*/

    std::cout << "All tests completed." << std::endl;
}

// Функция для печати дерева AST



int main() {
    run_tests();
    return 0;
}
