locals {
  certs = toset([
    "tgwebhook.${var.mainHost}",
    "tgwebhook2.${var.mainHost}",
    "traefik.${var.mainHost}",
    "monica.${var.mainHost}",
    "whoami.${var.mainHost}",
    var.registryHost,
    var.mainHost,
  ])
}
provider "acme" {
  server_url = "https://acme-v02.api.letsencrypt.org/directory"
}

resource "tls_private_key" "acme_private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.acme_private_key.private_key_pem
  email_address   = var.acmeEmail
}

resource "acme_certificate" "this" {
  for_each = local.certs
  account_key_pem           = acme_registration.reg.account_key_pem
  common_name               = each.value

  dns_challenge {
    provider = "selectel"
    config = {
      "SELECTEL_API_TOKEN" = var.selectelToken
    }
  }
}

/*resource "acme_certificate" "wildcard" {
  account_key_pem           = acme_registration.reg.account_key_pem
  common_name               = var.mainHost
  subject_alternative_names = ["*.${var.mainHost}"]

  dns_challenge {
    provider = "selectel"
    config = {
      "SELECTEL_API_TOKEN" = var.selectelToken
    }
  }
}*/

resource "ssh_resource" "certs" {
  for_each = local.certs
  host     = var.sshHost
  port     = var.sshPort
  user     = var.sshUser
  agent    = true
  timeout  = "1m"
  commands = [
    "echo '${acme_certificate.this[each.value].certificate_pem}${acme_certificate.this[each.value].issuer_pem}' > ${var.srvHomeDir}/certs/${each.value}.pem",
    "echo '${acme_certificate.this[each.value].private_key_pem}' > ${var.srvHomeDir}/certs/${each.value}.key",
  ]
}