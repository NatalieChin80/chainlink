#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username postgres -h postgres <<-EOSQL
    CREATE USER chainlink;
    CREATE DATABASE chainlink_test;
    GRANT ALL PRIVILEGES ON DATABASE chainlink_test TO chainlink;
EOSQL
