Provision a Deis Controller on Amazon EC2
=========================================

1. Install [knife-ec2][knifec2] with `gem install knife-ec2` or just
`bundle install` from the root directory of your deis repository:
```console
$ cd $HOME/projects/deis
$ gem install knife-ec2
Fetching: knife-ec2-0.6.4.gem (100%)
Successfully installed knife-ec2-0.6.4
1 gem installed
Installing ri documentation for knife-ec2-0.6.4...
Installing RDoc documentation for knife-ec2-0.6.4...
```

2. Export your EC2 credentials as environment variables and edit knife.rb
to read them:
```console
$ cat <<'EOF' >> $HOME/.bash_profile
export AWS_ACCESS_KEY=<your_aws_access_key>
export AWS_SECRET_KEY=<your_aws_secret_key>
EOF
$ source $HOME/.bash_profile
$ cat <<'EOF' >> $HOME/.chef/knife.rb
knife[:aws_access_key_id] = "#{ENV['AWS_ACCESS_KEY']}"
knife[:aws_secret_access_key] = "#{ENV['AWS_SECRET_KEY']}"
EOF
$ knife ec2 server list
Instance ID  Name  Public IP  Private IP  Flavor  Image  SSH Key  Security Groups  State
```

3. Download and install the [EC2 Command Line Tools][ec2cli] as described in
[AWS' documentation][ec2cli] and ensure they are available in your $PATH:
```console
$ ec2-describe-group
GROUP	sg-33d1045a	693041077886	default	default group
PERMISSION	693041077886	default	ALLOWS	tcp	0	65535	FROM	USER	693041077886	NAME default	ID sg-33d1045a	ingress
PERMISSION	693041077886	default	ALLOWS	udp	0	65535	FROM	USER	693041077886	NAME default	ID sg-33d1045a	ingress
PERMISSION	693041077886	default	ALLOWS	icmp	-1	-1	FROM	USER	693041077886	NAME default	ID sg-33d1045a	ingress
```

4. Run the provisioning script to create a new Deis controller:
```console
$ ./contrib/ec2/provision-ec2-controller.sh us-west-2
Creating security group: deis-controller
+ ec2-create-group deis-controller -d 'Created by Deis'
GROUP	sg-3c3a1c0c	deis-controller	Created by Deis
+ set +x
Authorizing TCP ports 22,80,443,514 from 0.0.0.0/0...
+ ec2-authorize deis-controller -P tcp -p 22 -s 0.0.0.0/0
...
ec2-203.0.113.33.us-west-2.compute.amazonaws.com
ec2-203-0-113-33.us-west-2.compute.amazonaws.com Chef Client finished, 74 resources updated
...
Instance ID: i-31c8d106
Flavor: m1.large
Image: ami-72e27c42
Region: us-west-2
Public DNS Name: ec2-203-0-113-33.us-west-2.compute.amazonaws.com
Public IP Address: 203.0.113.33
Run List: recipe[deis::controller]
...
```

[knifec2]: http://docs.opscode.com/plugin_knife_ec2.html
