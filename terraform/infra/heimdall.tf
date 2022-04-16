data "docker_registry_image" "heimdall" {
  name = "linuxserver/heimdall"
}
resource "docker_image" "heimdall" {
  name = data.docker_registry_image.heimdall.name
  pull_triggers = [data.docker_registry_image.heimdall.sha256_digest]
}

resource "docker_container" "heimdall" {
  image = docker_image.heimdall.latest
  name  = "heimdall"
  restart = "always"
  networks_advanced {
    name = "private"
  }
  volumes {
    host_path = "${var.srvHomeDir}/heimdall"
    container_path = "/config"
  }

  env = [
    "PGID=1000",
    "PUID=1000",
  ]
  labels {
    label = "traefik.http.routers.heimdall.rule"
    value = "Host(`${var.mainHost}`) && ClientIP(`${var.vpnIP}/32`)"
  }
  labels {
    label = "traefik.http.routers.heimdall.tls"
    value = "true"
  }
}