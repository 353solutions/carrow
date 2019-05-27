from distutils.core import setup
from Cython.Build import cythonize

from os.path import dirname, abspath
import numpy as np
import pyarrow as pa

here = dirname(abspath(__file__))


ext_modules = cythonize("bindings.pyx")

inc_dirs = [np.get_include(), pa.get_include(), here]
lib_dirs = [here] + pa.get_library_dirs()
libs = pa.get_libraries() + ["carrow_bindings"]

for ext in ext_modules:
    ext.include_dirs.extend(inc_dirs)
    ext.library_dirs.extend(lib_dirs)
    ext.libraries.extend(libs)

#     if os.name == 'posix':
#         ext.extra_compile_args.append('-std=c++11')

#     # Try uncommenting the following line on Linux
#     # if you get weird linker errors or runtime crashes
# ext.define_macros.append(("_GLIBCXX_USE_CXX11_ABI", "0"))


setup(ext_modules=ext_modules)
