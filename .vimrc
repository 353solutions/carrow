au FileType cpp setl tw=78
au FileType cpp setl tabstop=2
au FileType cpp setl shiftwidth=2
au FileType cpp setl softtabstop=2

" C++
func! CLangFormat()
    silent ! clang-format -i %
    e
endfunc
comm! CLangFormat call CLangFormat()

au BufWritePost *.cc call CLangFormat()
au BufWritePost *.h call CLangFormat()
