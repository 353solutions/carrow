# distutils: language=c++
# cython: language_level=3

from pyarrow.lib cimport *
cimport c_bindings

def try_build():
    table = c_bindings.CreateTable()
    l = callArrow(table)
    print(l)

cdef callArrow(void* table):
    cdef shared_ptr[CTable]* table_pointer = <shared_ptr[CTable]*>table
    if table_pointer == NULL:
        raise TypeError("not a table")
    return (table_pointer.get().num_columns(),table_pointer.get().num_rows())