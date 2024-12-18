cmake_minimum_required(VERSION 3.10)
project(Strawberry)

set(CMAKE_CXX_STANDARD 17)

find_package(FLEX REQUIRED)
find_package(BISON REQUIRED)

FLEX_TARGET(Lexer lexer/lexer.l ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.cpp DEFINES_FILE ${CMAKE_CURRENT_BINARY_DIR}/lexer.yy.h)
BISON_TARGET(Parser parser/parser.y ${CMAKE_BINARY_DIR}/parser.tab.cpp
        DEFINES_FILE ${CMAKE_BINARY_DIR}/parser.tab.h)

ADD_FLEX_BISON_DEPENDENCY(Lexer Parser)

set(SOURCE_FILES
        ${FLEX_Lexer_OUTPUTS}
        ${BISON_Parser_OUTPUTS}
        lexer/tokens.h
)

add_executable(Strawberry ${SOURCE_FILES})
include_directories(lexer)

target_include_directories(Strawberry PRIVATE ${CMAKE_SOURCE_DIR}/include)

add_executable(lexer_test
        tests/lexer_test.cpp
        ${FLEX_Lexer_OUTPUTS}
        lexer/tokens.h
)

target_include_directories(lexer_test PRIVATE ${CMAKE_SOURCE_DIR}/include)
target_include_directories(lexer_test PRIVATE ${CMAKE_CURRENT_BINARY_DIR})

enable_testing()
add_test(NAME LexerTest COMMAND lexer_test)
