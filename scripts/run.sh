#!/bin/sh
go build .
sudo ./brewTV -start &
tail -f ./brewTV.log
