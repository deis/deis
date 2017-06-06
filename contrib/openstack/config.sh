export DEIS_TEST_DOMAIN="xip.io"
export TF_VAR_stack_name="deis-${DEIS_ID}"

prompt "What glance image should I use?" TF_VAR_image coreos-stable
prompt "What flavor ?" TF_VAR_flavor 2
prompt "Path to your public key:" TF_VAR_public_key_path "${HOME}/.ssh/id_rsa.pub"
prompt "What is the name of your neutron network?" TF_VAR_network_name internal
prompt "What is the name of your floating IP pool?" TF_VAR_floatingip_pool external
prompt "what is your openstack auth URL?" TF_VAR_auth_url ${OS_AUTH_URL}
prompt "what is your openstack username?" TF_VAR_username ${OS_USERNAME}
prompt "what is your openstack password?" TF_VAR_password ${OS_PASSWORD}
prompt "what is your openstack tenant name?" TF_VAR_tenant_name ${OS_TENANT_NAME}

export TF_VAR_deis_root="${DEIS_ROOT}"
export TF_VAR_image="${GLANCE_IMAGE}"
export TF_VAR_flavor="${FLAVOR}"
export TF_VAR_public_key_path="${KEY}"
export TF_VAR_network_name="${NETWORK}"
export TF_VAR_floatingip_pool="${POOL}"
export TF_VAR_auth_url="${OS_AUTH_URL}"
export TF_VAR_username="${OS_USERNAME}"
export TF_VAR_password="${OS_PASSWORD}"
export TF_VAR_tenant_name="${OS_TENANT_NAME}"

rigger-save-vars DEIS_TEST_DOMAIN \
                 TF_VAR_deis_root \
                 TF_VAR_image \
                 TF_VAR_flavor \
                 TF_VAR_public_key_path \
                 TF_VAR_network_name \
                 TF_VAR_floatingip_pool \
                 TF_VAR_auth_url \
                 TF_VAR_username \
                 TF_VAR_password \
                 TF_VAR_tenant_name \
                 TF_VAR_stack_name


