resource "cloudkarafka_topic" "my_topic1" {
  instance_id        = cloudkarafka_instance.cluster.id
  name               = "mytopic1"
  replication_factor = 1
  partitions         = 12
  config             = {
	  segment_bytes   = 100000,
	  retention_bytes = 104857600,
	  cleanup_policy  = "delete"
  }
}
