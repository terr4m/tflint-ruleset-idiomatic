data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"]

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

data "aws_caller_identity" "current" {
  provider = aws.primary

  depends_on = [aws_organizations_account.main]
}

data "aws_availability_zones" "available" {
  count = 1

  state = "available"
}
