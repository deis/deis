set -eo pipefail

hello-bash() {
	echo "Hello world from Bash"
}

main() {
	echo "Arguments:" "$@"
	hello-bash | reverse
	curl -s https://api.github.com/repos/progrium/go-basher | json-pointer /owner/login
}
