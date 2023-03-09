locals {
  function_name = "make-order"
  src_path_make_order     = "${path.module}/../cmd/make-order/"
  binary_name  = local.function_name
  binary_path  = "${path.module}/tf_generated/${local.binary_name}"
  archive_path = "${path.module}/tf_generated/${local.function_name}.zip"
}

locals {

  function_name_complete_order = "complete-order"
  src_path_complete_order     = "${path.module}/../cmd/complete-order/"
  binary_path_complete_order = "${path.module}/tf_generated/${local.binary_name_complete_order}"
  binary_name_complete_order  = local.function_name_complete_order
  archive_path_complete_order = "${path.module}/tf_generated/${local.function_name_complete_order}.zip"
}