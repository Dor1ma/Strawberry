#ifndef STRAWBERRY_TOKENS_H
#define STRAWBERRY_TOKENS_H

enum Token {
    BERRY, SAY, RECIPE, INGREDIENT, IF, ELSE, FOR, IN, RANGE,
    TYPE_INT, TYPE_FLOAT, TYPE_BOOL, TYPE_STRING, TYPE_VOID, TYPE_LIST,
    OP_PLUS, OP_MINUS, OP_MULTIPLY, OP_DIVIDE, OP_GREATER, OP_LESS,
    OP_EQUAL, OP_NOT_EQUAL, OP_GREATER_EQUAL, OP_LESS_EQUAL,
    ASSIGN, DELIMITER, OP_UNKNOWN
};

typedef union {
    bool boolean;
    char* str;
    double num;
} YYSTYPE;

extern YYSTYPE yylval;

#endif //STRAWBERRY_TOKENS_H