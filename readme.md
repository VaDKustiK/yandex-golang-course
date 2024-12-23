# Калькулятор

Я сделал веб-сервис, который выполняет вычисление арифметических выражений. Он предоставляет API-эндпоинт, куда пользователь может отправить арифметическое выражение и получить результат.

## API Эндпоинт

### `POST /api/v1/calculate`

Этот эндпоинт принимает арифметическое выражение и возвращает результат или сообщение об ошибке.

#### Тело запроса
```
{
  "expression": "2+2*2"
}
```
expression: Строка, содержащая арифметическое выражение для вычисления. Может включать цифры, арифметические операторы (+, -, *, /) и скобки.

#### Тело ответа (успешное вычисление)
```
{
  "result": "6"
}
```

#### Тело ответа (ошибка 422 - неверное выражение)
```
{
  "error": "Expression is not valid"
}
```
Эта ошибка возникает, если выражение содержит недопустимые символы или синтаксические ошибки.

#### Тело ответа (ошибка 500 - внутренняя ошибка сервера)
```
{
  "error": "Internal server error"
}
```

## Пример использования

Для тестирования API можно использовать curl:

#### Успешное вычисление
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

#### Выражение неверно (ошибка 422)
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2a*2"
}'
```

#### Внутренняя ошибка сервера (ошибка 500)
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2/0"
}'
```

## Как использовать программу

1. Клонируй репозиторий
```
git clone https://github.com/VaDKustiK/yandex-golang-course.git
cd yandex-golang-course/calculator_service
```

2. Установи все зависимости
```
go mod init github.com/VaDKustiK/yandex-golang-course
```

3. Запускай
```
go run .
```
Это запустит сервер на порте 8080.