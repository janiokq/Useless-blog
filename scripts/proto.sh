#!/bin/bash
set -eu
source scripts/variables.sh

proto(){
    Parameter1=$1
    echo "开始编译proto文件${Parameter1}"
    echo "protoc include 路径为:${protoc_include_path}"
    dirname=./internal/proto/$2/$Parameter1
    swagger_dir=./deployments/config/swagger
    if [ ! -d $swagger_dir ];then
        mkdir -p $swagger_dir
    fi
    if [ -d $dirname ];then
		for f in $dirname/*.proto; do \
		    if [ -f $f ];then \
		        protoc -I. \
                -I$grpc_gateway_path/third_party/googleapis \
                -I$grpc_gateway_path \
                --proto_path=${protoc_include_path} \
                --grpc-gateway_out=. \
                --swagger_out=$swagger_dir \
                --swagger_out=. \
                --go_out=plugins=grpc:. $f; \
                echo compiled protoc: $f; \
            fi \
		done \
	fi
}


proto_inject() {
    echo "开始注入"
    dirname=./internal/proto/$2/$1
    if [ -d $dirname ];then
		for f in $dirname/*.pb.go; do \
		    if [ -f $f ];then \
                protoc-go-inject-tag -input=./$f; \
                echo "完成注入" protoc-go-inject-tag: $f; \
            fi \
		done \
	fi
}


# 用户
proto user v1
proto_inject user v1