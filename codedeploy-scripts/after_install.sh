#!/bin/bash

# create app directory
sudo mkdir /opt/go-rest-api-template

# copy application binary
sudo cp $HOME/go/src/github.com/leeprovoost/go-rest-api-template/go-rest-api-template /opt/go-rest-api-template

# copy fixtures,json file
sudo cp $HOME/go/src/github.com/leeprovoost/go-rest-api-template/fixtures.json /opt/go-rest-api-template
