locals {
  traefik_dynamic_conf = templatefile("traefik-dynamic-conf.yml", {
    certs = local.certs
  })
}
resource "ssh_resource" "traefik_init" {
  host     = var.sshHost
  port     = var.sshPort
  user     = var.sshUser
  agent    = true
  timeout  = "1m"
  commands = [
    "mkdir -p ${var.srvHomeDir}/registry/",
    "echo '${file("${path.module}/traefik.yml")}' > ${var.srvHomeDir}/traefik.yml",
    "echo '${local.traefik_dynamic_conf}' > ${var.srvHomeDir}/traefik-dynamic-conf/traefik.yml",
  ]
}
resource "docker_image" "traefik" {
  name = "traefik:v2.6"
}
resource "docker_container" "traefik" {
  depends_on = [ssh_resource.traefik_init]
  image      = docker_image.traefik.latest
  name       = "traefik"
  restart    = "always"
  networks_advanced {
    name = "bridge"
  }
  networks_advanced {
    name = docker_network.private.name
  }
  ports {
    internal = 80
    external = 80
  }
  ports {
    internal = 443
    external = 443
  }
  volumes {
    host_path      = "/var/run/docker.sock"
    container_path = "/var/run/docker.sock"
  }
  volumes {
    host_path      = "${var.srvHomeDir}/traefik.yml"
    container_path = "/etc/traefik/traefik.yml"
  }
  volumes {
    host_path      = "${var.srvHomeDir}/traefik-dynamic-conf"
    container_path = "/traefik-dynamic-conf"
  }

  volumes {
    host_path = "${var.srvHomeDir}/certs"
    container_path = "/certs"
  }

  env = [
    "SELECTEL_API_TOKEN=${var.selectelToken}",
  ]
  labels {
    label = "traefik.http.routers.traefik.rule"
    value = "Host(`traefik.${var.mainHost}`) && ClientIP(`${var.vpnIP}/32`)"
  }
  labels {
    label = "traefik.http.routers.traefik.service"
    value = "api@internal"
  }
  labels {
    label = "traefik.http.routers.traefik.entrypoints"
    value = "websecure"
  }
  labels {
    label = "traefik.http.routers.traefik.tls"
    value = "true"
  }
}

resource "docker_image" "whoami" {
  name = "registry.${var.mainHost}/whoami"
}
resource "docker_container" "whoami" {
  image = docker_image.whoami.latest
  name  = "whoami"
  labels {
    label = "traefik.http.routers.whoami.rule"
    value = "Host(`whoami.${var.mainHost}`)"
  }
  labels {
    label = "traefik.http.routers.whoami.tls"
    value = "true"
  }
}