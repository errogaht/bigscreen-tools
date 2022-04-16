resource "docker_image" "mysql" {
  name = "mysql:5.7"
}

resource "docker_container" "mysql" {
  image = docker_image.mysql.latest
  name  = "mysql"
  restart = "always"
  networks_advanced {
    name = docker_network.private.name
  }
  ports {
    internal = 3306
    external = 3306
  }
  volumes {
    volume_name = "mysql"
    container_path = "/var/lib/mysql"
  }
  volumes {
    host_path = "${var.srvHomeDir}/mysqldsock"
    container_path = "/var/run/mysqld"
  }
  labels {
    label = "traefik.enable"
    value = "false"
  }
  env = ["MYSQL_ROOT_PASSWORD=${var.dbPassword}"]
}