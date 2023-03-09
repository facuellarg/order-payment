locals {
  function_name_order_created = "order_created"
  src_path_order_created     = "${path.module}/cmd/order-created/"
  binary_name_order_created  = local.function_name_order_created
  binary_path_order_created = "${path.module}/tf_generated/${local.binary_name_order_created}"
  archive_path_order_created = "${path.module}/tf_generated/${local.function_name_order_created}.zip"
}
