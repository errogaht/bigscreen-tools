terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.16.0"
    }
    ssh = {
      source = "loafoe/ssh"
    }
    acme = {
      source  = "vancluever/acme"
      version = "~> 2.0"
    }
  }
}

provider "docker" {
  host     = "ssh://${var.sshUser}@${var.sshHost}:${var.sshPort}"
  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"]

  registry_auth {
    address  = var.registryHost
    username = var.registryUsername
    password = var.registryPassword
  }
}
resource "docker_network" "private" {
  name = "private"
  driver = "bridge"
  internal = false
}