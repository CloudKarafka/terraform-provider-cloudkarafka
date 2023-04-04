# Manage example cluster.
resource "cloudkarafka_instance" "cluster" {
  name   = "test"
  plan   = "dedicated_2-1"
  region = "amazon-web-services::us-east-1"
  kafka_version = "3.3.1"
  disk_size = 128
  tags = ["terraform", "testing"]
}
