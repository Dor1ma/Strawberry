#ifndef AST_H
#define AST_H

typedef enum {
    NODE_LITERAL,
    NODE_IDENTIFIER,
    NODE_BINARY_OP
} NodeType;

typedef struct ASTNode {
    NodeType type;
    union {
        struct {
            double value;
        } literal;
        struct {
            char* name;
        } identifier;
        struct {
            struct ASTNode* left;
            int op;
            struct ASTNode* right;
        } binary_op;
    } data;
} ASTNode;

#endif // AST_H
