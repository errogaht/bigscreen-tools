resource "docker_image" "postgres" {
  name = "postgres:latest"
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
  env = ["POSTGRES_USER=root", "POSTGRES_PASSWORD=secret"]
}