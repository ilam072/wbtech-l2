#!/bin/bash

# Папка для результатов
mkdir -p results

# Функция для сравнения двух файлов
compare_files() {
    local file1="$1"
    local file2="$2"
    local test_name="$3"

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

# TEST 1: простое поле
cut -f 2 testfile.txt > results/result_cut.txt
go run ../cmd/main.go -f 2 testfile.txt > results/result_app.txt
compare_files results/result_cut.txt results/result_app.txt "-f"

# TEST 2: диапазон полей
cut -f 2-3 testfile.txt > results/result_cut.txt
go run ../cmd/main.go -f 2-3 testfile.txt > results/result_app.txt
compare_files results/result_cut.txt results/result_app.txt "-f with interval"

# TEST 3: другой разделитель
cut -f 2 -d ':' testfile.txt > results/result_cut.txt
go run ../cmd/main.go -f 2 -d ':' testfile.txt > results/result_app.txt
compare_files results/result_cut.txt results/result_app.txt "-d"

# TEST 4: флаг -s
cut -f 1 -s testfile.txt > results/result_cut.txt
go run ../cmd/main.go -f 1 -s testfile.txt > results/result_app.txt
compare_files results/result_cut.txt results/result_app.txt "-s"
