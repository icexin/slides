package main

func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func main() {
	// 直接指针
	p := &Point{1, 2}
	p.ScaleBy(2)

	// 声明结构体后再用指针指向
	p1 := Point{1, 2}
	p2 := &p1
	p2.ScaleBy(2)

	// 使用结构体调用，隐式取地址
	p3 := Point{1, 2}
	p3.ScaleBy(2) // 等价于 (&p3).ScaleBy(2)
}
