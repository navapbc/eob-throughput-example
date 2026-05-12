#!/bin/bash

rm -f throughput 

go build

# Be noisy on the first run only to get the headers
noisy=1

for write in 0 1; do 
    for conc in 1 2 4 8; do
        for bundle in 10 20 50 100 200 500; do
            rm -f tput*.sqlite ; ./throughput filter -d ./bundles/ -r 10000 -c $conc -w $write -b $bundle -n $noisy
            if [ "$noisy" -eq 1 ]; then
                noisy=0;
            fi
        done
    done
done
