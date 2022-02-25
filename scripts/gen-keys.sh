#!/bin/bash

IP="192.168.1.99"
SUBJECT_CA="/C=US/ST=Georgia/L=Atlanta/O=water_sensor/OU=CA/CN=$IP"
SUBJECT_SERVER="/C=US/ST=Georgia/L=Atlanta/O=water_sensor/OU=Server/CN=$IP"
SUBJECT_CLIENT="/C=US/ST=Georgia/L=Atlanta/O=water_sensor/OU=Client/CN=$IP"

SKIP_CA=0
SKIP_SERVER=0
SKIP_CLIENT=0

function usage () {
  echo "Usage: gen-keys [ --skip_ca ] [ --skip_server ] [ --skip_client ]"
  exit 2
}

function generate_CA () {
  echo "$SUBJECT_CA"
  openssl req -x509 -nodes -sha256 -newkey rsa:2048 -subj "$SUBJECT_CA"  -days 365 -keyout ca.key -out ca.crt
}

function generate_server () {
  echo "$SUBJECT_SERVER"
  openssl req -nodes -sha256 -new -subj "$SUBJECT_SERVER" -keyout server.key -out server.csr
  openssl x509 -req -sha256 -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365
}

function generate_client () {
  echo "$SUBJECT_CLIENT"
  openssl req -new -nodes -sha256 -subj "$SUBJECT_CLIENT" -out client.csr -keyout client.key 
  openssl x509 -req -sha256 -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365
}

LONG=skip_ca,skip_server,skip_client

PARSED_ARGUMENTS=$(getopt -n gen-keys -o '' -l $LONG -- "$@")
VALID_ARGUMENTS=$?
if [ "$VALID_ARGUMENTS" != "0" ]; then
  usage
fi

eval set -- "$PARSED_ARGUMENTS"
while :
do
  case "$1" in
    --skip_ca)     SKIP_CA=1     ; shift ;;
    --skip_server) SKIP_SERVER=1 ; shift ;;
    --skip_client) SKIP_CLIENT=1 ; shift ;;
    --) shift ; break ;;
    *) echo "Unexpected option: $1 - how'd this happen 😱."
      usage ;;
  esac
done

if [ $SKIP_CA -ne 1 ]; then
  generate_CA
fi

if [ $SKIP_SERVER -ne 1 ]; then
  generate_server
fi

if [ $SKIP_CLIENT -ne 1 ]; then
  generate_client
fi