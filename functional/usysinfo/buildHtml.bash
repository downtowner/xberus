#!/bin/bash

# 该脚本用于生成html
# 前端仓库http://git.vnnox.net/loveRandy/dashboard

go get github.com/rakyll/statik

statik -src=./html
