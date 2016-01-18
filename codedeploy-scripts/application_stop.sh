#!/bin/bash

# stop application
kill -9 `cat /var/log/go-rest-api-template-pid.txt`

# clean up old files
rm -f /var/log/go-rest-api-template-pid.txt
rm -f /var/log/go-rest-api-template.log
rm -rf $GOPATH/src/github.com/leeprovoost/go-rest-api-template
rm -f /usr/bin/go-rest-api-template
rm -f /home/ec2-user/fixtures.json
