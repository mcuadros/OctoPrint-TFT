ARG IMAGE
FROM ${IMAGE}

ARG GO_TAGS
ENV GO_TAGS=${GO_TAGS}

RUN printf "deb http://archive.debian.org/debian/ jessie main\ndeb-src http://archive.debian.org/debian/ jessie main\ndeb http://security.debian.org jessie/updates main\ndeb-src http://security.debian.org jessie/updates main" > /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
  git build-essential \
  debhelper devscripts fakeroot git-buildpackage dh-make dh-systemd dh-golang \
  libcairo2-dev \
  libgtk-3-dev

# We cache go get gtk, to speed up builds.
#RUN go get -tags ${GO_TAGS} -v github.com/gotk3/gotk3/gtk/...

ADD . ${GOPATH}/src/github.com/mcuadros/OctoPrint-TFT/
#RUN go get -tags ${GO_TAGS} -v ./...

WORKDIR ${GOPATH}/src/github.com/mcuadros/OctoPrint-TFT/
