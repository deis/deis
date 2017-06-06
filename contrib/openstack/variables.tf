variable "instances" {
  default = "3"
}

variable "stack_name" {
  default = "deis"
}

variable "flavor" {
  default = "2"
}

variable "public_key_path" {
  description = "The path of the ssh pub key"
  default = "~/.ssh/id_rsa.pub"
}

variable "image" {
  description = "the image to use"
  default = "coreos-717.3.0"
}

variable "network_name" {
  description = "name of the internal network to use"
  default = "internal"
}

variable "floatingip_pool" {
  description = "name of the floating ip pool to use"
  default = "external"
}

variable "username" {
  description = "Your openstack username"
}

variable "password" {
  description = "Your openstack password"
}

variable "tenant" {
  description = "Your openstack tenant/project"
}

variable "auth_url" {
  description = "Your openstack auth URL"
}

variable "userdata" {
  description = "location of your user-data"
  default = "./contrib/coreos/user-data"
}
