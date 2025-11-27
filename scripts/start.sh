#!/bin/bash
set -euo pipefail

export CA_COMMON_NAME="LDNhanCA"
export SUBJ_C="VN"
export SUBJ_ST="5"
export SUBJ_L="HCM"
export SUBJ_O="SGU"
export SUBJ_OU="Dev"
export SECRETS_DIR="secrets"
export CERTS_DIR="${SECRETS_DIR}/certs"

mode=${1:-"dev"}

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

# Load environment variables
env_file="${DIR}/helper/env_config.sh"
echo "Processing $env_file"

# Load all needed function
source "$env_file"
source "${DIR}/helper/functions.sh"
source "${DIR}/helper/generate_certs.sh"

# Create .env
create_env_file

# Generate certificates for HTTPS
generate_root_ca
generate_cert_with_keystore_and_truststore "veg-store-backend" "veg-store-backend"
generate_cert_with_keystore_and_truststore "postgres" "postgres"

# Generate asymmetric key for JWT
generate_keypair

# Create custom images on Docker local
docker build -f docker/veg-store-backend/Dockerfile.builder -t veg-store/backend-builder:"${BACKEND_VERSION}" .
docker build -f docker/postgres/Dockerfile -t veg-store/postgres-full-ssl .

# Compose up all services with ${mode}
docker compose -f docker/docker-compose."${mode}".yml up --force-recreate -d