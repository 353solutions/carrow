# distutils: language=c++
# cython: language_level=3

from pyarrow.lib cimport *
from cython.operator cimport dereference as deref
cimport c_bindings

def try_build():
    table = c_bindings.CreateTable()
    l = callArrow(table)
    return l

cdef callArrow(void* table):
    cdef const shared_ptr[CTable]* table_pointer = <shared_ptr[CTable]*>table
    if table_pointer == NULL:
        raise TypeError("pointer is not a table")
    cdef shared_ptr[CTable] t = deref(table_pointer)
    print("bindings table: rows:",t.get().num_rows(),"columns:",t.get().num_columns())
    return pyarrow_wrap_table(t)