#!/bin/bash

function check_already_configured() {
    if [ -f "./terraform.tfstate" ] ; then
        echo
        echo "WARNING: Terraform state already exists!"
        echo "----------------------------------------"
        echo
        echo "A terraform state file already exists.  If any resources are currently being"
        echo "managed by Terraform, running this script will destroy the existing Terraform"
        echo "state and possibly leave AWS resources unable to be managed by Terraform."
        echo
        echo "This means you will not be able to automatically remove these instances using"
        echo "Terraform."
        echo
        echo "You should inspect the current state using the following command:"
        echo
        echo "    terraform show"
        echo
        echo "If you want to destroy all managed resources, run:"
        echo
        echo "    terraform destroy -var-file=config.tfvars -var-file=credentials.tfvars"
        echo
        echo "And then delete the terraform.tfstate file before running this script"
        echo
        echo "IMPORTANT NOTE: You *MUST* generate a new discovery URL next time your run"
        echo "this script if you destroy an existing cluster."
        echo
        exit 1
    fi
}

function display_intro_text() {
    cat <<INTRO

Provisioning a Deis Cluster
---------------------------

This script will walk you through the configuration and provisioning of a Deis
cluster on AWS using Terraform (http://terraform.io).

The following software must be installed and configured before using this
script:

 * awscli: AWS command line utility
        http://aws.amazon.com/cli/

 * terraform: Infrastructure managment
        https://terraform.io/downloads.html

 * PyYaml: Python YAML library
        Run: pip install pyyaml

Alternatively, you can create a cluster using Deis Pro for free (AWS charges
will be billed to your AWS account by Amazon) at https://try.deis.com/

INTRO
}


function check_awscli() {
    which aws >/dev/null 2>&1 || {
        echo "ERROR: awscli toos are not installed or are not in your path."
        echo "You can download awscli from http://aws.amazon.com/cli/"
        exit 1
    }
    aws --version 2>&1 | grep '^aws-cli/.*' >/dev/null || {
        echo "ERROR: The binary named aws in your path does not appear to be"
        echo "a valid version of the awscli tools."
        echo "You can download awscli from http://aws.amazon.com/cli/"
        exit 1
    }
    [ -f ~/.aws/config ] || {
        echo "ERROR: you awscli configuration does not exist in ~/.aws/config"
        echo "Please run:"
        echo "    aws configure"
        exit 1
    }
}

function get_aws_creds() {
    AWS_PROFILE_LIST=$(echo "default" ; cat ~/.aws/config | grep '\[profile .*\]' | sed 's/\[profile \(.*\)\]/\1/' | sort)
    if [ "$AWS_PROFILE_LIST" != "default" ] ; then
        echo "Select an AWS CLI profile configuration"
        echo "---------------------------------------"
        for profile in $AWS_PROFILE_LIST ; do 
            echo $profile
        done
        AWS_PROFILE=
        while [ "$AWS_PROFILE" = "" ] ; do
            echo
            read -p "Enter a profile configuration name [default]: " AWS_PROFILE
            if [ -z "$AWS_PROFILE" ] ; then
                AWS_PROFILE="default"
            fi
            echo "$AWS_PROFILE_LIST" | grep "$AWS_PROFILE" >/dev/null 2>&1
            if [ $? != "0" ] ; then
                echo "$AWS_PROFILE is not a valid profile name in your configuration"
                AWS_PROFILE=""
            fi
        done
        AWS_PROFILE_OPT=" --profile $AWS_PROFILE "
        echo
    else
        AWS_PROFILE_OPT=""
        AWS_PROFILE="default"
    fi
    AWS_ACCESS_KEY_ID=$(aws configure get $AWS_PROFILE_OPT aws_access_key_id)
    AWS_SECRET_ACCESS_KEY=$(aws configure get $AWS_PROFILE_OPT aws_secret_access_key)
    AWS_DEFAULT_REGION=$(aws configure $AWS_PROFILE_OPT get region)
    
    failure=
    if [ -z "$AWS_ACCESS_KEY_ID" ] ; then
        echo "error: the access key is not specified for the '$AWS_PROFILE' profile"
        failure=yes
    fi
    if [ -z "$AWS_SECRET_ACCESS_KEY" ] ; then
        echo "error: the secret key is not specified for the '$AWS_PROFILE' profile"
        failure=yes
    fi
    if [ -z "$AWS_DEFAULT_REGION" ] ; then
        echo "error: the default region is not specified for the '$AWS_PROFILE' profile"
        failure=yes
    fi
    if [ "$failure" = "yes" ] ; then
        echo "Please run:"
        echo "   aws configure$AWS_PROFILE_OPT"
        exit 1
    fi
}

function select_region() {
    echo
    echo "Select an AWS Region"
    echo "--------------------"
    AWS_AVAILABLE_REGIONS=$(aws --output text --profile $AWS_PROFILE ec2 --query 'Regions[].[RegionName]' describe-regions | sort -r)
    for region in $AWS_AVAILABLE_REGIONS ; do
        echo $region
    done
    AWS_REGION=""
    while [ "$AWS_REGION" = "" ] ; do
        echo
        echo "Select the region in which you would like to build the Deis cluster."
        echo
        read -p "Enter a region [${AWS_DEFAULT_REGION}]: " AWS_REGION
        echo "$AWS_AVAILABLE_REGIONS" | grep "$AWS_REGION" >/dev/null 2>&1
        if [ -z "$AWS_REGION" ] ; then
            AWS_REGION=$AWS_DEFAULT_REGION
        fi
        if [ $? != "0" ] ; then
            echo "$AWS_REGION is not a valid region identifier"
            AWS_REGION=""
        fi
    done
    echo
}

function select_vpc_config() {
    cat <<VPC_TYPES
VPC Configuration Options
-------------------------
There are 3 options for the AWS VPC configuration for Deis nodes.  You can use
an existing VPC on your AWS account, or a new VPC with public or private Deis
cluster nodes can be created.

1. Existing VPC:  You will need to enter a VPC ID and a Subnet ID to specify
   where launch the Deis cluster nodes.  For simple evaluation purposes, this
   is not recommended to ensure that any existing resources are not
   modified.

2. New VPC with Public Nodes:  For simply evaluating Deis, this is the
   recommended option.  This will create a new VPC with public Internet facing
   Deis nodes with routable IP addresses.  Security groups will be created that
   limit access to the specific ports required for Deis.  AWS charges will
   apply for the number of Deis nodes you specify as well as 1 Elastic Load
   Balancer.

3. New VPC with Private Nodes:  AWS best practices recommend using a private
   subnet for instances whenever possible.  This option will create a private
   subnet for the Deis cluster nodes as well as a public facing NAT instance
   and bastion instance.  The Deis cluster nodes will not be directly
   accessible.  AWS charges will apply for the number of Deis nodes you
   specify, an Elastic Load Balancer and the two extra instances (NAT and
   Bastion)

VPC_TYPES


    VPC_TYPE=
    while [ -z "$VPC_TYPE" ] ; do
        read -p "Select a configuration [2]: " VPC_TYPE
        if [ -z "$VPC_TYPE" ] ; then
            VPC_TYPE=2
        fi
        case $VPC_TYPE in
            1)
                VPC_TYPE="existing"
                ;;
            2)
                VPC_TYPE="public"
                return
                ;;
            3)
                VPC_TYPE="private"
                return
                ;;
            *)
                VPC_TYPE=
                echo
                echo "Please enter 1, 2, or 3"
                echo
                ;;
        esac
    done
    VPC_LIST=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-vpcs --query 'Vpcs[].VpcId')
    if [ -z "$VPC_LIST" ] ; then
        echo
        echo "You do not have any existing VPCs in this region.  A new VPC must be created"
        echo "in order to continue."
        echo
        echo "1. VPC with Public Subnet for Deis Nodes"
        echo "2. VPC with Private Subnet for Deis Nodes (with a NAT and Bastion instance)"
        echo
        VPC_TYPE
        while [ -z "$VPC_TYPE" ] ; do
            read -p "Select a configuration [1]: " VPC_TYPE
            case $VPC_TYPE in
                1)
                    VPC_TYPE="public"
                    return
                    ;;
                2)
                    VPC_TYPE="private"
                    return
                    ;;
                *)
                    VPC_TYPE=
                    echo
                    echo "Please enter 1 or 2"
                    echo
                    ;;
            esac
        done
    fi
    echo
    echo "Select which VPC to use for Deis cluster nodes"
    echo "----------------------------------------------"
    echo
    for vpcs in $VPC_LIST ; do 
            echo "$vpcs"	
    done
    echo
    VPC_ID=
    while [ -z "$VPC_ID" ] ; do 
        read -p "VPC ID: " VPC_ID
        echo "$VPC_LIST" | grep "$VPC_ID" >/dev/null 2>&1
        if [ $? != 0 ] ; then
            echo
            echo "$VPC_ID is not a valid VPC ID"
            echo
            VPC_ID=""
        fi
    done
}

