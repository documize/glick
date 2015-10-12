# Copyright 2013-2014 Documize (http://www.documize.com)
# script to run "go vet", "golint" and "errcheck" to check correctness, style and error handling
echo "go vet:"
go vet github.com/documize/glick github.com/documize/glick/glpie
echo "go tool vet:"
go tool vet -shadowstrict .
echo "golint:"
golint *.go */*.go
echo "errcheck:"
errcheck -blank=true ./...
