set tabstop=2
set shiftwidth=2
set softtabstop=2

" C++
func! CLangFormat()
    silent ! clang-format -i %
    e
endfunc
comm! CLangFormat call CLangFormat()

au BufWritePost *.cc call CLangFormat()
au BufWritePost *.h call CLangFormat()
