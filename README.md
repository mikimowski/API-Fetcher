HTTP Server providing REST API for CRUD operations.
API allows to (un)subscribe to given urls and retrive saved data. Data is fetched from subscribed URLs on given interval.
The key component is [subscriber](https://github.com/mikimowski/API-Fetcher/tree/master/subscriber) module that manages subscriptions. Neat implementation of timers is leveraged so that goroutines are spawned only when needed and no conflicts occur

### Database
Two options:
* memory database
* mongodb database

### Running in docker
From main directory:<br>
`docker-compose up -d --build`<br>
This will use mongodb
* mongodb listens on :27017
* server listens on :8080

#### Clean up
* `docker-compose down -v --rmi all --remove-orphans`

### Running locally

#### With memory database
* no setup required, install necessary packages and run
* supports one program argument `debug`. This controls logging. Debug mode is more verbose.
    `/go/bin/TWFjaWVqLU1pa3XFgmE debug`

#### With MongoDB
* setup mongoDB
* set ENV variable `MONGO_URL`. That is mongodb connection URL
  for instance, `MONGO_URL=mongodb://localhost:27017`

#### Tests
Some basic and naive testing is implemented.

#### API examples
* add subscription `curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url": "https://httpbin.org/range/15","interval":60}'`
* list subscriptions `curl -s 127.0.0.1:8080/api/fetcher`
* list history `curl -si 127.0.0.1:8080/api/fetcher/1/history`
* delete subscription `curl -s 127.0.0.1:8080/api/fetcher/12 -X DELETE`
* update subscription<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"interval":6}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"url": "https://httpbin.org/range/10"}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"url": "https://httpbin.org/range/10", "interval":6}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"id":1, "url": "https://httpbin.org/range/10", "interval":6}'`<br>
    * invalid: `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"id":42, "url": "https://httpbin.org/range/10", "interval":6}'`<br>
* *id* unique integer > 0
* *url* valid url string
* *interval* integer > 0
