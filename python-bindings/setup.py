from distutils.core import setup
from Cython.Build import cythonize

import os
from os.path import dirname, abspath
import numpy as np
import pyarrow as pa

here = dirname(abspath(__file__))


ext_modules = cythonize("bindings.pyx")

for ext in ext_modules:

    ext.include_dirs.append('./')
    ext.libraries.extend(["carrow"])

    # The Numpy C headers are currently required
    ext.include_dirs.append(np.get_include())
    ext.include_dirs.append(pa.get_include())
    ext.libraries.extend(pa.get_libraries())
    ext.library_dirs.extend(pa.get_library_dirs())

    if os.name == 'posix':
        ext.extra_compile_args.append('-std=c++11')

    # Try uncommenting the following line on Linux
    # if you get weird linker errors or runtime crashes
    # ext.define_macros.append(("_GLIBCXX_USE_CXX11_ABI", "0"))


setup(ext_modules=ext_modules)
