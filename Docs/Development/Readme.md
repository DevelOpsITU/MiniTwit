# Development


To start the environment locally make a copy of the [.envExample](../../.envExample), call it ".env" and adjust the
environment variables accordingly.

Then from the root run 'docker-compose up -d'

and verify that everything is up with the command 'docker-compose ps', it should give the following output:

```shell
$ docker-compose ps
    Name                  Command               State                 Ports               
------------------------------------------------------------------------------------------
data-scraper   python3 -u ./main.py             Up
fluentd        tini -- /bin/entrypoint.sh ...   Up      0.0.0.0:24224->24224/tcp, 5140/tcp
grafana        /run.sh                          Up      0.0.0.0:3000->3000/tcp
loki           /usr/bin/loki -config.file ...   Up      0.0.0.0:3100->3100/tcp
minitwit       ./minitwit                       Up      0.0.0.0:8080->8080/tcp
nginx          nginx -g daemon off;             Up      0.0.0.0:80->80/tcp
postgres       docker-entrypoint.sh postgres    Up      0.0.0.0:5433->5432/tcp
prometheus     /bin/prometheus --config.f ...   Up      0.0.0.0:9090->9090/tcp
```

Now will the application be accessible from [localhost:80](localhost:80) through Nginx and 
[localhost:8080](localhost:8080) for direct access. 

Postgres runs on [localhost:5433](localhost:5433) and internally in docker on port 5432.

Grafana runs on [localhost:3000](localhost:3000)

prometheus runs  on [localhost:9090](localhost:9090)



## Use of SQLite
To start the application with an in memory sqlite database export this environment variable

````Shell
export DB_CONNECTION_STRING=file::memory:
````