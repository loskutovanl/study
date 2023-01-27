# HTTP-сервис

Принимает входящие соединения с JSON-данными и обрабатывает их.

1. Обработчик создания пользователя.

```
POST /users/new HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"name":"some name","age":"24","friends":[]}
```
Запрос возвращает ID пользователя в формате json и статус 201.

2. Обработчик, который делает друзей из двух пользователей.

```
POST /users/befriend HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"source_id":"1","target_id":"2"}
```
Запрос возвращает статус 200 и сообщение «username_1 и username_2 теперь друзья».

3. Обработчик, который удаляет пользователя.

```
DELETE /users/delete HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"target_id":"1"}
```
Запрос возвращает 200 и имя удалённого пользователя.

4. Обработчик, который возвращает всех друзей пользователя.

```
GET /users/user_id/friends HTTP/1.1
Host: localhost:8080
Connection: close
```
Запрос возвращает всех друзей пользователя с id = user_id в формате json.

5. Обработчик, который обновляет возраст пользователя.

```
PUT /users/user_id HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"new_age":"28"}
```
Запрос возвращает 200 и сообщение «возраст пользователя успешно обновлён».