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

go build -o jwt .

rm data/client_script
rm data/server_script

for file in $files
do
  shasum -a 512 data/$file | awk '{printf $1}' > data/$file.hash
	#openssl dgst -sign pri.pem -sha256 -out data/$file.byte data/$file.hash
	#base64 data/$file.byte > data/$file.sign
  ./jwt SHA512 data/$file.hash > data/$file.jwt

  echo 'nohup hey -n 100000000 -c 100 -z 10s -m POST -H "Gateway-Jwt: '`cat data/$file.jwt`'" -D 'data/$file' http://10.168.0.5:9090/api/xxx/yyy/v1 &' >> data/client_script

  echo 'nohup gudong start -H='"'"'Gateway-Apiresp:{"hashed":"'`cat data/$file.hash`'","resp":{"return_code":0,"return_msg":"success"}}'"'"' --body-file=data/'$file' &' >> data/server_script
done

