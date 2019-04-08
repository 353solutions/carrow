# distutils: language=c++
# cython: language_level=3

from pyarrow.lib cimport *
cimport c_bindings

def try_build():
    # Just an example function accessing both the pyarrow Cython API
    # and the Arrow C++ API
    #cdef shared_ptr[CArray] arr = pyarrow_unwrap_array(obj)
    #if arr.get() == NULL:
    #    raise TypeError("not an array")
    #return arr.get().length()
    print("hello")
    c_bindings.Build()

