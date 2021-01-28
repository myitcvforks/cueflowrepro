package test

import "encoding/json"

A: {
	input: """
		{
			"hello": "world"
		}
		"""
	output: string

	unmarshalled: json.Unmarshal(output)
}

B: {
	input:  A.unmarshalled.hello
	output: string
}
