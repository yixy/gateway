#!/bin/bash

mkdir data

if [ ! -f "data/10B" ]; then
  dd if=/dev/urandom of=data/10B bs=10 count=1
fi
if [ ! -f "data/100B" ]; then
  dd if=/dev/urandom of=data/100B bs=100 count=1
fi
if [ ! -f "data/1K" ]; then
  dd if=/dev/urandom of=data/1K bs=1024 count=1
fi
if [ ! -f "data/10K" ]; then
  dd if=/dev/urandom of=data/10K bs=10240 count=1
fi
if [ ! -f "data/500K" ]; then
  dd if=/dev/urandom of=data/500K bs=512000 count=1
fi
if [ ! -f "data/1M" ]; then
  dd if=/dev/urandom of=data/1M bs=1048576 count=1
fi
if [ ! -f "data/10M" ]; then
  dd if=/dev/urandom of=data/10M bs=1048576 count=10
fi
if [ ! -f "data/100M" ]; then
  dd if=/dev/urandom of=data/100M bs=1048576 count=100
fi
if [ ! -f "data/500M" ]; then
  dd if=/dev/urandom of=data/500M bs=1048576 count=500
fi
if [ ! -f "data/1G" ]; then
  dd if=/dev/urandom of=data/1G bs=1048576 count=1024
fi

files="100B 100M 10B  10K  10M  1G   1K   1M   500K 500M"

for file in $files
do
  shasum -a 256 data/$file | awk '{printf $1}' > data/$file.hash
	#openssl dgst -sign pri.pem -sha256 -out data/$file.byte data/$file.hash
	#base64 data/$file.byte > data/$file.sign
	go build -o jwt .
  ./jwt sha256 data/$file.hash > data/$file.jwt
done

