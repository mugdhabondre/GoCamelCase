#!/usr/bin/env bash
go install -x oxforddict/oxforddict.go
go install -x gocamelcaseimpl/gocamelcaseimpl.go
go build -o build/camelcaseapp -x httphandler.go



