#! /bin/bash

# example: ./scover.sh ./libs/rdb/...

t="/tmp/go-cover.$$.tmp"
go test -coverprofile=$t $@ && go tool cover -html=$t && unlink $t