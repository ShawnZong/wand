package main

import future.keywords

default allow = false

allow {
	is_admin
}

optional[{key: msg}] {
	is_admin

	key := "$..ContainerDefinitions..[?(@.Image=='busybox')].Name"
	msg := "example appended comment"
}

optional[{key: msg}] {
	is_admin

	key := "example_key2"
	msg := "example message2"
}

is_admin {
	# "admin" in input.subject.groups
	true
}