function select_keypair() {
    KEYPAIR_LIST=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-key-pairs --query 'KeyPairs[].KeyName' | sort)
    if [ -z "$KEYPAIR_LIST" ] ; then
        echo "No SSH keypairs currently exist in this region.  A new keypair must be created to"
        echo "continue."
        echo
        create_keypair
    else
        echo
        echo "Select an SSH key pair"
        echo "----------------------"
        echo
        for keyname in $KEYPAIR_LIST ; do
            echo $keyname
        done
        KEY_NAME=
        while [ -z "$KEY_NAME" ] ; do
            echo
            read -p "Select a key name or press enter to generate a new keypair: " KEY_NAME
            if [ -z "$KEY_NAME" ] ; then
                create_keypair
            else
                echo "$KEYPAIR_LIST" | grep "$KEY_NAME" >/dev/null 2>&1
                if [ $? != "0" ] ; then
                    echo "$KEY_NAME is not a valid key name"
                    KEY_NAME=""
                fi
            fi
        done
    fi
}

function create_keypair() {
    echo
    KEY_LOCATION="${HOME}/.ssh"
    if [ ! -d "$KEY_LOCATION" ] ; then
        echo "WARNING: ${HOME}/.ssh directory does not exist.  Key will be stored in current"
        echo "directory."
        KEY_LOCATION="./"
        echo
    fi
    KEY_NAME=
    while [ -z "$KEY_NAME" ] ; do
        read -p "Enter a keypair name [deis]: " KEY_NAME
        if [ -z "$KEY_NAME" ] ; then
            KEY_NAME="deis"
        fi
        echo "$KEY_NAME" | grep '[^-a-zA-Z0-9]' > /dev/null 2>&1
        if [ $? = 0 ] ; then
            echo "Only alphanumeric characters and dashes (-) are allowed."
            KEY_NAME=
        else
            if [ -f "${KEY_LOCATION}/${KEY_NAME}.key" ] ; then
                echo "${KEY_LOCATION}/${KEY_NAME}.key already exists, please choose a different"
                echo "keypair name."
                echo
                KEY_NAME=
            else
                PRIVATE_KEY=$(aws --output json --region $AWS_REGION --profile $AWS_PROFILE ec2 create-key-pair --key-name "$KEY_NAME" --query 'KeyMaterial' | sed 's/"//g')
                if [ -z "$PRIVATE_KEY" ] ; then
                    KEY_NAME=
                    echo
                    echo "ERROR: could not create new key."
                fi
                echo -e "$PRIVATE_KEY" > "${KEY_LOCATION}/${KEY_NAME}.key"
            fi
        fi
    done
}

