我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

HTTP Server providing REST API for CRUD operations.<br>
API allows to (un)subscribe to given urls and retrive saved data. Data is fetched from subscribed URLs on given interval.<br>
The key component is [subscriber](https://github.com/mikimowski/API-Fetcher/tree/master/subscriber) module that manages subscriptions. Neat implementation of timers is leveraged so that goroutines are spawned only when needed and no conflicts occur.

### Database
Two options:
* [memory database](https://github.com/mikimowski/API-Fetcher/blob/master/data/memory_database.go)
* [mongodb database](https://github.com/mikimowski/API-Fetcher/blob/master/data/dao_mongo.go)

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
    * `[{"id":1,"url":"https://httpbin.org/range/15","interval":60}, {"id":2,"url": "https://httpbin.org/delay/10","interval":120}]`
* list history `curl -si 127.0.0.1:8080/api/fetcher/1/history`
    * `HTTP/1.1 200 OK [{"response": "abcdefghijklmno", "duration": 0.571, "created_at": 1559034638.31525,}, {"response": null, "duration": 5,"created_at": 1559034938.623,}, ]`
* delete subscription `curl -s 127.0.0.1:8080/api/fetcher/12 -X DELETE`
    * `curl -s 127.0.0.1:8080/api/fetcher/12 -X DELETE $ curl -s 127.0.0.1:8080/api/fetcher/12/history -i`
    * `HTTP/1.1 404 Not Found`
* update subscription<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"interval":6}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"url": "https://httpbin.org/range/10"}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"url": "https://httpbin.org/range/10", "interval":6}'`<br>
    * `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"id":1, "url": "https://httpbin.org/range/10", "interval":6}'`<br>
    * invalid: `curl -si 127.0.0.1:8080/api/fetcher/1 -X PATCH -d '{"id":42, "url": "https://httpbin.org/range/10", "interval":6}'`<br>

#### Additional assumptions
* *id* unique integer > 0
* *url* valid url string
* *interval* integer > 0
* POST payload is limited to 1MB
    * curl -si 127.0.0.1:8080/api/fetcher -X POST -d 'more than 1MB of data' 
    * HTTP/1.1 413 Request Entity Too Large
* fetching data from url has timeout set to 5s
