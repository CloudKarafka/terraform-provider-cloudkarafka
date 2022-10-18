terraform {
  required_providers {
    cloudkarafka = {
      source = "cloudkarafka/cloudkarafka"
      version = "~>1.0"
    }
  }
}

// Configuration of the CloudKarafka provider
provider "cloudkarafka" {
  apikey  = "<api-key>"
}

// CloudKarafka Cluster
resource "cloudkarafka_instance" "instance" {
  name = "Test cluster"
  plan = "mouse-1"
  region = "amazon-web-services::eu-north-1"
  tags = ["test"]
  kafka_version = "2.8.0"
  vpc_subnet = "10.90.92.0/24"
}
