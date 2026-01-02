#!/bin/bash
cd /app

goose -dir ./database/migrations "$@"
