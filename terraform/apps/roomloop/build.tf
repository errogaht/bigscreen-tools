terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "2.16.0"
    }
    ssh = {
      source = "loafoe/ssh"
    }
    acme = {
      source  = "vancluever/acme"
      version = "~> 2.0"
    }
  }
}


resource "null_resource" "binary" {
  triggers = {
    time = timestamp()
  }
  provisioner "local-exec" {
    working_dir = "./../../../"
    command = "go build -o images/roomloop/roomloop app.go"
  }
}

resource "null_resource" "onchange" {
  depends_on = [null_resource.binary]
  triggers = {
    binary = filemd5("./../../../images/roomloop/roomloop")
    dockerfile = filemd5("./../../../images/roomloop/Dockerfile")
  }
  provisioner "local-exec" {
    working_dir = "./../../../"
    command = "docker build -t ${var.registryHost}/roomloop ./images/roomloop/ && docker push ${var.registryHost}/roomloop"
  }
}