function get_bastion_ami() {
    DEFAULT_BASTION_IMAGE_ID=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-images --owners "099720109477" --filters Name="name",Values='ubuntu/images/hvm/ubuntu-trusty-14.04-amd64-server*' --query 'sort_by(Images, &CreationDate)[-1].ImageId')
    echo
    echo "Select a Bastion Instance Image"
    echo "-------------------------------"
    echo
    echo "The default AMI below is the latest EBS backed Ubuntu 14.04 image"
    echo
    aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-images --image-ids "$DEFAULT_BASTION_IMAGE_ID" --query 'sort_by(Images, &CreationDate)[-1].Name'
    echo 
    BASTION_IMAGE_ID=
    while [ -z "$BASTION_IMAGE_ID" ] ; do
        read -p "Enter an AMI ID for the bastion instance [${DEFAULT_BASTION_IMAGE_ID}]: " BASTION_IMAGE_ID
        if [ -z "$BASTION_IMAGE_ID" ] ; then
            BASTION_IMAGE_ID=$DEFAULT_BASTION_IMAGE_ID
        fi
        aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-images --image-ids "$BASTION_IMAGE_ID" > /dev/null 2>&1
        if [ $? != 0 ] ; then
            BASTION_IMAGE_ID=
            echo
            echo "error: the specified image ID does not exist"
        fi
    done
    echo
}

