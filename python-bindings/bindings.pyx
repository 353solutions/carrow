# distutils: language=c++

from pyarrow.lib cimport *
from carrow.lib import *


def build():
    # Just an example function accessing both the pyarrow Cython API
    # and the Arrow C++ API
    #cdef shared_ptr[CArray] arr = pyarrow_unwrap_array(obj)
    #if arr.get() == NULL:
    #    raise TypeError("not an array")
    #return arr.get().length()
    carrow.Build()

