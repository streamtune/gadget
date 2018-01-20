package utils

const ConstRepositoryDirectory string = "~/.gadget"
const ConstRepositoryFile string = "gadget.db"
const ConstDebugMode bool = false
const ConstRestServerMode string = "debug"
const ConstRestPort int = 9080

/*
	"dockerEndpoint" : "tcp://192.168.99.100:2376",
	"localDockerEndpoint": "unix:///var/run/docker.sock"

	export DOCKER_TLS_VERIFY="1"
	export DOCKER_HOST="tcp://192.168.99.100:2376"
	export DOCKER_CERT_PATH="/home/gianluca/.docker/machine/machines/dev"
	export DOCKER_MACHINE_NAME="dev"
*/

const ConstUseDockerMachine bool = false
const ConstMachineDockerEndpoint string = "tcp://192.168.99.100:2376"
const ConstMachineDockerCertFile string = "/home/gianluca/.docker/machine/machines/dev/cert.pem"
const ConstMachineDockerKeyFile string = "/home/gianluca/.docker/machine/machines/dev/key.pem"
const ConstMachineDockerCAFile string = "/home/gianluca/.docker/machine/machines/dev/ca.pem"

const ConstLocalDockerEndpoint string = "unix:///var/run/docker.sock"
