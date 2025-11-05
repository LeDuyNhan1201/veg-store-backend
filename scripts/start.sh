#!/bin/bash
set -euo pipefail

export CA_COMMON_NAME="LDNhanCA"
export SUBJ_C="VN"
export SUBJ_ST="5"
export SUBJ_L="HCM"
export SUBJ_O="SGU"
export SUBJ_OU="Dev"

export CERTS_DIR="certs"

mode=${1:-"dev"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

env_file="${DIR}/helper/env_config.sh"
echo "Processing $env_file"

source "$env_file"
source "${DIR}/helper/functions.sh"
source "${DIR}/helper/generate_certs.sh"

create_env_file
generate_root_ca
generate_keystore_and_truststore "go-clean-arch" "go-clean-arch"

docker build -f docker/Dockerfile.builder -t learn-go/go-clean-arch-builder .

docker compose -f docker/docker-compose."${mode}".yml up --force-recreate -d