#!/bin/bash
set -euo pipefail

mode=${1:-"dev"}

service_name=${2:-"veg-store-backend"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

env_file="${DIR}/helper/env_config.sh"
echo "Processing $env_file"

source "$env_file"
source "${DIR}/helper/functions.sh"

create_env_file

docker compose -f docker/docker-compose."${mode}".yml up "${service_name}" --force-recreate -d