function get_existing_subnet() {
    echo
    echo "Select a subnet for running Deis nodes"
    echo "--------------------------------------"
    AWS_VPC_SUBNET_LIST=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-subnets --query 'Subnets[?VpcId==`'$VPC_ID'`].SubnetId')
    aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-subnets --query 'Subnets[?VpcId==`'$VPC_ID'`] | [?State==`available`].[SubnetId,CidrBlock,AvailabilityZone]'
    echo
    SUBNET_ID=
    while [ "$SUBNET_ID" = "" ] ; do
        echo
        read -p "Select a subnet ID for cluster nodes: " SUBNET_ID
        echo "$AWS_VPC_SUBNET_LIST" | grep "$SUBNET_ID" >/dev/null 2>&1
        if [ $? != "0" ] ; then
            echo "$SUBNET_ID is not a valid subnet id"
            SUBNET_ID=""
        fi
    done
    AWS_AZ=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-subnets --subnet-ids "$SUBNET_ID" --query 'Subnets[].AvailabilityZone')
    ROUTE_TABLE_ID=$(aws --output text --region $AWS_REGION ec2 describe-route-tables --filters Name="vpc-id",Values="$VPC_ID" --query 'RouteTables[].Associations[?SubnetId==`'$SUBNET_ID'`].RouteTableId')
    if [ -z "$ROUTE_TABLE_ID" ] ; then
        # Subnet uses the main routing table
        ROUTE_TABLE_ID=$(aws --output text --region $AWS_REGION ec2 describe-route-tables --filter Name="vpc-id",Values="$VPC_ID" --query 'RouteTables[].Associations[?Main==`true`].RouteTableId')
    fi
    DESTINATION=$(aws --output text --region $AWS_REGION ec2 describe-route-tables --route-table-ids "$ROUTE_TABLE_ID" --query 'RouteTables[].Routes[?DestinationCidrBlock==`0.0.0.0/0`].[InstanceId, GatewayId]')
    if [ -z "$DESTINATION" ] ; then
        echo
        echo " * Warning: unable to determine if the specified subnet is public or private"
        echo " * since the routing table associated with the subnet does not have a default"
        echo " * gateway."
        echo " *"
        echo " * Assuming a public subnet configuration"
        echo
        VPC_TYPE=public
        read -p "Press ENTER to continue" THROWAWAY
    else
        VPC_TYPE=
        if [[ "$DESTINATION" =~ i- ]] ; then
            # Default gateway is an instance
            VPC_TYPE="existing-private"
        elif [[ "$DESTINATION" =~ (i|c)gw- ]] ; then
            # Internet gateway or customer gateway (VPN).  Assume a public style
            # configuration for subnets with a VPN default route 
            VPC_TYPE="existing-public"
        fi

        if [ "$VPC_TYPE" = "existing-private" ] ; then
            echo
            echo "The subnet you selected is behind a NAT instance.  Public addresses will"
            echo "not be assigned."
            echo
        elif [ "$VPC_TYPE" = "existing-public" ] ; then
            echo
            echo "The subnet you selected is directly connected to the Internet.  Public"
            echo "addresses will be assigned to each cluster node."
            echo
        else
            echo
            echo "ERROR: a default route for the routing table associated with the subnet"
            echo "exists, but could not determine it's destination type (IGW/NAT)" 
            echo
            exit 1
        fi
    fi
}

function select_az() {
    echo
    echo "Availability Zones in $AWS_REGION"
    echo "---------------------------------"
    AWS_REGIONS_AZ=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-availability-zones --query 'AvailabilityZones[?State==`available`].[ZoneName]' | sort)
    AWS_FIRST_AZ=""
    for az in $AWS_REGIONS_AZ ; do
        if [ -z "$AWS_FIRST_AZ" ] ; then
            AWS_FIRST_AZ=$az
        fi
        echo $az
    done
    AWS_AZ=""
    while [ "$AWS_AZ" = "" ] ; do
        echo
        read -p "Enter an availability zone [${AWS_FIRST_AZ}]: " AWS_AZ
        if [ -z "$AWS_AZ" ] ; then
            AWS_AZ=$AWS_FIRST_AZ
        fi
        echo "$AWS_REGIONS_AZ" | grep "$AWS_AZ" >/dev/null 2>&1
        if [ $? != "0" ] ; then
            echo "$AWS_AZ is not a valid availability zone"
            AWS_AZ=""
        fi
    done
    echo
}

