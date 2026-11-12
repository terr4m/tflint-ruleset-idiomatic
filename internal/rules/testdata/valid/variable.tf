variable "environment" {
  description = "Deployment environment name."
  type        = string
  nullable    = false
  default     = "dev"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be dev, staging, or prod."
  }
}

variable "ports" {
  description = "Port configuration."

  type = object({
    external = number
    internal = number
    protocol = optional(string, "tcp")
  })

  sensitive = false

  default = {
    external = 8300
    internal = 8300
  }
}

variable "token" {
  description = "API token."
  type        = string
  nullable    = false
  sensitive   = true
  ephemeral   = true
}

variable "legacy_region" {
  description = "Deprecated region input."
  type        = string
  default     = "us-east-1"
  deprecated  = "Use var.region instead."
}
