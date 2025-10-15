#!/usr/bin/env bash

# Getting env
curl --http1.1 -s -H "X-Vault-Token: $vaulttoken" \
    https://vault.HIDDEN.com:8200/v1/kv/data/HIDDEN/bff/$app_env | jq -r .data.data.env >.env

# Checking if env file exist
if ! [ -s /app/.env ]; then
    echo "ERROR: '/app/.env' file does not exist or is empty"
    exit 1
fi

# Starting app
./bff
