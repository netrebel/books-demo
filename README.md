# books-demo
simple JSON Rest API demo

# How to run with docker-compose

* `docker-compose build app` to rebuild the app with code changes
* `docker-compose up` starts the app with redis
* `docker restart books-demo_prometheus_1` restarts prometheus
* `docker exec -it books-demo_prometheus_1 sh` enter running container
* `docker logs -f books-demo_prometheus_1` tail logs

Note: Redis will cache the value even if the container is stopped.

# URLs

* App: http://localhost:9000
* Prometheus: http://localhost:9090
* Grafana: http://localhost:3000

## Endpoints

POST /books

```sh
curl --location --request POST 'http://localhost:9000/books' \
--header 'Content-Type: application/json' \
--data-raw '{
    "isbn": "1234",
    "title": "Sky",
    "author": {
        "first_name" : "John",
        "last_name" : "Doe"
    }
}'
```

GET /books

```sh
curl --location --request GET 'http://localhost:9000/books'
```

GET /books/{id}
```sh
curl --location --request GET 'http://localhost:9000/books/{id}'
```

DELETE /books/{id}
```sh
curl --location --request DELETE 'http://localhost:9000/books/{id}'
```

PUT /books/{id}
```sh
curl --location --request PUT 'http://localhost:9000/books/{id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "isbn": "1234",
    "title": "Sky",
    "author": {
        "first_name" : "Miguel",
        "last_name" : "Reyes"
    }
}'
```

GET prometheus/

```sh
curl --location --request GET 'http://localhost:9000/prometheus'
```

# Load tests

```
hey -z 5m -q 5 -m GET -H "Accept: text/html" http://127.0.0.1:9000
```

Graphana dashboard: https://gist.github.com/netrebel/b64c4af742c11d61a0fdf0763979515a

# Refs:

* Monitoring implementation. See https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang