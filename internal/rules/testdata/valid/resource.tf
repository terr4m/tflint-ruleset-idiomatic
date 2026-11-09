resource "aws_instance" "web" {
  count = 2

  provider = aws.primary

  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"
  subnet_id     = var.subnet_id

  tags = {
    Name = "web-${count.index}"
  }

  connection {
    host = self.public_ip
    type = "ssh"
  }

  provisioner "local-exec" {
    command = "echo created"
  }

  provisioner "local-exec" {
    command = "echo configured"
  }

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [tags]
  }

  depends_on = [aws_security_group.web]
}

resource "aws_s3_bucket" "logs" {
  for_each = toset(var.bucket_names)

  bucket = each.value
}

resource "aws_instance" "checked" {
  ami = "ami-0123456789abcdef0"

  lifecycle {
    precondition {
      condition     = var.enabled
      error_message = "The instance must be enabled."
    }
  }
}
