# HTTP Stream Store Service

HTTP Stream Store Service -- простой сервис хранения потоков данных c RESTful интерфейсом. Обеспечивает прием, сохранение и возврат потоков данных для передачи их между HTTP-клиентами.

## Особенности реализации и использования
Потоки данных принимаются в теле HTTP-запросов и сохраняются в файлы на сервере. Каждый поток соответствует одному запросу и сохраняется в один файл. Размер потока ограничен лимитами файловой системы сервера. При сохранении потока ему присваивается уникальный идентификатор, который возвращается клиенту.
После успешного сохранения, поток может быть запрошен клиентом при предъявлении идентификатора. После чтения потока он не удаляется. Для удаления потока необходим отдельный запрос. 

Одновременное чтение и запись одного потока невозможны. Для чтения потока, его запись должна быть завершена. Идентификатор потока возвращается только после успешного завершения записи. Возможно параллельное чтение одного потока несколькик клиентами. Также возможна одновременная запись нескольких различных потоков.

Перезапись или дополнение существующего потока не возможны. Запрос на сохранение потока всегда приводит к созданию нового потока с уникальным идентификатором. (в перспективе будет возможно объединение потоков). В случае неуспешной записи (прерывание операции), поток удалется. Также возможны удаление существующего потока и создание нового.

Запрос на удаление потока не приводит к прерыванию текущих операций чтения. 

Запись потоков сихронизируются с файловой системой, поэтому возможна параллельная работа нескольких экземпляров сервиса с одной и той же локальной файловой системой или с распределенной файловой системой в режиме синхронной записи. 

### TO DO (не реализовано)

Потоки организованы в очереди, в каждой очереди ведется список всех потоков, упорядоченный по времени создания потока (FIFO). 
Доступны операции получения списка потоков очереди, получения первого потока в очереди, а такде следующего или предыдущего потока относительно заданного. 

## Интерфейс (v1)

### Параметры командной строки

```sh
$ httpstream -h
HTTP Stream Store Service -- простой сервис хранения потоков данных c RESTful интерфейсом.
Usage of httpstream:
  -addr string
        Service listennig address & port (default "0.0.0.0:8080")
  -d    Debug mode
  -path string
        Streams storage path (default "DATA")
```

### Сохранение нового потока 

Метод POST:

```sh
$curl -d "sоme data" localhost:8080/v1/streams
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> POST /v1/streams HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> Content-Length: 9
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 9 out of 9 bytes
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Fri, 02 Jun 2017 12:27:05 GMT
< Content-Length: 65
< 
{"id":"20ef43f7-4bfd-4b86-9204-7c857a09a940","status":"created"}
```
Создан поток с иднтификатором 1393da65-aabf-48b5-a58e-248bf5b9eac3, содержащий "some data".

Метод PUT работает аналогично (загрузка файла):

```sh
$ curl -v -T "./testfile" localhost:8080/v1/streams
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> PUT /v1/streams HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> Content-Length: 5877
> Expect: 100-continue
> 
< HTTP/1.1 100 Continue
} [5877 bytes data]
* We are completely uploaded and fine
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Date: Fri, 02 Jun 2017 12:28:03 GMT
< Content-Length: 65
< 
{ [65 bytes data]
{"id":"55763224-c834-4c11-8966-1511512b4b4e","status":"created"}
```
Создан поток с иднтификатором 55763224-c834-4c11-8966-1511512b4b4e, содержащий данные файла "testfile".

### Получение данных потока

Метод GET:

```sh
$ curl -v localhost:8080/v1/streams/851db6b3-b94e-49dc-a087-4b7e15d8bfcb
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /v1/streams/851db6b3-b94e-49dc-a087-4b7e15d8bfcb HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Content-Length: 9
< Content-Type: text/plain; charset=utf-8
< Last-Modified: Fri, 02 Jun 2017 12:24:07 GMT
< Date: Fri, 02 Jun 2017 12:32:09 GMT
< 
* Connection #0 to host localhost left intact
same data
``` 

В случае отсутствующего потока возвращается ошибка 404:

```sh
$ curl -v localhost:8080/v1/streams/error-in-stream-id
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET /v1/streams/error-in-stream-id HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> 
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Fri, 02 Jun 2017 12:35:07 GMT
< Content-Length: 19
< 
404 page not found
```

### Удаление потока

Метод DELETE:

```sh
$ curl -X DELETE localhost:8080/v1/streams/851db6b3-b94e-49dc-a087-4b7e15d8bfcb
{"id":"851db6b3-b94e-49dc-a087-4b7e15d8bfcb","status":"deleted"}
```

В случае отсутствущего потока выдается ошибка 404 c пояснением причины ошибки удаления:
```sh
$ curl -v -XDELETE localhost:8080/v1/streams/851db6b3-b94e-49dc-a087-4b7e15d8bfcb
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> DELETE /v1/streams/851db6b3-b94e-49dc-a087-4b7e15d8bfcb HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.47.0
> Accept: */*
> 
< HTTP/1.1 404 Not Found
< Content-Type: application/json; charset=utf-8
< Date: Fri, 02 Jun 2017 12:42:22 GMT
< Content-Length: 107
< 
{"error":{"Op":"remove","Path":"DATA/851db6b3-b94e-49dc-a087-4b7e15d8bfcb.data","Err":2},"status":"error"}
```





