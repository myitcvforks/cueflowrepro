package test

A: {
	input:  "foobar"
	output: string
}

B: {
	input:  "interpolated \(A.output)"
	output: string
}
