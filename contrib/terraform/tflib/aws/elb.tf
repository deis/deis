
resource "aws_elb" "deis_web_elb" {
    name = "deis-elb"
    subnets = ["${aws_subnet.deis_public.id}"]

    listener {
        instance_port = 80
        instance_protocol = "http"
        lb_port = 80
        lb_protocol = "http"
    }

    listener {
        instance_port = 443
        instance_protocol = "tcp"
        lb_port = 443
        lb_protocol = "tcp"
    }

    listener {
        instance_port = 2222
        instance_protocol = "tcp"
        lb_port = 2222
        lb_protocol = "tcp"
    }

    health_check {
        healthy_threshold = 4
        unhealthy_threshold = 2
        timeout = 5
        target = "HTTP:80/health-check"
        interval = 15
    }
}

output "elb-dns-name" {
    value = "${aws_elb.deis_web_elb.dns_name}"
}
