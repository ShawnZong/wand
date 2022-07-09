package example.authz

default hello = false

optional[{"key": key, "msg": msg}] {
	is_admin

	key := "example_key3"
	msg := "example message3"
}

hello {
	m := input.message
	m == "world"
}
