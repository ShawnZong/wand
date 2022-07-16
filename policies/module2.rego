package main

default hello = false

optional[{key: msg}] {
	is_admin

	key := "example_key3"
	msg := "example message3"
}

optional[{key: msg}] {
	is_admin

	key := "example_key3"
	msg := "example message4"
}

hello {
	m := input.message
	m == "world"
}
