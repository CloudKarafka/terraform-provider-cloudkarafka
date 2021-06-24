terraform {
  required_providers {
    cloudkarafka = {
      source = "github.com/e-conomic/cloudkarafka"
    }
  }
}

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