function select_instance_type() {
    echo
    echo "Select an instance type for Deis cluster nodes"
    echo "----------------------------------------------"
    echo
    echo "The minimum required memory for a Deis cluster node is 4GB.  However"
    echo "the recommended minimum is 8GB of memory."
    echo 
    echo "The recommended minimum instance types are m4.large, m3.large, c4.xlarge,"
    echo "c3.2xlarge, or r3.2xlarge"
    echo
    AWS_INSTANCE_TYPE=
    read -p "Select an instance type [m3.large]: " AWS_INSTANCE_TYPE
    if [ -z "$AWS_INSTANCE_TYPE" ] ; then
        AWS_INSTANCE_TYPE="m3.large"
    fi
}

function select_node_count() {
    echo
    echo "Select the number of Deis cluster nodes"
    echo "---------------------------------------"
    echo 
    echo "The minimum number of cluster node members is 3.  Deis does not require an"
    echo "an odd number of nodes.  If an even number is specified, the last node will"
    echo "be used as a standby for all quorum decisions."
    echo
    CLUSTER_SIZE=
    while [ -z "$CLUSTER_SIZE" ] ; do
        read -p "Cluster size [3]: " CLUSTER_SIZE
        if [ -z "$CLUSTER_SIZE" ] ; then
            CLUSTER_SIZE=3
        fi
        if ! [ "$CLUSTER_SIZE" -eq "$CLUSTER_SIZE" ] 2>/dev/null || [ "$CLUSTER_SIZE" -lt 3 ]; then
            echo
            echo "A number is greater than 3 is required."
            echo
            CLUSTER_SIZE=""
        fi
    done
}

function select_node_management() {
    echo
    echo "Select Node Management Configuration"
    echo "------------------------------------"
    echo
    echo "There are two options for managing Deis nodes.  Nodes can be managed directly"
    echo "by Terraform or with an AWS Auto Scaling Group."
    echo
    echo "1. Terraform Managed Deis Nodes"
    echo "        Nodes managed directly by Terraform allow for better infrastructure"
    echo "        visibility by Terraform.  Each node will be tracked directly in the"
    echo "        Terraform state file and scaling up cluster size is done via the"
    echo "        Terraform configuration.  However, if a cluster node dies, it will"
    echo "        not automatically be restarted until a Terraform apply is done."
    echo
    echo "2. AWS Auto Scaling Group"
    echo "        An Auto Scaling Group will enforce the configured number of nodes at"
    echo "        all times.  If a node dies, the ASG will automatically start a new"
    echo "        node to replace the dead instance immediately.  However, Terraform"
    echo "        cannot directly manage individual nodes because the instance"
    echo "        creation is managed by the ASG.  Terraform will not have any"
    echo "        information about Deis nodes other than the number of nodes configured"
    echo "        in the Auto Scaling Group"
    echo
    NODE_MANAGEMENT_TYPE=
    while [ -z "$NODE_MANAGEMENT_TYPE" ] ; do
        read -p "Select a node management option: [1]: " NODE_MANAGEMENT_TYPE
        if [ -z "$NODE_MANAGEMENT_TYPE" ] ; then
            NODE_MANAGEMENT_TYPE=1
        fi
        case $NODE_MANAGEMENT_TYPE in
            1) 
                NODE_MANAGEMENT_TYPE=tf
                ;;
            2) 
                NODE_MANAGEMENT_TYPE=asg
                ;;
            *)
                NODE_MANAGEMENT_TYPE=
                echo
                echo "Please select 1 Terraform managed nodes or 2 for AutoScalingGroup"
                echo
                ;;
        esac
    done

}

