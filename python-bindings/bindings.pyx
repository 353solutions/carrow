# distutils: language=c++

from pyarrow.lib cimport *

# failed to load like this
# cdef extern from 'libcarrow.h':
#     void Build()


def try_build():
    # Just an example function accessing both the pyarrow Cython API
    # and the Arrow C++ API
    #cdef shared_ptr[CArray] arr = pyarrow_unwrap_array(obj)
    #if arr.get() == NULL:
    #    raise TypeError("not an array")
    #return arr.get().length()
    print("hello")

