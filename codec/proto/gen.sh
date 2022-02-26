#!/usr/bin/env bash

#================================================================
#   Copyright (C) 2022. All rights reserved.
#   
#   file : gen.sh
#   coder: zemanzeng
#   time : 2022-02-13 11:15:17
#   desc : 协议生成脚本
#
#================================================================

protoc -I=./ --go_out=./ service.proto
cp -r codec/proto/*.go . && \
    rm -rf codec