function enable_advanced_options() {
    echo "Enable Advanced Options"
    echo "-----------------------"
    echo
    echo "Most basic configurations should be suitable for most users for Deis evaluation"
    echo "and simple uses; however, you can further customize your cluster."
    echo
    ADVANCED_OPTIONS=
    while [ -z "$ADVANCED_OPTIONS" ] ; do
        read -p "Would you like to enable advanced options [y/N]: " ADVANCED_OPTIONS
        if [ -z "$ADVANCED_OPTIONS" ] ; then
            ADVANCED_OPTIONS=n
            return
        fi
        if [ "$ADVANCED_OPTIONS" = "y" -o "$ADVANCED_OPTIONS" = "Y" ] ; then
            ADVANCED_OPTIONS=y
            return
        fi
        if [ "$ADVANCED_OPTIONS" = "n" -a "$ADVANCED_OPTIONS" = "N" ] ; then
            ADVANCED_OPTIONS=n
            return
        fi
    done
}

function select_virt_type() {
    INSTANCE_TYPE_FAMILY=$(echo $AWS_INSTANCE_TYPE | sed 's/^\([^.]*\)\..*$/\1/')
    if [ "$INSTANCE_TYPE_FAMILY" != "m3" -a "$INSTANCE_TYPE_FAMILY" != "c3" ] ; then
        VIRT_TYPE="hvm"
        return
    fi
    echo
    echo "Select Virtualization Type"
    echo "--------------------------"
    echo
    echo "1. paravirtual (recommended)"
    echo "2. hvm"
    echo
    VIRT_TYPE=
    while [ -z $VIRT_TYPE ] ; do
        read -p "Virtualization type [1]: " VIRT_TYPE
        if [ -z "$VIRT_TYPE" ] ; then
            VIRT_TYPE=1
        fi
        case $VIRT_TYPE in
            1) 
                VIRT_TYPE=paravirtual
                ;;
            2) 
                VIRT_TYPE=hvm
                ;;
            *)
                VIRT_TYPE=
                echo
                echo "Please select 1 for paravirtual or 2 for hvm"
                echo
                ;;
        esac
    done
}

function check_coreos_userdata() {
    DISCOVERY_URL=
    echo
    echo "Configure an etcd discovery URL"
    echo "-------------------------------"
    if [ -f "../coreos/user-data" ] ; then
        DISCOVERY_URL=$(grep -e "^[ ]\+discovery" ../coreos/user-data | sed 's/^[ ]*discovery: //')
        echo "An existing etcd discovery URL was found in the ../coreos/user-data file."
        echo "The URL is: $DISCOVERY_URL"
        echo
        echo "To prevent any new instances from interfering with an existing Deis cluster, it"
        echo "is recommended to generate a new URL."
        echo
        while [ -z "$GENERATE_ETCD_URL" ] ; do
            read -p "Generate a new etcd URL now [Y/n]: " GENERATE_ETCD_URL
            if [ -z "$GENERATE_ETCD_URL" ] ; then
                GENERATE_ETCD_URL=y
            fi
            if [ "$GENERATE_ETCD_URL" = "y" -o "$GENERATE_ETCD_URL" = "Y" ] ; then
                GENERATE_ETCD_URL=y
            elif [ "$GENERATE_ETCD_URL" != "n" -a "$GENERATE_ETCD_URL" != "N" ] ; then
                GENERATE_ETCD_URL=
            fi
        done
    fi
    if [ -z "$DISCOVERY_URL" -o "$GENERATE_ETCD_URL" = "y" ] ; then
        echo -n "Generating etcd discovery URL... "
        $(cd ../.. && make discovery-url > /dev/null)
        echo "done"
        echo
    fi
    echo -n "Generating user-data for instances..."
    python ./gen-userdata.py > coreos-user-data.txt
    echo "done"
}

function get_docker_volume_size() {
    DISCOVERY_URL=
    echo
    echo "Configure Docker EBS Volume Size"
    echo "--------------------------------"
    echo
    DOCKER_VOLUME_SIZE=
    while [ -z "$DOCKER_VOLUME_SIZE" ] ; do
        read -p "Enter the size (in GB) to host your Docker containers [100]: " DOCKER_VOLUME_SIZE
        if [ -z "$DOCKER_VOLUME_SIZE" ] ; then
            DOCKER_VOLUME_SIZE=100
        fi
        echo $DOCKER_VOLUME_SIZE | grep '^[0-9]\+$' > /dev/null
        if [ $? != 0 ] ; then
            DOCKER_VOLUME_SIZE=
            echo
        fi
    done
}


