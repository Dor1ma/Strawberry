#ifndef AST_H
#define AST_H

#define NODE_LITERAL 1
#define NODE_IDENTIFIER 2
#define NODE_BINARY_OP 3
#define NODE_LIST 4

typedef struct ASTNode ASTNode;

typedef struct {
    double value;  // Для числовых литералов
    char* string_value;  // Для строковых литералов
} LiteralData;

typedef struct {
    char* name;
} IdentifierData;

typedef struct {
    ASTNode* left;
    int op;  // Операция
    ASTNode* right;
} BinaryOpData;

typedef struct {
    ASTNode** elements;  // Массив указателей на элементы списка
    int element_count;   // Количество элементов
} ListData;

struct ASTNode {
    int type;  // Тип узла
    union {
        LiteralData literal;  // Литералы
        IdentifierData identifier;  // Идентификаторы
        BinaryOpData binary_op;  // Бинарные операции
        ListData list;  // Списки
    } data;
    ASTNode* next;  // Следующий узел в списке (для построения цепочек)
};

// Функции для создания AST
ASTNode* create_literal(double value);
ASTNode* create_literal_string(const char* value);
ASTNode* create_identifier(char* name);
ASTNode* create_binary_op(ASTNode* left, int op, ASTNode* right);
ASTNode* create_list(ASTNode* elements);

ASTNode* append_statements(ASTNode* list, ASTNode* statement);
ASTNode* create_variable_declaration(const char* name, ASTNode* value);
ASTNode* create_typed_variable_declaration(const char* name, const char* type, ASTNode* value);
ASTNode* create_if_statement(ASTNode* condition, ASTNode* body, ASTNode* elif_branch, ASTNode* else_branch);
ASTNode* create_for_loop(ASTNode* variable, ASTNode* start, ASTNode* end, ASTNode* body);
ASTNode* create_function_declaration(const char* name, ASTNode* params, const char* return_type, ASTNode* body);
ASTNode* append_parameters(ASTNode* list, ASTNode* parameter);
ASTNode* create_parameter(const char* name, const char* type);
#endif // AST_H
