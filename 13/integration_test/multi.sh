#!/bin/bash

mkdir -p results

compare_files() {
    local file1="$1"
    local file2="$2"
    local test_name="$3"

    # Приведение окончаний строк к LF и удаление пустых строк
    local tmp1=$(mktemp)
    local tmp2=$(mktemp)
    sed 's/\r$//' "$file1" | sed '/^$/d' > "$tmp1"
    sed 's/\r$//' "$file2" | sed '/^$/d' > "$tmp2"

    if diff -u "$tmp1" "$tmp2" >/dev/null; then
        echo "TEST $test_name PASSED"
    else
        echo "TEST $test_name FAIL"
    fi

    rm -f "$tmp1" "$tmp2"
}

# first test
cut -f 3 -d ':' testfile.txt  > results/result_cut_multi1.txt
go run ../cmd/main.go -f 3 -d ':' testfile.txt > results/result_app_multi1.txt
compare_files results/result_cut_multi1.txt results/result_app_multi1.txt "-f -d"

# second test
cut -f 1 -s testfile.txt > results/result_cut_multi2.txt
go run ../cmd/main.go -f 1 -s  testfile.txt > results/result_app_multi2.txt
compare_files results/result_cut_multi2.txt results/result_app_multi2.txt "-f -s"

# third test
cut -f 3 -d ':' -s testfile.txt  > results/result_cut_multi3.txt
go run ../cmd/main.go -f 3 -d ':' -s testfile.txt > results/result_app_multi3.txt
compare_files results/result_cut_multi3.txt results/result_app_multi3.txt "-f -d -s"
