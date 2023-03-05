#!/usr/bin/bash

## this is a temporary solution... I need to organize this into a make command

set -o allexport && source .env && go run main.go
