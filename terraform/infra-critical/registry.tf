resource "docker_image" "registry" {
  name = "registry:2"
}
resource "ssh_resource" "registry_init" {
  host         = var.sshHost
  port = var.sshPort
  user         = var.sshUser
  agent        = true
  timeout = "1m"
  commands = [
    "mkdir -p ${var.srvHomeDir}/registry/storage",
  ]
}
resource "docker_container" "registry" {
  depends_on = [ssh_resource.registry_init]
  image = docker_image.registry.latest
  name  = "registry"
  restart = "always"
  volumes {
    host_path = "${var.srvHomeDir}/registry"
    container_path = "/config"
    read_only = true
  }

  volumes {
    host_path = "${var.srvHomeDir}/registry/storage"
    container_path = "/var/lib/registry"
    read_only = false
  }

  env = [
    "REGISTRY_AUTH=htpasswd",
    "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm",
    "REGISTRY_AUTH_HTPASSWD_PATH=/config/htpasswd",
  ]

  labels {
    label = "traefik.http.routers.registry.rule"
    value = "Host(`${var.registryHost}`)"
  }
  labels {
    label = "traefik.http.routers.registry.tls"
    value = "true"
  }
}