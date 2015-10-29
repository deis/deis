variable "instances" {
  default = "3"
}

variable "prefix" {
  default = "deis"
}

variable "region" {
  default = "sfo1"
}

variable "size" {
  default = "8GB"
}

variable "ssh_keys" {
  description = "The ssh fingerprint of the ssh key you'll be using"
}

variable "token" {
  description = "Your DigitalOcean auth token"
}
