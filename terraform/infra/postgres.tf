resource "docker_image" "postgres" {
  name = "postgres:14.2"
}

resource "docker_network" "web" {
  name = "web"
  internal = false

}

resource "docker_network" "internal" {
  name = "internal"
  internal = true
}
resource "docker_container" "postgres" {
  image = docker_image.postgres.latest
  name  = "postgres"
  restart = "always"
  ports {
    internal = 5432
    external = 5432
  }
  volumes {
    volume_name = "postgres"
    container_path = "/var/lib/postgresql/data"
  }
  volumes {
    volume_name = "${var.srvHomeDir}/pgsock"
    container_path = "/var/run/postgresql"
  }
  labels {
    label = "traefik.enable"
    value = "false"
  }
  env = ["POSTGRES_USER=root", "POSTGRES_PASSWORD=secret"]
}