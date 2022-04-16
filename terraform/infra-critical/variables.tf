variable "selectelToken" {
  type = string
}
variable "registryUsername" {
  type = string
  sensitive = true
}
variable "registryPassword" {
  type = string
  sensitive = true
}
variable "acmeEmail" {
  type = string
}
variable "dbPassword" {
  type = string
}
variable "sshHost" {
  type = string
  sensitive = true
}
variable "sshUser" {
  type = string
  sensitive = true
}
variable "vpnIP" {
  type = string
  sensitive = true
}
variable "registryHost" {
  type = string
}
variable "srvHomeDir" {
  type = string
}
variable "monicaSecret" {
  type = string
}
variable "monicaHashSalt" {
  type = string
}
variable "mainHost" {
  type = string
}
variable "sshPort" {
  sensitive = true
}