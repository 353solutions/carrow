# distutils: language=c++
# cython: language_level=3

from pyarrow.lib cimport *
cimport c_bindings

def try_build():
    table = c_bindings.CreateTable()
    l = callArrow(table)
    print(l)

cdef callArrow(void* table):
    cdef CTable* table_pointer = <CTable*>table
    if table_pointer == NULL:
        raise TypeError("not an array")
    return (table_pointer.num_columns(),table_pointer.num_rows())