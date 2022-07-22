package main

is_encrypted {
	input.Configuration.Protocol == "ssl"
}

optional[{key: msg}] {
	not is_encrypted

	key := "$.Configuration.Protocol"
	msg := "secure connection not available, please establish ssl"
}

prohibited[{key: msg}] {
	input.Configuration.Port == "80"

	key := "$.Configuration.Port"
	msg := "HTTP is prohibited, please use HTTPS instead"
}

mandatory[{key: msg, "templateRef": templateRef}] {
	input.Configuration.Type == "web"

	key := "$.Configuration"
	msg := "A web configuration needs to configure server"
	templateRef := "templates/template.yaml"
}
