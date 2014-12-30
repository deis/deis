#!/bin/bash

# generate management cert

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -config cert.conf -keyout azure-cert.pem -out azure-cert.pem -config cert.conf
openssl  x509 -outform der -in azure-cert.pem -out azure-cert.cer
