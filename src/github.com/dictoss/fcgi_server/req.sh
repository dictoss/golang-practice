#!/bin/bash
cat hello2_request.json | curl -v http://192.168.22.102:8081/gofcgi/json/hello2/ -H 'Content-Type:application/json' --data @-
