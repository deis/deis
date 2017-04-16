
variable "access_key" {}

variable "secret_key" {}

variable "nat_ami" {
    default = "ami-2dae821d"
}

variable "bastion_ami" {
    default = "ami-a9e2da99"
}

variable "region" {
    default = "us-east-1"
}

variable "availability_zone" {
    default = "us-east-1b"
}

variable "key_name" {}

variable "virt_type" {
	default = "pv"
}

variable "instance_type" {
	default = "m3.large"
}

variable "cluster_size" {
	default = "3"
}

variable "docker_volume_size" {
	default = "500"
}

