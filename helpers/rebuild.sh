#!/bin/bash

rm -f bin/iterator
go build -o bin/iterator .
chmod +x bin/iterator
./bin/iterator
