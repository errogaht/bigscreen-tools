resource "docker_image" "roomloop" {
  name = "${var.registryHost}/roomloop"
}

resource "docker_container" "roomloop" {
  depends_on = [docker_container.postgres]
  image = docker_image.roomloop.latest
  name  = "roomloop"

  volumes {
    host_path = "${var.srvHomeDir}/app/.env"
    container_path = "/.env"
  }
  volumes {
    host_path = "${var.srvHomeDir}/pgsock"
    container_path = "${var.srvHomeDir}/pgsock"
  }
}