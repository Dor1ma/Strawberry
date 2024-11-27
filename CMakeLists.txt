cmake_minimum_required(VERSION 3.10)
project(Strawberry)

set(CMAKE_CXX_STANDARD 17)

find_package(FLEX REQUIRED)
find_package(BISON REQUIRED)

FLEX_TARGET(Lexer src/lexer/lexer.l ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.cpp
        DEFINES_FILE ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.h)

BISON_TARGET(Parser src/parser/parser.y ${CMAKE_BINARY_DIR}/parser.tab.cpp
        DEFINES_FILE ${CMAKE_CURRENT_BINARY_DIR}/parser.tab.h)

ADD_FLEX_BISON_DEPENDENCY(Lexer Parser)

set(SOURCE_FILES
        ${FLEX_Lexer_OUTPUTS}
        ${BISON_Parser_OUTPUTS}
        src/parser/ast.h
        src/parser/tokens.h
        src/parser/ast.c
)

# Добавляем каталог src/parser в пути для include
include_directories(${CMAKE_SOURCE_DIR}/src/parser)

add_executable(Strawberry ${SOURCE_FILES})
target_include_directories(Strawberry PRIVATE ${CMAKE_SOURCE_DIR}/src/lexer)
target_include_directories(Strawberry PRIVATE ${CMAKE_SOURCE_DIR}/src/parser)
