#!/usr/bin/env bats

setup() {
  REPO_ROOT="$(cd "$(dirname "$BATS_TEST_FILENAME")"/.. && pwd)"
  ORDNA="$REPO_ROOT/src/ordna"
}

@test "prints help and exits 0" {
  run bash "$ORDNA" --help
  [ "$status" -eq 0 ]
  [[ "$output" == *"Usage: ordna"* ]]
}

@test "errors when missing source/dest" {
  run bash "$ORDNA" -m
  [ "$status" -ne 0 ]
  [[ "$output" == *"must specify at least one source"* ]]
}

