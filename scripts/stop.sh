#!/bin/bash
set -euo pipefail

mode=${1:-"dev"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
env_file="${DIR}/helper/env_config.sh"
source "$env_file"

docker compose -f docker/docker-compose."${mode}".yml down -v

echo "Removing certs..."

sudo rm -rf certs/*
sudo rm -rf secrets/*
sudo rm -rf docker/.env

docker rmi veg-store/backend-builder:"${BACKEND_VERSION}"
docker rmi veg-store/postgres-full-ssl

echo "Done."