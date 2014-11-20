#!/bin/bash

if [ ! -v OH_URL ]
then
	echo Need OH_URL environment variable
	exit 1
fi

tls=""
cert=""
key=""
if [ -v OH_TLS ]
then
	if [ $OH_TLS == "true" ]
	then
		tls="-tls"
		cert="-cert /etc/automated/cert.pem"
		key="-key /etc/automated/key.pem"
	fi
fi

/go/bin/automated -ohurl $OH_URL $tls $cert $key
