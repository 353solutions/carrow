from bindings import try_build

table = try_build()
print(table)
c1 = table[0]
c2 = table[1]
print(c1.data)
print(c2.data)
