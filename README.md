# PolicyChampion
PolicyChampion is a command line tool to implement policy champion model. It helps to refine configuration files based on Rego policies. The policy champion model provides a mean for senior developer to share their best practices centrally and then junior developers can apply those best practices on their own. It's a process of knowledge share.

PolicyChampion reads a configuration file and then execute policies, based on evaluation results of policies, it performs actions to modify the configuration and generate an updated file. The policies are written in a policy language called [Rego]{https://www.openpolicyagent.org/docs/latest/policy-language/}, they evaluated by [Open Policy Agent]{https://www.openpolicyagent.org}.

PolicyChampion supports 3 types of policies:
* optional: It appends comment to specific configuration. Developers can use this policy to add referece materials or communicate with each other.
* prohibited: It removes specific configuration fields and add comment as explanation. Developers can use this policy to remove some configurations which violate security practices or company policies.
* mandatory: It appends template configuration into an existing configuration and add comment as explanation. Developers can use this policy to add some missing or recommended configurations.

**Currently, PolicyChampion only supports configuration file in YAML.**

# Options
## -i
This option sets path of input YAML file.

**Example**:
``` bash
policychampion -i /path/to/my/folder/configuration.yaml
```
## -o
This option sets path for output updated YAML file, default is ./updated_`<`input file name`>`.

**Example**:
``` bash
policychampion -i /path/to/my/folder/configuration.yaml -o /path/to/output/folder/output.yaml
```
## -p
This option tells policychampion where to locate policy files. By default, it looks for folder `policies` in current path.

**Example**:

Assumed that policies are stored in folder `/path/to/policy/folder`.

``` bash
policychampion -i /path/to/my/folder/configuration.yaml -p /path/to/policy/folder
```

## -namespace
Each policy has a namespace, by default, PolicyChampion looks for policies specified in the `main` namespace, but one can change to look for another namespace too.

**Example**:
``` bash
policychampion -i /path/to/my/folder/configuration.yaml -namespace master
```

# Usage
To learn how to write polices in Rego, please check the official [documentation]{https://www.openpolicyagent.org/docs/latest/policy-language/} from OPA.
PolicyChampion supports 3 types of policy rules: `optional`, `prohibited`, and `mandatory`:
* `optional`:
``` css
optional[{key: msg}] {
	<rego rules>

	key := <yaml key>
	msg := <comment to be added>
}
```
* `prohibited`:
``` css
prohbited[{key: msg}] {
	<rego rules>

	key := <yaml key>
	msg := <comment to be added>
}
```
* `mandatory`:
``` css
mandatory[{key: msg, "templateRef": templateRef }] {
	<rego rules>

    key := <yaml key>
	msg := <comment to be added>
	templateRef := <path of yaml template to be appended after key>
}
```


For example, given a policy `policies/example.rego` and input configuration YAML `example.yaml`:
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

prohbited[{key: msg}] {
	input.Configuration.Port == "80"

	key := "$.Configuration.Port"
	msg := "HTTP is prohibited, please use HTTPS instead"
}

mandatory[{key: msg, "templateRef": templateRef}] {
	input.Configuration.Type == "web"

	key := "$.Configuration"
	msg := "A web configuration needs to configure server"
	templateRef := "path/to/template/folder/serverConfiguration.yaml"
}
```