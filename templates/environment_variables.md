---
page_title: "Environment Variables"
title: "Environment Variables"
linkTitle: "Environment Variables"
---

<!--
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Environment Variables Overview

*This Guide will describe two common ways on how to use environment variables with terraform. For more information you may check out the [terraform docs](https://developer.hashicorp.com/terraform/cli/config/environment-variables) directly.*

<br>

## Export or Directly Set During Terraform Command

*This section will desribe how to pass in environment variables directly during the terraform command.*

main.tf
```terraform
// Create a variable which will be assigned during the command
variable "foo" {
  type = string
}

data "example_datasource" "all" {
  filter {
    // Reference that variable how you would normally in terraform 
    name = [var.foo]
  }
}
```

Using the above example, to set the environment variable for `foo` you must use the prefix of `TF_VAR_` + the name of the variable i.e: `TF_VAR_foo`. This will translate to setting the variable during the command.

> Example Command: `TF_VAR_foo="example-name" terraform plan`

or

> Export and then run the terraform command
```
export TF_VAR_foo="example-name"
```


This will execute terraform plan with the value of `example-name` for the variable `foo`.

<br>

## Using the External Terraform Provider

*This section will describe how to use environment variables using the [external](https://registry.terraform.io/providers/hashicorp/external/latest/docs/data-sources/external) terraform provider*

env.sh
```
cat << EOF
{
 "foo": "test"
}
EOF
```

main.tf
```terraform
terraform {
  required_providers {
    external = {
      source = "hashicorp/external"
      version = "2.3.4"
    }
  }
}

# Run the script to get the environment variables of interest.
# This is a datasource, so it will run at plan time.
data "external" "env" {
  program = ["bash", "env.sh"]
}

# Show the results of running the datasource. This is a map of environment
# variable names to their values.
output "env" {
  value = data.external.env.result
}

data "example_datasource" "all" {
  filter {
    name = [data.external.env.result["foo"]]
  }
}
```

In this example there are 2 files. 

1. `env.sh` This creates the input for the environment variables you would like to pass into and use in the terraform configuration `main.tf`
2. `main.tf` this defines the `external` datasource which takes in the output from the `external` datasource and then it can be used in the `example_datasource` filter i.e `data.external.env.result["foo"]`
 