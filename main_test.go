package main

import (
  "testing"
)

func sum(a, b int) int {
  return a+b
}

func TestSum(t *testing.T) {
  c := sum(5, 5)
  if c != 10 {
    t.Fail()
  }

}
