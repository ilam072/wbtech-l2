#!/bin/bash

#first test
sort -k 2 files/K.txt > results/K_sort.txt
go run ../cmd/main.go -k 2 files/K.txt > results/K_app.txt

DIFF=$(diff results/K_sort.txt results/K_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -k PASSED"
else
    echo "TEST -k FAIL"
fi

#second test
sort -n files/N.txt > results/N_sort.txt
go run ../cmd/main.go -n files/N.txt > results/N_app.txt

DIFF=$(diff results/N_sort.txt results/N_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -n PASSED"
else
    echo "TEST -n FAIL"
fi

#third test
sort -r files/K.txt > results/R_sort.txt
go run ../cmd/main.go -r files/K.txt > results/R_app.txt

DIFF=$(diff results/R_sort.txt results/R_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -r PASSED"
else
    echo "TEST -r FAIL"
fi

#fourth test
sort -u files/U.txt > results/U_sort.txt
go run ../cmd/main.go -u files/U.txt > results/U_app.txt

DIFF=$(diff results/U_sort.txt results/U_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -u PASSED"
else
    echo "TEST -u FAIL"
fi


#fifth test
sort -M files/M.txt > results/M_sort.txt
go run ../cmd/main.go -M files/M.txt > results/M_app.txt

DIFF=$(diff results/M_sort.txt results/M_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -M PASSED"
else
    echo "TEST -M FAIL"
fi

#sixth test
sort -b files/B.txt > results/B_sort.txt
go run ../cmd/main.go -b files/B.txt > results/B_app.txt

DIFF=$(diff results/B_sort.txt results/B_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -b PASSED"
else
    echo "TEST -b FAIL"
fi

#seventh test
sort -c files/C.txt > results/C_sort.txt 2>&1
go run ../cmd/main.go -c files/C.txt > results/C_app.txt

DIFF=$(diff results/C_sort.txt results/C_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -c PASSED"
else
    echo "TEST -c FAIL"
fi

# eighth test
sort -h files/H.txt > results/H_sort.txt
go run ../cmd/main.go -H files/H.txt > results/H_app.txt

DIFF=$(diff results/H_sort.txt results/H_app.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -H PASSED"
else
    echo "TEST -H FAIL"
fi