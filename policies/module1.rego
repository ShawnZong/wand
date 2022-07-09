package example.authz

import future.keywords

default allow = false

allow {
	is_admin
}

optional[{"key": key, "msg": msg}] {
	is_admin

	key := "example_key"
	msg := "example message"
}

optional[{"key": key, "msg": msg}] {
	is_admin

	key := "example_key2"
	msg := "example message2"
}

is_admin {
	"admin" in input.subject.groups
}
