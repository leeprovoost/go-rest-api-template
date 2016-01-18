#!/bin/bash

# stop application
kill -9 `cat /var/log/go-rest-api-template-pid.txt`

# clean up old files
rm -f /var/log/go-rest-api-template-pid.txt
rm -f /var/log/go-rest-api-template.log
rm -f /opt/go-rest-api-template/go-rest-api-template
rm -f /opt/go-rest-api-template/fixtures.json
