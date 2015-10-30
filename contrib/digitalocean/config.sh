export DEIS_TEST_DOMAIN="xip.io"
export DO_TOKEN
export DO_SSH_FINGERPRINT

function check-do-token {
  if [ ! -z ${1} ]; then
    local token="${1}"
    curl --fail -X GET -H "Authorization: Bearer ${token}" \
                          "https://api.digitalocean.com/v2/account" &> /dev/null
  else
    return 1
  fi
}

function check-do-ssh-key {
  local token="${1}"
  local ssh_fingerprint="${2}"

  if [ ! -z ${ssh_fingerprint} ]; then
    curl --fail \
         -X GET \
         -H \
         "Authorization: Bearer ${token}" \
         "https://api.digitalocean.com/v2/account/keys/${ssh_fingerprint}" &> /dev/null
  else
    return 1
  fi
}

while true; do
  password-prompt "What DigitalOcean token should I use?" DO_TOKEN

  if ! check-do-token "${DO_TOKEN:-}"; then
    rerun_log error "Couldn't login to DigitalOcean using this API token. :-("
    unset DO_TOKEN
  else
    rerun_log info "Successfully logged into DigitalOcean!"
    break
  fi
done

while true; do
  ssh-private-key-prompt "What private SSH key should I use when creating DigitalOcean droplets?" SSH_PRIVATE_KEY_FILE

  export DO_SSH_FINGERPRINT="$(ssh-fingerprint "${SSH_PRIVATE_KEY_FILE}")"

  if ! check-do-ssh-key "${DO_TOKEN}" "${DO_SSH_FINGERPRINT:-}"; then
    rerun_log error "Couldn't find the fingerprint for this key in DigitalOcean."

    cat <<EOF
  Upload the public key by pressing the "Add SSH Key" button on
  your DigitalOcean security page:

    https://cloud.digitalocean.com/settings/security

  Or pick a different key...

EOF

    unset SSH_PRIVATE_KEY_FILE
  else
    rerun_log info "This SSH key is correctly configured for use with DigitalOcean!"
    break
  fi
done

rigger-log "DO_SSH_FINGERPRINT set to ${DO_SSH_FINGERPRINT}"

export TF_VAR_deis_root="${DEIS_ROOT}"
export TF_VAR_ssh_keys="${DO_SSH_FINGERPRINT}"
export TF_VAR_prefix="deis-${DEIS_ID}"

rigger-save-vars DEIS_TEST_DOMAIN \
                 DO_SSH_FINGERPRINT \
                 TF_VAR_deis_root \
                 TF_VAR_prefix \
                 TF_VAR_ssh_keys
