#!/usr/bin/env bash
docker compose run app migrate create -dir sql -ext sql $1