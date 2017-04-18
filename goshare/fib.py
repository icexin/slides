def fib(n):
    a, b = 0, 1
    for i in range(n):
        a, b = b, a + b
        yield a

for x in fib(10):
    print x
