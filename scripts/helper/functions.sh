#!/bin/bash

create_env_file() {
    # Clean old content before overwrite
    : > docker/.env

    echo CERT_SECRET=$CERT_SECRET >> docker/.env

    echo POSTGRES_USER=$POSTGRES_USER >> docker/.env
    echo POSTGRES_PASSWORD=$POSTGRES_PASSWORD >> docker/.env

    echo GIN_MODE=$GIN_MODE >> docker/.env
    echo MODE=$MODE >> docker/.env

}

