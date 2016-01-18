#!/bin/bash

sudo nohup /opt/go-rest-api-template/go-rest-api-template -fixtures=/opt/go-rest-api-template/fixtures.json -port=80 > /var/log/go-rest-api-template.log 2>&1&
echo $! > /var/log/go-rest-api-template-pid.txt
