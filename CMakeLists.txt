cmake_minimum_required(VERSION 3.10)
project(Strawberry)

set(CMAKE_CXX_STANDARD 17)

find_package(FLEX REQUIRED)
find_package(BISON REQUIRED)

FLEX_TARGET(Lexer src/lexer/lexer.l ${CMAKE_CURRENT_BINARY_DIR}/lexer.cpp)
BISON_TARGET(Parser src/parser/parser.y ${CMAKE_CURRENT_BINARY_DIR}/parser.cpp)

ADD_FLEX_BISON_DEPENDENCY(Lexer Parser)

add_executable(Strawberry
        ${FLEX_Lexer_OUTPUTS}
        ${BISON_Parser_OUTPUTS}
        src/lexer/lexer.cpp
        src/lexer/lexer.h
        src/lexer/tokens.h
        #src/main.cpp
        #src/llvm_codegen/codegen.cpp
        #src/utils/logger.cpp
)

# Подключение библиотек LLVM
#find_package(LLVM REQUIRED CONFIG)
#include_directories(${LLVM_INCLUDE_DIRS})
#add_definitions(${LLVM_DEFINITIONS})
#target_link_libraries(my_language LLVM)
