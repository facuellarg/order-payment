locals {
  function_name = "process-payment"
  src_path      = "${path.module}/cmd/process-payment/"
  binary_name  = local.function_name
  binary_path  = "${path.module}/tf_generated/${local.binary_name}"
  archive_path = "${path.module}/tf_generated/${local.function_name}.zip"
}

output "binary_path" {
  value = local.binary_path
}