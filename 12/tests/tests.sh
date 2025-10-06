#!/bin/bash

#first test
grep -A 2 grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -A 2 grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -A PASSED"
else
    echo "TEST -A FAIL"
fi

#second test
grep -B 2 grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -B 2 grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -B PASSED"
else
    echo "TEST -B FAIL"
fi

#third test
grep -C 2 grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -C 2 grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -C PASSED"
else
    echo "TEST -C FAIL"
fi

#fourth test
grep -c grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -c  grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -c PASSED"
else
    echo "TEST -c FAIL"
fi


#fifth test
grep -i grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -i grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -i PASSED"
else
    echo "TEST -i FAIL"
fi

#sixth test
grep -v grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -v grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -v PASSED"
else
    echo "TEST -v FAIL"
fi

#seventh test
grep -F grep testfile.txt | grep -v '^--$' > results/result_grepF.txt
go run ../cmd/main.go -F grep testfile.txt > results/result_appF.txt

DIFF=$(diff results/result_appF.txt results/result_grepF.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -F PASSED"
else
    echo "TEST -F FAIL"
fi

#eigths test
grep -n grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -n grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -n PASSED"
else
    echo "TEST -n FAIL"
fi