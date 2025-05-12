# Strawberry
Programming language for Compliers and Platforms subject at ITMO University

## Оглавление
- [Примеры синтаксиса](#примеры-синтаксиса)
    - [Объявление переменных](#объявление-переменных)
    - [Условные операторы](#условные-операторы)
    - [Циклы](#циклы)
    - [Функции](#функции)
- [Useful info](#useful-info)
    - [Git](#git)
- [Support](#support)

## Примеры синтаксиса

### Объявление переменных

```plaintext
var a = 0;
var b = "Strawberry";
var c = 10 * 2;
var d = [1, 2, 3, 4, 5];
```

### Условные операторы
```plaintext
if (x > 5) {
    return "X is greater than 5";
}
else {
    return "X is less than 5";
}
```

### Циклы
```plaintext
var arr = [3, 2, 1, 5, -10, 12];
var n = 6;
for (var i = 0; i < n; i = i + 1) {
    for (var j = 0; j < n - i - 1; j = j + 1) {
        if (arr[j] > arr[j + 1]) {
            var temp = arr[j];
            arr[j] = arr[j + 1];
            arr[j + 1] = temp;
        }
    }
}
```

### Функции
```plaintext
fun a() {
    var v = 1;
    v = v * 2 + 3;
    return v + 4;
}

var b = a();
print b;
```

## Useful info

### Git

1. Бранчимся от main. Ветку именуем в виде <тип работы в ветке>/краткое_описание. Пример:
    ```plaintext
    feature/lexer
    ```
2. Обязательно прогоняем имеющиеся unit-тесты
3. Если тесты прошли успешно, кидаем PR. **Самостоятельно не мерджим!**
4. 


## Support
Если вы понимаете, что в другие части проекта нужно добавить что-то, без чего
ваша часть работать не будет или вам нужно больше информации насчет конкретного компонента -
тегайте в чате
