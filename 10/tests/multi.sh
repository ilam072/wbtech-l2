#!/bin/bash

# first test
sort -nk 2 files/K.txt > results/KN_sort.txt
go run ../cmd/main.go -nk 2 files/K.txt > results/KN_app.txt

DIFF=$(diff results/KN_sort.txt results/KN_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -nk PASSED"
else
    echo "TEST -nk FAIL"
fi

# second test
sort -uM files/M.txt > results/uM_sort.txt
go run ../cmd/main.go -uM files/M.txt > results/uM_app.txt

DIFF=$(diff results/uM_sort.txt results/uM_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -uM PASSED"
else
    echo "TEST -uM FAIL"
fi

# third test
sort -nr files/N.txt > results/NR_sort.txt
go run ../cmd/main.go -nr files/N.txt > results/NR_app.txt

DIFF=$(diff results/NR_sort.txt results/NR_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -nr PASSED"
else
    echo "TEST -nr FAIL"
fi

# fourth test
sort -Mr files/M.txt > results/MR_sort.txt
go run ../cmd/main.go -Mr files/M.txt > results/MR_app.txt

DIFF=$(diff results/NR_sort.txt results/NR_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -nr PASSED"
else
    echo "TEST -nr FAIL"
fi
