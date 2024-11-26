#include "C:\\Users\\ZGHTL\\CLionProjects\\Strawberry\\parser.tab.h"
#include "../../cmake-build-debug/lexer.yy.h"
#include <windows.h> // Для GetTempPath

ASTNode* root; // Объявление переменной root

// Прототип функции serialize_ast
void serialize_ast(ASTNode* node, char* buffer, size_t size);

// Функция для проверки теста
void run_test(const char* test_name, const char* input, const char* expected_output) {
    // Получаем путь к временной директории
    char temp_filename[MAX_PATH];
    GetTempPath(MAX_PATH, temp_filename);
    strcat(temp_filename, "temp_input.txt"); // Имя временного файла

    FILE* temp_input = fopen(temp_filename, "w+");
    if (!temp_input) {
        fprintf(stderr, "Failed to create temporary input file for test '%s'\n", test_name);
        return;
    }
    fputs(input, temp_input);
    rewind(temp_input);

    yyin = temp_input;
    root = NULL;  // Ensure root is initialized to NULL before parsing

    if (yyparse() == 0 && root != NULL) {
        char actual_output[1024];
        serialize_ast(root, actual_output, sizeof(actual_output));

        if (strcmp(actual_output, expected_output) == 0) {
            printf("[PASS] %s\n", test_name);
        } else {
            printf("[FAIL] %s\n", test_name);
            printf("  Expected: %s\n", expected_output);
            printf("  Got:      %s\n", actual_output);
        }
    } else {
        printf("[FAIL] %s\n", test_name);
        printf("  Parser failed to process input: %s\n", input);
    }

    fclose(temp_input);
    remove(temp_filename); // Удаление временного файла
}

void serialize_ast(ASTNode* node, char* buffer, size_t size) {
    if (!node) {
        snprintf(buffer, size, "NULL");
        return;
    }

    switch (node->type) {
        case NODE_LITERAL:
            snprintf(buffer, size, "LITERAL(%f)", node->data.literal.value);
            break;
        case NODE_IDENTIFIER:
            snprintf(buffer, size, "IDENTIFIER(%s)", node->data.identifier.name);
            break;
        case NODE_BINARY_OP: {
            char left[256], right[256];
            serialize_ast(node->data.binary_op.left, left, sizeof(left));
            serialize_ast(node->data.binary_op.right, right, sizeof(right));
            snprintf(buffer, size, "BINARY_OP(%s, '%c', %s)", left, node->data.binary_op.op, right);
            break;
        }
        default:
            snprintf(buffer, size, "UNKNOWN");
            break;
    }
}

int main() {
    run_test("Test 1: Literal number", "42;", "LITERAL(42.000000)");
    run_test("Test 2: Identifier", "x;", "IDENTIFIER(x)");
    run_test("Test 3: Binary operation", "3 + 5;", "BINARY_OP(LITERAL(3.000000), '+', LITERAL(5.000000))");
    return 0;
}
