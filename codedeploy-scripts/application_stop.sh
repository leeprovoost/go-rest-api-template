#!/bin/bash

# stop application
kill -9 `cat /var/log/go-rest-api-template-pid.txt`

# clean up old files
rm /var/log/go-rest-api-template-pid.txt
rm /var/log/go-rest-api-template.log
rm -rf cd $GOPATH/src/github.com/leeprovoost/go-rest-api-template
rm /home/ec2-user/go/bin/go-rest-api-template
rm /home/ec2-user/go/bin/fixtures.json
