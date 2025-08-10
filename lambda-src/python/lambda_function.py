#!/usr/bin/env python3
from cffi import FFI


def lambda_handler(_event, _context):
    ffi = FFI()
    ffi.cdef("int abs(int);")
    clib = ffi.dlopen(None)
    return clib.abs(-123)


if __name__ == "__main__":
    lambda_handler(None, None)
