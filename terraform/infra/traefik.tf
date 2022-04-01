resource "ssh_resource" "traefik_init" {
  host     = var.sshHost
  port     = var.sshPort
  user     = var.sshUser
  agent    = true
  timeout  = "1m"
  commands = [
    "mkdir -p ${var.srvHomeDir}/registry/",
    "echo '${file("${path.module}/traefik.yml")}' > ${var.srvHomeDir}/traefik.yml",

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
  ports {
    internal = 8080
    external = 8080
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

  env = [
    "SELECTEL_API_TOKEN=${var.selectelToken}",
  ]
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
  labels {
    label = "traefik.http.routers.whoami.tls.certresolver"
    value = "default"
  }
  labels {
    label = "traefik.http.routers.whoami.tls.domains[0].main"
    value = var.mainHost
  }
  labels {
    label = "traefik.http.routers.whoami.tls.domains[0].sans"
    value = "*.${var.mainHost}"
  }
}