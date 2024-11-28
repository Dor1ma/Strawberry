cmake_minimum_required(VERSION 3.10)
project(Strawberry)

set(CMAKE_CXX_STANDARD 17)

find_package(FLEX REQUIRED)
find_package(BISON REQUIRED)

# Добавляем директорию src/parser как include для Flex и Bison
include_directories(${CMAKE_SOURCE_DIR}/src/parser)

# Генерация файлов Flex и Bison
FLEX_TARGET(Lexer src/lexer/lexer.l ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.cpp
        DEFINES_FILE ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.h)

BISON_TARGET(Parser src/parser/parser.y ${CMAKE_CURRENT_BINARY_DIR}/parser.tab.cpp
        DEFINES_FILE ${CMAKE_CURRENT_BINARY_DIR}/parser.tab.h)

ADD_FLEX_BISON_DEPENDENCY(Lexer Parser)

# Источники и заголовки для библиотек
set(PARSER_SOURCES
        ${BISON_Parser_OUTPUTS}
        src/parser/ast.c
)
set(PARSER_HEADERS
        src/parser/ast.h
        src/parser/tokens.h
        ${CMAKE_CURRENT_BINARY_DIR}/parser.tab.h
)

set(LEXER_SOURCES
        ${FLEX_Lexer_OUTPUTS}
)
set(LEXER_HEADERS
        ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.h
)

# Создаем статическую библиотеку для парсера
add_library(ParserLib STATIC ${PARSER_SOURCES})
target_include_directories(ParserLib PUBLIC ${CMAKE_SOURCE_DIR}/src/parser ${CMAKE_CURRENT_BINARY_DIR})

# Создаем статическую библиотеку для лексера
add_library(LexerLib STATIC ${LEXER_SOURCES})
target_include_directories(LexerLib PUBLIC ${CMAKE_SOURCE_DIR}/src/lexer ${CMAKE_SOURCE_DIR}/src/parser ${CMAKE_CURRENT_BINARY_DIR})

# Добавляем внешние модули
add_subdirectory(external)
