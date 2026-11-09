output "instance_ids" {
  description = "IDs of the web instances."
  sensitive   = false
  value       = aws_instance.web[*].id

  precondition {
    condition     = length(aws_instance.web) > 0
    error_message = "At least one instance is required."
  }

  depends_on = [aws_security_group.web]
}

output "token" {
  description = "The API token."
  ephemeral   = true
  deprecated  = "Read the token from the secret store instead."
  value       = var.token
}
