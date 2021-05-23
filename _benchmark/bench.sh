#!/bin/bash

echo "Dynamic route test"
wrk -c100 -d10 -t10 http://localhost:8080/some/page/xjj

echo
echo "Static route test"
wrk -c100 -d10 -t10 http://localhost:8080/other/page/path
