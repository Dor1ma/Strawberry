#ifndef STRAWBERRY_TOKENS_H
#define STRAWBERRY_TOKENS_H

enum Token {
    BERRY = 300, SAY, RECIPE, IF, ELSE, ELIF, FOR, IN, RANGE, RETURN, ARROW,
    TYPE_INT, TYPE_FLOAT, TYPE_BOOL, TYPE_STRING, TYPE_VOID, TYPE_LIST,
    OP_PLUS, OP_MINUS, OP_MULTIPLY, OP_DIVIDE, OP_GREATER, OP_LESS,
    OP_EQUAL, OP_NOT_EQUAL, OP_GREATER_EQUAL, OP_LESS_EQUAL, OP_MODULO,
    ASSIGN, DELIMITER, ID, OP_UNKNOWN,
    OPEN_PAREN, CLOSE_PAREN, OPEN_BRACE, CLOSE_BRACE, OPEN_BRACKET, CLOSE_BRACKET, COMMA, COLON,
    BISON_T_BOOLEAN, BISON_T_IDENTIFIER, BISON_T_STRING, BISON_T_NUMBER,
};

typedef union {
    bool boolean;
    char* str;
    double num;
} YYSTYPE;

extern YYSTYPE yylval;

#endif //STRAWBERRY_TOKENS_H