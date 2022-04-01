data "docker_registry_image" "roomloop" {
  name = "${var.registryHost}/roomloop"
}
resource "docker_image" "roomloop" {
  name = data.docker_registry_image.roomloop.name
  pull_triggers = [data.docker_registry_image.roomloop.sha256_digest]
}

resource "docker_container" "roomloop" {
  count = 1
  depends_on = [docker_container.postgres]
  image = docker_image.roomloop.latest
  name  = "roomloop"
  restart = "always"
  volumes {
    host_path = "${var.srvHomeDir}/app/.env"
    container_path = "/.env"
  }
  volumes {
    host_path = "${var.srvHomeDir}/pgsock"
    container_path = "${var.srvHomeDir}/pgsock"
  }
  labels {
    label = "traefik.enable"
    value = "false"
  }
}

resource "docker_container" "tgwebhook" {
  count = 1
  depends_on = [docker_container.postgres]
  image = docker_image.roomloop.latest
  name  = "tgwebhook"
  restart = "always"
  volumes {
    host_path = "${var.srvHomeDir}/app/.env"
    container_path = "/.env"
  }
  volumes {
    host_path = "${var.srvHomeDir}/pgsock"
    container_path = "${var.srvHomeDir}/pgsock"
  }
  entrypoint = ["/roomloop", "tghook"]

  labels {
    label = "traefik.http.routers.tghook.rule"
    value = "Host(`tgwebhook.${var.mainHost}`)"
  }
  labels {
    label = "traefik.http.routers.tghook.tls"
    value = "true"
  }
  labels {
    label = "traefik.http.routers.tghook.tls.certresolver"
    value = "default"
  }
  labels {
    label = "traefik.http.routers.tghook.tls.domains[0].main"
    value = var.mainHost
  }
  labels {
    label = "traefik.http.routers.tghook.tls.domains[0].sans"
    value = "*.${var.mainHost}"
  }
}