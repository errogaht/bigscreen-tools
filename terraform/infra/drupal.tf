data "docker_registry_image" "drupal" {
  name = "drupal:rc-php8.1"
}
resource "docker_image" "drupal" {
  name = data.docker_registry_image.drupal.name
  pull_triggers = [data.docker_registry_image.drupal.sha256_digest]
}

resource "docker_container" "drupal" {
  count = 0
  image = docker_image.drupal.latest
  name  = "drupal"
  restart = "always"
  networks_advanced {
    name = "private"
  }
  volumes {
    host_path = "${var.srvHomeDir}/drupal/modules"
    container_path = "/var/www/html/modules"
  }
  volumes {
    host_path = "${var.srvHomeDir}/drupal/profiles"
    container_path = "/var/www/html/profiles"
  }
  sysctls = {
    "net.ipv4.ip_unprivileged_port_start"=0
  }
  user = "1000:1000"
  volumes {
    host_path = "${var.srvHomeDir}/drupal/sites"
    container_path = "/var/www/html/sites"
  }
  volumes {
    host_path = "${var.srvHomeDir}/drupal/themes"
    container_path = "/var/www/html/themes"
  }

  labels {
    label = "traefik.http.routers.drupal.rule"
    value = "Host(`drupal.${var.mainHost}`)"
  }
  labels {
    label = "traefik.http.routers.drupal.tls"
    value = "true"
  }
}