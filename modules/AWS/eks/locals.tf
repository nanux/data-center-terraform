
locals {
  # When the framework get uninstalled, there is no appNode attribute for node_group and so next install will break
  # because eks is still is not created but terraform tries to evaluate the cluster_asg_name
  cluster_asg_name = try(module.eks.node_groups.appNodes.resources[0].autoscaling_groups[0].name, null)

  autoscaler_service_account_namespace = "kube-system"
  autoscaler_service_account_name      = "cluster-autoscaler-aws-cluster-autoscaler-chart"

  autoscaler_version = "9.16.0"
}