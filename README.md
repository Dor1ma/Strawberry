# Strawberry
Programming language for Compliers and Platforms subject at ITMO University

## Оглавление
- [Примеры синтаксиса](#примеры-синтаксиса)
    - [Объявление переменных](#объявление-переменных)
    - [Условные операторы](#условные-операторы)
    - [Циклы](#циклы)
    - [Функции](#функции)
- [Гайд по разработке](#гайд-по-разработке)
    - [Chocolatey](#chocolatey)
    - [Flex и Bison](#flex-и-bison)
    - [Структура проекта](#структура-проекта)
    - [CMakeLists.txt](#cmakelists-txt)
    - [Git](#git)
- [Support](#support)

## Примеры синтаксиса

### Объявление переменных

```plaintext
berry x = 10
berry name = "Strawberry"
berry age int = 25
berry b = 10 % 2
berry l list = [1, 2, 3, 4, 5]
```

### Условные операторы
```plaintext
if x > 5 {
    say "X is greater than 5"
} elif x == 5 {
    say "X is equal to 5"
} else {
    say "X is less than 5"
}
```

### Циклы
```plaintext
for i in range(1, 11) {
    say i
}
```

### Функции
```plaintext
recipe divide(a: float, b: float) -> result: float {
    if b == 0 {
        return 0
    }
    return a / b
}

berry result = divide(10.0, 0.0)

recipe greet(name: string) -> void {
    say "Hello, " + name
}

greet("World")
```

## Гайд по разработке

Для того, чтобы была возможность собрать проект, необходимо установить
следующие зависимости:

- chocolatey
- win_flex
- win_bison

Короткая инструкция по их установке

### Chocolatey
Chocolatey - это пакетный менеджер для Windows. Будем использовать его, т.к. он сильно
упрощает установку используемых библиотек. Как установить:

1. Запускаем Powershell от имени администратора
2. Запускаем команду
    ```plaintext
    Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
    ```
3. Проверяем, что chocolatey успешно установился:
    ```plaintext
    choco --version
    ```

### Flex и Bison

1. С помощью choco мы можем поставить flex и bison одной командой:
    ```plaintext
    choco install winflexbison
    ```
2. Проверим, что все корректно установилось:
    ```plaintext
    win_flex --version
    win_bison --version
    ```

### Структура проекта
При разработке прошу придерживаться следующей структуры:
```plaintext
Strawberry/
│
├── src/                  # Основной исходный код проекта
│   ├── lexer/            # Модуль лексера
│   │   ├── lexer.l
│   │   └── tokens.h
│   │
│   │
│   ├── parser/           # Модуль парсера
│   │   ├── parser.y
│   │   └── остальные файлы
│   │
│   │
│   ├── llvm_codegen/     # Модуль генерации кода с использованием LLVM
│   │   └── файлы
│   │
│   │
│   └── main.cpp          # Точка входа для компилятора/интерпретатора (добавим позднее)
│
│
├── tests/                # Тесты для компонентов (юнит-тесты, интеграционные тесты)
│   ├── lexer_tests.cpp
│   ├── parser_tests.cpp
│   └── codegen_tests.cpp
│
├── cmake-build-debug/    # Директория для сборки (не хранится в системе контроля версий)
│
└── CMakeLists.txt        # Файл сборки проекта с CMake или другой системой сборки

```

### CMakeLists.txt
CMakeLists.txt это файл для подключения зависимостей и настройки
исходников для сборки.

На момент готовности лексера в этом файле **уже** подключены Flex и Bison. LLVM
нужно будет подключать отдельно.

**Если вам нужно подключить или добавить что-нибудь новое, по возможности не удаляйте предыдущие настройки конфига!**

### Git

1. Бранчимся от main. Ветку именуем в виде <тип работы в ветке>/краткое_описание. Пример:
    ```plaintext
    feature/lexer
    ```
2. Обязательно руками прогоняем тесты
3. Если тесты прошли успешно, кидаем PR. **Самостоятельно не мерджим!**


## Support
Если вы понимаете, что в другие части проекта нужно добавить что-то, без чего
ваша часть работать не будет или вам нужно больше информации насчет конкретного компонента -
тегайте в чате