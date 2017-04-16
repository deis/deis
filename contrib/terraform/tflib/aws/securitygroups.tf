
resource aws_security_group "deis_ssh" {
    name = "deis-allow-ssh"
    description = "Allow SSH - Managed by Deis Terraform"
    vpc_id = "${aws_vpc.deis.id}"

    ingress {
        from_port = 22
        to_port = 22
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    tags {
        Name = "deis-allow-ssh"
    }
}

resource aws_security_group "deis_control" {
    name = "deis-allow-control-port"
    description = "Allow the port required for deisctl - Managed by Deis Terraform"
    vpc_id = "${aws_vpc.deis.id}"

    ingress {
        from_port = 2222
        to_port = 2222
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    tags {
        Name = "deis-allow-control-port"
    }
}

resource aws_security_group "deis_web" {
    name = "deis-allow-web"
    description = "Allow the ports required for http and https - Managed by Deis Terraform"
    vpc_id = "${aws_vpc.deis.id}"

    ingress {
        from_port = 80
        to_port = 80
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port = 443
        to_port = 443
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }
    tags {
        Name = "deis-allow-web"
    }
}

resource aws_security_group "intra_vpc" {
    name = "deis-allow-internal-vpc"
    description = "Allow all cluster nodes to talk to one another - Managed by Deis Terraform"
    vpc_id = "${aws_vpc.deis.id}"

    ingress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        self = true
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

    tags {
        Name = "deis-allow-internal-vpc"
    }
}
