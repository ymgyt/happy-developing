format="%-10s OK\n"


function install_kick() {
    (which kick 1>/dev/null && printf "${format}" "kick") || go get -u github.com/isomorphicgo/kick
}

function install_circlecicli() {
    which circleci 1>/dev/null  && printf "${format}" "circleci" && return 0;
    sudo curl -fLSs https://circle.ci/cli | bash
}

function install_vbox() {
    which virtualbox 1>/dev/null && printf "${format}" "virtualbox" && return 0;
    brew cask install virtualbox
}

function install_minikube() {
    which minikube 1>/dev/null && printf "${format}" "minikube" && return 0;
    brew cask install minikube
}

function install_kubectl() {
    which kubectl 1>/dev/null && printf "${format}" "kubectl" && return 0;
    brew install kubernetes-cli
}

function install_task() {
    which task 1>/dev/null && printf "${format}" "task" && return 0;
    go get -u -v github.com/go-task/task/cmd/task
}

install_kick
install_circlecicli
install_minikube
install_kubectl
install_task