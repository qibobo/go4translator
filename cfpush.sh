#!/bin/bash
set -e

cf delete -f go4translator
cf ds -f mytl
make build
cf cs language_translator standard mytl
cf push go4translator -b binary_buildpack -c './go4translator' --no-start
# cf set-env go4translator TRANSLATOR_NAME mytl
cf bs go4translator mytl
cf start go4translator
url="https://$(cf app go4translator | grep --color=no "routes:" | awk -F ':            ' '{print $2}')?model=en-zh&text=we%20are%20good"
curl "$url"