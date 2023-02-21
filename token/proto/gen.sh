#!/usr/bin/env bash

#================================================================
#   Copyright (C) 2021. All rights reserved.
#   
#   file : gen.sh
#   coder: zemanzeng
#   time : 2021-09-06 10:11:18
#   desc : 协议生成脚本
#
#================================================================


protoc -I=./ --go_out=./ business_token.proto
cp -r business_token/proto/business_token.pb.go . && \
    rm -rf business_token
