# Wand
Wand is a command line tool to implement the Policy Champion Model (PCM). It helps to refine configuration files based on Rego policies. PCM provides a mean for a senior developer to share their best practices centrally and then junior developers can apply those best practices on their own. It's a process of knowledge sharing.

Wand reads a configuration file and then executes policies, based on the evaluation results of policies, it performs actions to modify the configuration and generate an updated file. The policies are written in a policy language called [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/), and they are evaluated by [Open Policy Agent](https://www.openpolicyagent.org).

Wand supports 3 types of policies:
* optional: It appends comments to specific configurations. Developers can use this policy to add reference materials or communicate with each other.
* prohibited: It removes specific configuration fields and adds comments as an explanation. Developers can use this policy to remove some configurations which violate security practices or company policies.
* mandatory: It appends template configuration into an existing configuration and adds comments as an explanation. Developers can use this policy to add some missing or recommended configurations.

**Currently, Wand only supports configuration file in YAML.**

# Options
## -i
This option sets the path of the input YAML file.

**Example**:
``` bash
wand -i /path/to/my/folder/configuration.yaml
```
## -o
This option sets the path for output updated YAML file, default is ./updated_`<`input file name`>`.

**Example**:
``` bash
wand -i /path/to/my/folder/configuration.yaml -o /path/to/output/folder/output.yaml
```
## -p
This option tells Wand where to locate policy files. By default, it looks for folder `policies` in the current path.

**Example**:

Assumed that policies are stored in folder `/path/to/policy/folder`.

``` bash
wand -i /path/to/my/folder/configuration.yaml -p /path/to/policy/folder
```

## -namespace
Each policy has a namespace, by default, Wand looks for policies specified in the `main` namespace, but one can change to look for another namespace too.

**Example**:
``` bash
wand -i /path/to/my/folder/configuration.yaml -namespace master
```

# Usage
To learn how to write policies in Rego, please check the official [documentation](https://www.openpolicyagent.org/docs/latest/policy-language/) from OPA.
Wand supports 3 types of policy rules: `optional`, `prohibited`, and `mandatory`. They look for YAML nodes by [JSONPath syntax](https://github.com/vmware-labs/yaml-jsonpath), and edit found Nodes accordingly.

* `optional`:
variable `key` specifies which key in YAML you want to modify
variable `msg` is a comment that will be appended to the key
``` css
optional[{key: msg}] {
 <rego rules>

 key := <yaml key>
 msg := <comment to be added>
}
```
* `prohibited`:
a prohibited rule will remove YAML nodes
variable `key` specifies which key in YAML you want to modify
variable `msg` is a comment that will be appended to the key
``` css
prohbited[{key: msg}] {
 <rego rules>

 key := <yaml key>
 msg := <comment to be added>
}
```
* `mandatory`:
a mandatory rule will concatenate template YAML after the chosen key.
variable `key` specifies after which key in YAML you want to concatenate
variable `msg` is a comment that will be appended to the key
variable `templateRef` specifies the path of template yaml file
``` css
mandatory[{key: msg, "templateRef": templateRef }] {
 <rego rules>

 key := <yaml key>
 msg := <comment to be added>
 templateRef := <path of yaml template to be appended after key>
}
```


For example, given input configuration YAML `example/input.yaml`, a policy `policies/example.rego`, and template YAML `templates/template.yaml`:
`example/input.yaml`:
```yaml
Configuration:
  Type: "web"
  Protocol: "foo"
  Port: "80"
```

`policies/example.rego`:
``` css
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
```

`template/template.yaml`:
```yaml
Name: website-service
Essential: true
Image: foo
Memory: 128
Environment:
 - Name: URL
 Value: bar
PortMappings:
 - ContainerPort: 8000
LogConfiguration:
 LogDriver: awslogs
 Options:
 awslogs-group: !Ref CloudWatchLogsGroup
 awslogs-region: !Ref AWS::Region
```

Run command:
```bash
wand -i example/input.yaml -o output.yaml -p policies
```
Wand will generate an output yaml file `output.yaml`, which will append comments and yaml Nodes according to Rego rules:
`output.yaml`:
```yaml
Configuration:
  Type: "web"
  Protocol: "foo"
  # secure connection not available, please establish ssl
  Port: null
  # HTTP is prohibited, please use HTTPS instead

  # A web configuration needs to configure server
  Name: website-service
  Essential: true
  Image: foo
  Memory: 128
  Environment:
    - Name: URL
  Value: bar
  PortMappings:
    - ContainerPort: 8000
  LogConfiguration:
  LogDriver: awslogs
  Options:
  awslogs-group: !Ref CloudWatchLogsGroup
  awslogs-region: !Ref AWS::Region
```