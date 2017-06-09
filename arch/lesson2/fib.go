func fib(n int) int {
	for i := 0; i < n; i++ {
		x, y = y, x+y
	}
	return x
}
