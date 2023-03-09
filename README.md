# HTTP-service

Handles incoming JSON

1. Handler that creates user.

```
POST /users/new HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"name":"some name","age":"24","friends":[]}
```
The request returns user ID as JSON and 201 status code.

2. Handler that makes two users friends.

```
POST /users/befriend HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"source_id":"1","target_id":"2"}
```
The request returns 200 status code and message «username_1 и username_2 теперь друзья».

3. Handler that deletes user.

```
DELETE /users/delete HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"target_id":"1"}
```
The request returns 200 status code and the name of deleted user.

4. Handler that gets all friends of the user.

```
GET /users/user_id/friends HTTP/1.1
Host: localhost:8080
Connection: close
```
The request returns JSON of all friends of the user with id equal to user_id.

5. Handler that updates user age.

```
PUT /users/user_id HTTP/1.1
Content-Type: application/json; charset=utf-8
Host: localhost:8080
{"new_age":"28"}
```
The request returns 200 status code and message «возраст пользователя успешно обновлён».
