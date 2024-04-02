#!/bin/sh
go build .
go run ./brewTV -start &
tail -f ./brewTV.log
