module "network" {
  source                     = "terraform-aws-modules/vpc/aws"
  version                    = "5.8.1"
  ignore_nested_deprecations = true

  for_each = var.networks

  providers = {
    aws = aws.primary
  }

  cidr = each.value.cidr
  name = each.key

  tags = {
    Environment = "dev"
    Name        = each.key
  }

  depends_on = [module.base]
}

module "storage" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "4.1.0"

  count = var.enable_storage ? 1 : 0

  bucket = "example-storage"
}
