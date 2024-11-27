#include "ast.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

ASTNode* create_node(const char* name) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_IDENTIFIER;
    node->data.identifier.name = strdup(name);
    node->next = NULL;
    return node;
}

ASTNode* create_typed_variable_declaration(const char* name, const char* type, ASTNode* value) {
    ASTNode* node = create_node("typed_var_decl");
    node->type = NODE_BINARY_OP; // Используем NODE_BINARY_OP, так как нам нужно два дочерних узла
    node->data.binary_op.left = create_node(name); // имя переменной
    node->data.binary_op.right = value; // значение переменной

    // В этой части можно хранить тип переменной в строке или в другом месте, если нужно
    ASTNode* type_node = create_node(type);
    // Например, вы можете сохранить тип в виде отдельного узла, если необходимо
    node->data.binary_op.left->data.identifier.name = strdup(type);
    free(type_node);  // Освобождаем временный узел для типа
    return node;
}

ASTNode* create_if_statement(ASTNode* condition, ASTNode* body, ASTNode* elif_branch, ASTNode* else_branch) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_LIST; // Используем NODE_LIST для представления всех веток
    node->data.list.element_count = 4;
    node->data.list.elements = (ASTNode**)malloc(4 * sizeof(ASTNode*));
    node->data.list.elements[0] = condition;  // Условие
    node->data.list.elements[1] = body;       // Тело if
    node->data.list.elements[2] = elif_branch; // elif
    node->data.list.elements[3] = else_branch; // else
    return node;
}

ASTNode* create_for_loop(ASTNode* variable, ASTNode* start, ASTNode* end, ASTNode* body) {
    ASTNode* node = create_node("for_loop");
    node->type = NODE_LIST;
    node->data.list.elements = (ASTNode**)malloc(4 * sizeof(ASTNode*));
    node->data.list.elements[0] = variable; // переменная цикла
    node->data.list.elements[1] = start; // начальное значение
    node->data.list.elements[2] = end; // конечное значение
    node->data.list.elements[3] = body; // тело цикла
    node->data.list.element_count = 4; // количество элементов
    return node;
}

ASTNode* create_function_call(const char* name, ASTNode* args) {
    ASTNode* node = create_node("func_call");
    node->type = NODE_BINARY_OP; // Используем NODE_BINARY_OP для пары "имя функции" и "аргументы"
    node->data.binary_op.left = create_node(name); // имя функции
    node->data.binary_op.right = args; // аргументы функции
    return node;
}

ASTNode* create_argument_list(ASTNode* expr) {
    expr->next = NULL;
    return expr;
}

ASTNode* append_arguments(ASTNode* list, ASTNode* expr) {
    if (!list) return create_argument_list(expr);
    ASTNode* current = list;
    while (current->next) current = current->next;
    current->next = create_argument_list(expr); // Создаём новый элемент списка
    return list;
}

ASTNode* append_statements(ASTNode* list, ASTNode* statement) {
    if (!list) return statement;
    ASTNode* current = list;
    while (current->next) current = current->next;
    current->next = statement;
    return list;
}

ASTNode* create_variable_declaration(const char* name, ASTNode* value) {
    ASTNode* node = create_node("var_decl");
    node->type = NODE_BINARY_OP;
    node->data.binary_op.left = create_node(name); // имя переменной
    node->data.binary_op.right = value; // значение
    return node;
}

ASTNode* create_function_declaration(const char* name, ASTNode* params, const char* return_type, ASTNode* body) {
    ASTNode* node = create_node("func_decl");
    node->type = NODE_LIST;
    node->data.list.elements = (ASTNode**)malloc(4 * sizeof(ASTNode*));
    node->data.list.elements[0] = create_node(name); // имя функции
    node->data.list.elements[1] = params; // параметры
    node->data.list.elements[2] = create_node(return_type); // возвращаемый тип
    node->data.list.elements[3] = body; // тело функции
    node->data.list.element_count = 4;
    return node;
}

ASTNode* append_parameters(ASTNode* list, ASTNode* parameter) {
    return append_statements(list, parameter);
}

ASTNode* create_parameter(const char* name, const char* type) {
    ASTNode* node = create_node("parameter");
    node->type = NODE_BINARY_OP;
    node->data.binary_op.left = create_node(name); // имя параметра
    node->data.binary_op.right = create_node(type); // тип параметра
    return node;
}

ASTNode* create_literal_string(const char* value) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_LITERAL;
    node->data.literal.string_value = strdup(value);
    return node;
}

ASTNode* create_list(ASTNode* elements) {
    ASTNode* node = (ASTNode*)malloc(sizeof(ASTNode));
    node->type = NODE_LIST;

    // Подсчитываем количество элементов в списке
    int count = 0;
    ASTNode* current = elements;
    while (current) {
        count++;
        current = current->next;
    }

    // Выделяем память под массив указателей
    node->data.list.elements = (ASTNode**)malloc(count * sizeof(ASTNode*));
    node->data.list.element_count = count;

    // Заполняем массив указателей
    current = elements;
    for (int i = 0; i < count; i++) {
        node->data.list.elements[i] = current;
        current = current->next;
    }

    return node;
}