  #!/bin/bash

#first test
grep -A 2 -B 2 grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -A 2 -B 2 grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -A -B PASSED"
else
    echo "TEST -A -B FAIL"
fi

#second test
grep -A 2 -B 1 -C 1 grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -A 2 -B 1 -C 1 grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -A -B -C PASSED"
else
    echo "TEST -A -B -C FAIL"
fi

#third test
grep -C 2 -n -i grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -C 2 -n -i grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -C 2 -n -i PASSED"
else
    echo "TEST -C 2 -n -i FAIL"
fi

#fourth test
grep -n -F grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -n -F  grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -n -F PASSED"
else
    echo "TEST -n -F FAIL"
fi


#fifth test
grep -i -c -v grep testfile.txt | grep -v '^--$' > results/result_grep.txt
go run ../cmd/main.go -i -c -v grep testfile.txt > results/result_app.txt

DIFF=$(diff results/result_app.txt results/result_grep.txt)
if [ "$DIFF" = "" ]
then
    echo "TEST -i -c -v PASSED"
else
    echo "TEST -i -c -v FAIL"
fi