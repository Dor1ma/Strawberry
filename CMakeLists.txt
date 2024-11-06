cmake_minimum_required(VERSION 3.10)
project(Strawberry)

set(CMAKE_CXX_STANDARD 17)

find_package(FLEX REQUIRED)
find_package(BISON REQUIRED)

FLEX_TARGET(Lexer src/lexer/lexer.l ${CMAKE_CURRENT_BINARY_DIR}/lexer.cpp)
#BISON_TARGET(Parser src/parser/parser.y ${CMAKE_CURRENT_BINARY_DIR}/parser.cpp)

# Устанавливаем зависимость, чтобы лексер генерировался до парсера
#ADD_FLEX_BISON_DEPENDENCY(Lexer Parser)

set(SOURCE_FILES
        ${FLEX_Lexer_OUTPUTS}
 #       ${BISON_Parser_OUTPUTS}
        src/lexer/lexer.cpp
        src/lexer/lexer.h
        src/lexer/tokens.h
        src/include/tokens.h
        src/lexer/lexer.cpp
)

add_executable(Strawberry ${SOURCE_FILES})
include_directories(src/lexer)

target_include_directories(Strawberry PRIVATE ${CMAKE_SOURCE_DIR}/include)

add_executable(lexer_test
        tests/lexer_test.cpp
        ${FLEX_Lexer_OUTPUTS}
        src/lexer/lexer.h
        src/lexer/tokens.h
        src/lexer/lexer.cpp
)

target_include_directories(lexer_test PRIVATE ${CMAKE_SOURCE_DIR}/include)

enable_testing()
add_test(NAME LexerTest COMMAND lexer_test)
