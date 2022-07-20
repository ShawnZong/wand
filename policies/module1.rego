package main

import future.keywords

default allow = false

allow {
	is_admin
}

optional[{key: msg}] {
	is_admin

	key := "$..ContainerDefinitions..Image"
	msg := "optional msg 1"
}

optional[{key: msg}] {
	is_admin

	key := "$..ContainerDefinitions..Image"
	msg := "optional msg 3"
}
optional[{key: msg}] {
	is_admin

	key := "$..Volumes"
	msg := "optional msg 2"
}

prohibited[{key: msg}] {
	is_admin

	key := "$..SourcePath"
	msg := "prohibited msg 1"
}

prohibited[{key: msg}] {
	is_admin

	key := "$..Volumes"
	msg := "prohibited msg 2"
}

mandatory[{key:msg,"templateRef":templateRef}]{
	is_admin

	key:="$..ContainerDefinitions"
	msg:="mandatory msg 1"
	templateRef:="../../templates/ref.yaml"
}

mandatory[{key:msg,"templateRef":templateRef}]{
	is_admin

	key:="$..Volumes"
	msg:="mandatory msg 2"
	templateRef:="../../templates/ref2.yaml"
}

is_admin = true
