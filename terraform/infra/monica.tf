data "docker_registry_image" "monica" {
  name = "monica:3"
}
resource "docker_image" "monica" {
  name = data.docker_registry_image.monica.name
  pull_triggers = [data.docker_registry_image.monica.sha256_digest]
}

resource "docker_container" "monica" {
  image = docker_image.monica.latest
  name  = "monica"
  restart = "always"
  networks_advanced {
    name = "private"
  }
  volumes {
    host_path = "${var.srvHomeDir}/monica/storage"
    container_path = "/var/www/html/storage"
  }
  volumes {
    host_path = "${var.srvHomeDir}/mysqldsock"
    container_path = "/var/run/mysqld"
  }

  env = [
    "DB_HOST=mysql",
    "DB_PORT=3306",
    "DB_USERNAME=root",
    "DB_PASSWORD=${var.dbPassword}",
    "APP_ENV=production",
    "APP_URL=https://monica.${var.mainHost}",
    "APP_DEBUG=false",
    "APP_KEY=${var.monicaSecret}",
    "HASH_SALT=${var.monicaHashSalt}",
  ]


  labels {
    label = "traefik.http.routers.monica.rule"
    value = "Host(`monica.${var.mainHost}`)"
  }
  labels {
    label = "traefik.http.routers.monica.tls"
    value = "true"
  }
}