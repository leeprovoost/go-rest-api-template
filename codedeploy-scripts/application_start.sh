#!/bin/bash

sudo nohup /home/ec2-user/go/bin/go-rest-api-template -fixtures=/home/ec2-user/go/bin/fixtures.json -port=80 > /var/log/go-rest-api-template.log 2>&1&
echo $! > /var/log/go-rest-api-template-pid.txt
