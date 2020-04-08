# Terraform provider for CloudKarafka

Setup your CloudKarafka cluster with Terraform

## Install

```sh
git clone https://github.com/cloudkarafka/terraform-provider.git
cd terraform-provider
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.

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



