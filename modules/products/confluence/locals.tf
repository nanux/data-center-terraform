locals {
  product_name = "confluence"

  helm_chart_repository         = "https://atlassian.github.io/data-center-helm-charts"
  confluence_helm_chart_version = var.confluence_configuration["helm_version"]

  confluence_software_resources = {
    "minHeap" : var.confluence_configuration["min_heap"]
    "maxHeap" : var.confluence_configuration["max_heap"]
    "cpu" : var.confluence_configuration["cpu"]
    "mem" : var.confluence_configuration["mem"]
  }

  rds_instance_name = format("atlas-%s-%s-db", var.environment_name, local.product_name)

  # if the domain wasn't provided we will start Bamboo with LoadBalancer service without ingress configuration
  use_domain          = length(var.ingress) == 1
  product_domain_name = local.use_domain ? "${local.product_name}.${var.ingress[0].ingress.domain}" : null
  # ingress settings for Bamboo service
  ingress_with_domain = yamlencode({
    ingress = {
      create = "true"
      host   = local.product_domain_name
    }
  })

  service_as_loadbalancer = yamlencode({
    confluence = {
      service = {
        type = "LoadBalancer"
      }
    }
    ingress = {
      https = false
    }
  })

  ingress_settings   = local.use_domain ? local.ingress_with_domain : local.service_as_loadbalancer
  storage_class_name = "efs-cs"

  license_settings = yamlencode({
    confluence = {
      license = {
        secretName = kubernetes_secret.license_secret.metadata[0].name
      }
    }
  })

}