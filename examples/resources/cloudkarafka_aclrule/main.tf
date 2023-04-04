resource "cloudkarafka_aclrule" "myuser_rule1" {
  instance_id = cloudkarafka_instance.cluster.id
  username     = cloudkarafka_user.myuser.name
  operation        = "read"
  resource         = "topic"
  resource_pattern = "sample-"
  resource_pattern_type = "prefixed"
}
