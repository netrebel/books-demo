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

* Prometheus: http://localhost:9090