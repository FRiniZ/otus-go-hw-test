#!/usr/bin/env bash

result=$(echo $FOO)
expected="foo"

[ "${result}" = "${expected}" ] || (echo -e "invalid output: ${result}" && exit 1)

echo "PASS"


