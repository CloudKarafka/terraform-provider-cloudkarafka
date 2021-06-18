# Terraform provider for CloudKarafka

Setup your CloudKarafka cluster with Terraform

## Install

```sh
git clone https://github.com/e-conomic/terraform-provider-cloudkarafka
cd terraform-provider-cloudkarafka
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.


Generally provider plugins are expected to be fetched from registry.terraform.io, We can create custom provider plugins and inform terraform to load them from local filesystem rather than from registry. More about this at https://www.terraform.io/docs/cli/config/config-file.html#provider-installation

In brief, the custom provider plugins are expected to be in following path/formats on various OS.
```
Windows:
%APPDATA%/terraform.d/plugins and %APPDATA%/HashiCorp/Terraform/plugins

Mac OS X: 
$HOME/.terraform.d/plugins/, ~/Library/Application Support/io.terraform/plugins, and /Library/Application Support/io.terraform/plugins

Linux and other Unix-like systems:
$HOME/.terraform.d/plugins/
```

Under these directrories the further provider path is expected as follows along with strict binary naming convention:
```
HOSTNAME/NAMESPACE/TYPE/VERSION/TARGET/terraform-provider-TYPE_VERSION
```

So for an example system of linux with amd64 architecture and cloud karafka version 0.9.0, the expected path is :
```
~/.terraform.d/plugins/github.com/e-conomic/cloudkarafka/0.9.0/linux_amd64/terraform-provider-cloudkarafka_v0.9.0
```

## Cloud Karafka API

There are two api's for cloud karafka:
Customer API - https://docs.cloudkarafka.com/
Console API - https://docs.cloudkarafka.com/cloudkarafka_api.html


The `cloudkarafka` provider is abstraction over Customer API and `cloudkarafka_instance` resource is abstraction over https://docs.cloudkarafka.com/#create-an-instance api. The params for the resource are similar to the create api.

## Requirements 
To access Cloud karafka api, you need api token which can be generated at https://customer.cloudkarafka.com/apikeys


## Example

```hcl
provider "cloudkarafka" {}

resource "cloudkarafka_instance" "kafka_bat" {
  name   = "terraform-provider-test"
  plan   = "bat-3"
  region = "amazon-web-services::us-east-1"
  vpc_subnet = "10.201.0.0/24"
}

output "kafka_brokers" {
  value = "${cloudkarafka_instance.kafka_bat.brokers}"
}
```



