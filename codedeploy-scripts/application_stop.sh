#!/bin/bash

if [ -d $HOME/go/src/github.com/leeprovoost/go-rest-api-template]; then
  rm -rf $HOME/go/src/github.com/leeprovoost/go-rest-api-template
fi
if [ -f /var/log/go-rest-api-template-pid.txt ]; then
  kill -9 `cat /var/log/go-rest-api-template-pid.txt`
  rm -f /var/log/go-rest-api-template-pid.txt
fi
if [ -f /var/log/go-rest-api-template.log ]; then
  rm -f /var/log/go-rest-api-template.log
fi
if [ -d /opt/go-rest-api-template ]; then
  rm -rf /opt/go-rest-api-template
fi
