resource "cloudkarafka_kafkaconfig" "kafkaconfig" {
  instance_id = cloudkarafka_instance.cluster.id
  auto_create_topics_enable = true
  num_io_threads = 10
}

