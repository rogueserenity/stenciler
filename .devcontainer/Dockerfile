FROM mcr.microsoft.com/devcontainers/go:1-1.22-bookworm

# Install httpie
RUN curl -SsL https://packages.httpie.io/deb/KEY.gpg |\
  gpg --dearmor -o /usr/share/keyrings/httpie.gpg &&\
  echo "deb [arch=amd64 signed-by=/usr/share/keyrings/httpie.gpg] https://packages.httpie.io/deb ./" >\
  /etc/apt/sources.list.d/httpie.list && apt-get update && apt-get install -y \
  httpie python3 python3-dev pipx

USER vscode:vscode

# install ginkgo
RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.17.3

# install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh |\
  sh -s -- -b $(go env GOPATH)/bin v1.58.1

# install mockery
RUN go install github.com/vektra/mockery/v2@v2.43.1

# install cobra-cli
RUN go install github.com/spf13/cobra-cli@v1.3.0

# configure user environment
WORKDIR /home/vscode
ADD userhome/ .
