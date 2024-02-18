#!/bin/sh

#* variables
PROTO_PATH=./idl/proto
PROTO_OUT=./dto
IDL_PATH=./idl
DOC_OUT=./docs

mkdir -p ${DOC_OUT}/html
mkdir -p ${DOC_OUT}/markdown

for folder in $PROTO_PATH/*; do
	doc_out=${DOC_OUT}/${folder}
	protoc \
		${folder}/*.proto \
		-I=/usr/local/include \
		--proto_path=${folder} \
		--go_out=:${PROTO_OUT} \
		--validate_out=lang=go:${folder} \
		--go-grpc_out=:${PROTO_OUT} \
		--grpc-gateway_out=:${PROTO_OUT} \
		--openapiv2_out=:${doc_out}/swagger \
		--custom_out=:${PROTO_OUT} \
		--doc_out=:${doc_out}/html --doc_opt=html,index.html
done
