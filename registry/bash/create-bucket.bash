set -eo pipefail

main() {
  # ensure registry bucket exists
  /app/bin/create_bucket ${BUCKET_NAME}
}