function write_terraform_config() {
    NAT_IMAGE_ID=$(aws --output text --region $AWS_REGION --profile $AWS_PROFILE ec2 describe-images --filter Name="owner-alias",Values="amazon" --filter Name="name",Values="amzn-ami-vpc-nat*" Name="virtualization-type",Values="paravirtual" --query 'sort_by(Images, &CreationDate)[-1].ImageId')
    echo 'region = "'$AWS_REGION'"' > config.tfvars
    echo 'key_name = "'$KEY_NAME'"' >> config.tfvars
    echo 'nat_ami = "'$NAT_IMAGE_ID'"' >> config.tfvars
    echo 'bastion_ami = "'$BASTION_IMAGE_ID'"' >> config.tfvars
    echo 'cluster_size = "'$CLUSTER_SIZE'"' >> config.tfvars
    echo 'availability_zone = "'$AWS_AZ'"' >> config.tfvars
    echo 'docker_volume_size = "'$DOCKER_VOLUME_SIZE'"' >> config.tfvars

    echo 'access_key = "'$AWS_ACCESS_KEY_ID'"' > credentials.tfvars
    echo 'secret_key = "'$AWS_SECRET_ACCESS_KEY'"' >> credentials.tfvars
    

    cp tflib/aws/coreos-images.tf tflib/aws/elb.tf \
       tflib/aws/provider.tf tflib/aws/securitygroups.tf \
       tflib/aws/variables.tf tflib/aws/vpc.tf .
    if [ "$NODE_MANAGEMENT_TYPE" = "asg" ] ; then
        cp tflib/aws/nodes-asg.tf .
    else
        cp tflib/aws/nodes-tf.tf .
        cp tflib/aws/elb-tf_override.tf .
    fi

    if [ "$VPC_TYPE" = "private" ] ; then
        cp tflib/aws/elb-private_override.tf tflib/aws/vpc-private.tf .
        if [ "$NODE_MANAGEMENT_TYPE" = "asg" ] ; then
            cp tflib/aws/nodes-asg-private_override.tf .
        else
            cp tflib/aws/nodes-tf-private_override.tf .
        fi
            
    else 
        rm -f elb-private_override.tf vpc-private.tf nodes-asg-private_override.tf \
            nodes-tf-private_override.tf
    fi
}

function apply_terraform() {
    echo
    echo "Configuration Complete"
    echo "----------------------"
    echo
    echo "The configuration is complete and the terraform configuration has been copied"
    echo "to `pwd`"
    echo
    echo "You can now use the standard terraform commands (plan/apply/show) in this directory."
    echo
    echo "To see the plan, run:"
    echo "    terraform plan -var-file=config.tfvars -var-file=credentials.tfvars"
    echo
    echo "To apply the configuration and provision the cluster, run:"
    echo "    terraform apply -var-file=config.tfvars -var-file=credentials.tfvars"
    echo
    PROVISION_CLUSTER=
    while [ -z "$PROVISION_CLUSTER" ] ; do
        read -p "Would you like to provision the cluster now [Y/n]: " PROVISION_CLUSTER
        if [ -z "$PROVISION_CLUSTER" ] ; then
            PROVISION_CLUSTER=y
        fi
        if [ "$PROVISION_CLUSTER" = "y" -o "$PROVISION_CLUSTER" = "Y" ] ; then
            terraform apply -var-file=config.tfvars -var-file=credentials.tfvars
            exit
        fi
        if [ "$PROVISION_CLUSTER" != "n" -a "$PROVISION_CLUSTER" != "N" ] ; then
            PROVISION_CLUSTER=
        fi
    done
}

check_already_configured
display_intro_text
check_awscli
get_aws_creds
enable_advanced_options
select_region
if [ "$ADVANCED_OPTIONS" = "n" ] ; then
    select_az
    VPC_ID=
    VPC_TYPE=public
else
    select_vpc_config
    if [ -z "$VPC_ID" ] ; then
        if [ "$VPC_TYPE" = "private" ] ; then
            get_bastion_ami
        fi
        select_az
    else
        get_existing_subnet
    fi
fi
select_instance_type
select_keypair
select_virt_type
check_coreos_userdata
get_docker_volume_size
select_node_count
if [ "$ADVANCED_OPTIONS" = "y" ] ; then
    select_node_management
else
    NODE_MANAGEMENT_TYPE=tf
fi
write_terraform_config
apply_terraform

