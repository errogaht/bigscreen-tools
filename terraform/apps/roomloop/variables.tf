variable "selectelToken" {
  type = string
  description = "selectel DNS api token"
}
variable "registryUsername" {
  type = string
  sensitive = true
}
variable "registryPassword" {
  type = string
  sensitive = true
}

variable "sshHost" {
  type = string
  sensitive = true
}
variable "sshUser" {
  type = string
  sensitive = true
}
variable "registryHost" {
  type = string
  sensitive = true
}
variable "srvHomeDir" {
  type = string
}
variable "mainHost" {
  type = string
}
variable "sshPort" {
  sensitive = true
}