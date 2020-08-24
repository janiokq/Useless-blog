#!/bin/bash
set -u

#项目相关的
ProjectName=${ProjectName:-"https://github.com/janiokq/Useless-blog"}
Version=${Version:-"unknow"}
TARGET=${TARGET:-'main'}

#执行环境
GOPROXY=${GOPROXY:-"https://goproxy.cn"}
#go mod是否开启
GO111MODULE=${GO111MODULE:-"auto"}
#GOPATH的路径
GOPATH=${GOPATH:-${HOME}"/go"}
#其他软件的安装目录
soft_dir=${soft_dir:-${HOME}}
#go安装的版本
go_version=${go_version:-"1.15"}
#protoc的版本
protoc_version=${protoc_version:-"3.12.4"}
#protoc引用的路
protoc_include_path=${protoc_include_path:-"/usr/local/Cellar/protobuf/${protoc_version}/include"}
#cloc版本
cloc_version=${cloc_version:-"1.86"}
#执行文件路径
cmd_path=${cmd_path:-"${GOPATH}/bin"}
grpc_gateway_path="${GOPATH}/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.14.7"
mkdir -p ${GOPATH}/bin
mkdir -p ${GOPATH}/src
#将环境变量存入本地环境配置
echo "GOPROXY=${GOPROXY}" >>${HOME}/.profile
echo "protoc_include_path=${protoc_include_path}" >>${HOME}/.profile
echo "GO111MODULE=${GO111MODULE}" >>${HOME}/.profile
echo "GOPATH=${GOPATH}" >>${HOME}/.profile
echo "PATH=${soft_dir}/go/bin:${GOPATH}/bin:${PATH}" >>${HOME}/.profile

#手动执行
#source ~/.profile
