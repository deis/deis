export SUGGEST_DEIS_TEST_DOMAIN="local3.deisapp.com"
export DEIS_TEST_SSH_KEY="${HOME}/.vagrant.d/insecure_private_key"
export SUGGEST_DEIS_SSH_KEY="${DEIS_TEST_SSH_KEY}"

rigger-save-vars SUGGEST_DEIS_TEST_DOMAIN SUGGEST_DEIS_SSH_KEY DEIS_TEST_SSH_KEY
