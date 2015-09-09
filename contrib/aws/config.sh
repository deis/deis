export SUGGEST_DEV_REGISTRY="registry.hub.docker.com"
export STACK_TAG="${STACK_TAG:-test}-${DEIS_ID}"
export DEIS_NUM_INSTANCES=${DEIS_NUM_INSTANCES:-3}
export STACK_NAME="${STACK_NAME:-deis-${STACK_TAG}}"

prompt "Enter your AWS key:" AWS_ACCESS_KEY_ID
prompt "Enter your AWS secret key:" AWS_SECRET_ACCESS_KEY

rigger-save-vars AWS_ACCESS_KEY_ID \
                 AWS_SECRET_ACCESS_KEY \
                 STACK_NAME \
                 STACK_TAG
