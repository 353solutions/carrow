from bindings import try_build
from ctypes import cdll
from os.path import dirname, abspath

here = dirname(abspath(__file__))
lib = cdll.LoadLibrary(f'{here}/carrow_bindings.so')

try_build()