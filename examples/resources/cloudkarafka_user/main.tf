resource "cloudkarafka_user" "myuser" {
  instance_id = cloudkarafka_instance.cluster.id
  name        = "user1"
  type        = "sasl"
}
