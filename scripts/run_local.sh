#!/bin/sh
go build .
./brewTV -start &
tail -f ./brewTV.log
