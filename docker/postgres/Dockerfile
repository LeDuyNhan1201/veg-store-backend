FROM postgres:latest

COPY secrets/certs/postgres/server.key /var/lib/postgresql/18/main/server.key
COPY secrets/certs/postgres/server.crt /var/lib/postgresql/18/main/server.crt
COPY secrets/certs/ca/ca.crt /var/lib/postgresql/18/main/ca.crt

RUN chmod 600 /var/lib/postgresql/18/main/server.key
RUN chown postgres:postgres /var/lib/postgresql/18/main/server.key