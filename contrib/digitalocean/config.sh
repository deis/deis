export DEIS_TEST_DOMAIN="xip.io"
export DO_TOKEN
export DO_SSH_FINGERPRINT

prompt "What Digital Ocean token should I use?" DO_TOKEN
prompt "What Digital Ocean ssh fingerprint should I use?" DO_SSH_FINGERPRINT

export TF_VAR_deis_root="${DEIS_ROOT}"
export TF_VAR_token="${DO_TOKEN}"
export TF_VAR_ssh_keys="${DO_SSH_FINGERPRINT}"
export TF_VAR_prefix="deis-${DEIS_ID}"

rigger-save-vars DEIS_TEST_DOMAIN \
                 DO_TOKEN \
                 DO_SSH_FINGERPRINT \
                 TF_VAR_deis_root \
                 TF_VAR_prefix \
                 TF_VAR_ssh_keys \
                 TF_VAR_token
