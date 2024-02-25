package main

//helper

func delay(d, n int) []int {
	var ds = make([]int, n)
	for i := 0; i < n; i++ {
		ds[i] = d
	}
	return ds
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func gcd(a, b int) int {
	if a == 0 {
		return b
	}
	for b != 0 {
		h := a % b
		a = b
		b = h
	}
	return a
}

func last[T any](arr []T) T {
	return arr[len(arr)-1]
}

func round(f float64) int {
	if f > 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}
