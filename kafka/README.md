## Usage
Basically I'm just trying to make Kafka work...

## Background
There's three team in different countries, a data team, a produce team, a process team.
They have their own services and needs to deal with the same data, services could go offline at anytime, and we couldn't afford the consequences of losing data.
So we have to build **high available** microservices with Apache Kafka.

### 1. Build
```shell
docker-compose up -d
```

### 2. Start producer in terminal 1
```shell
go run main.go
```

### 3. Start service in terminal 2
```shell
go run processor/main.go
```

### 4. Start another service in terminal 3
```shell
go run datateam/main.go
```

If everything works fine, producer keeps produce messages into the queue, and both service will consume messages from the **same** queue.

Now shutdown one of the two services in terminal 2 & 3,  wait a few seconds, then reboot it. It will consume messages from the latest position when it's offline(high available).
