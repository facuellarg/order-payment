locals {
  function_name = "make-order"
  function_name_complete_order = "complete-order"
  src_path_make_order     = "${path.module}/cmd/make-order/"
  src_path_complete_order     = "${path.module}/cmd/complete-order/"
  binary_name  = local.function_name
  binary_name_complete_order  = local.function_name_complete_order
  binary_path  = "${path.module}/tf_generated/${local.binary_name}"
  binary_path_complete_order = "${path.module}/tf_generated/${local.binary_name_complete_order}"
  archive_path = "${path.module}/tf_generated/${local.function_name}.zip"
  archive_path_complete_order = "${path.module}/tf_generated/${local.function_name_complete_order}.zip"
}

output "binary_path" {
  value = local.binary_path
}