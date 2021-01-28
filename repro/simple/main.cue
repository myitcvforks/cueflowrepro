package test

A: {
	input:  "foobar"
	output: string
}

B: {
	input:  A.output
	output: string
}
