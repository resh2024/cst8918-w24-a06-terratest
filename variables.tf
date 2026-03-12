# Define config variables
variable "labelPrefix" {
  type        = string
  description = "resh0004"
}

variable "region" {
  default = "canadacentral"
}

variable "admin_username" {
  type        = string
  default     = "azureadmin"
  description = "The username for the local user account on the VM."
}
