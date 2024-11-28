# external/CMakeLists.txt
# Внешний файл для работы с лексером и парсером
add_executable(ExternalLexer ${CMAKE_SOURCE_DIR}/tests/lexer/lexer_test.cpp)
target_link_libraries(ExternalLexer PRIVATE LexerLib ParserLib)

add_executable(ExternalParser ${CMAKE_SOURCE_DIR}/tests/parser/parser_test.cpp)
target_link_libraries(ExternalParser PRIVATE ParserLib LexerLib)

# Заголовочные файлы в качестве доступных внешним приложениям
target_include_directories(ExternalLexer PRIVATE ${CMAKE_SOURCE_DIR}/src/lexer ${CMAKE_CURRENT_BINARY_DIR})
target_include_directories(ExternalParser PRIVATE ${CMAKE_SOURCE_DIR}/src/parser ${CMAKE_CURRENT_BINARY_DIR})
