#!/bin/bash

if [ -d $HOME/go/src/github.com/leeprovoost/go-rest-api-template]; then
  rm -rf $HOME/go/src/github.com/leeprovoost/go-rest-api-template
fi

if [ -f /var/log/go-rest-api-template-pid.txt ]; then
  # Stop application
  kill -9 `cat /var/log/go-rest-api-template-pid.txt`
  # Then delete all files
  rm -f /var/log/go-rest-api-template-pid.txt
fi
if [ -f /var/log/go-rest-api-template.log ]; then
  rm -f /var/log/go-rest-api-template.log
fi
if [ -f /opt/go-rest-api-template/go-rest-api-template ]; then
  rm -f /opt/go-rest-api-template/go-rest-api-template
fi
if [ -f /opt/go-rest-api-template/fixtures.json ]; then
  rm -f /opt/go-rest-api-template/fixtures.json
fi
