#!/bin/bash

# stop application
sudo service go-rest-api-template stop

# clean up old files
rm -rf cd $GOPATH/src/github.com/leeprovoost/go-rest-api-template
rm /home/ec2-user/go/bin/go-rest-api-template
rm /home/ec2-user/go/bin/fixtures.json
