#!/bin/python3
import sys
from typing import Generator


def read_input() -> Generator[str, str, None]:
    for line in sys.stdin:
        yield line.strip("\n")


def print_as_it_comes(input: Generator[str, str, None]):
    for s in input:
        if s == "close":
            print("exiting")
            exit(0)
        print(s)


if __name__ == "__main__":
    print_as_it_comes(read_input())
