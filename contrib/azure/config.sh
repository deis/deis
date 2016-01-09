# use COREOS_VERSION in azure deploy json file
source "${DEIS_ROOT}/contrib/utils.sh"
export COREOS_VERSION

export DEIS_TEST_DOMAIN="xip.io"
export AZURE_INSTANCE_SIZE="${AZURE_INSTANCE_SIZE:-Standard_A3}"
export AZURE_DOMAIN_PREFIX="deis${DEIS_ID}"
export AZURE_RESOURCE_GROUP="deis${DEIS_ID}"
export AZURE_STORAGE_NAME="deis${DEIS_ID}"
export AZURE_REGION="West US"

rigger-save-vars AZURE_DOMAIN_PREFIX \
                 AZURE_INSTANCE_SIZE \
                 AZURE_REGION \
                 AZURE_RESOURCE_GROUP \
                 AZURE_STORAGE_NAME \
                 COREOS_VERSION \
                 DEIS_TEST_DOMAIN
