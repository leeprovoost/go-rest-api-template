#!/bin/bash

sudo nohup /home/ec2-user/go/bin/go-rest-api-template > /var/log/go-rest-api-template.log 2>&1&
echo $! > /var/log/go-rest-api-template-pid.txt
