########## UBUNTU BASE ##########
# Contains common utils that are useful everywhere, even in runtime containers
FROM ubuntu:23.10 AS ubuntu_base

RUN apt-get update && \
    apt-get satisfy -y bash-completion curl dnsutils file gzip iputils-ping less net-tools tar vim wget which && \
    apt-get clean

########## NODEJS BASE ##########
# To use in a plain Ubuntu image:
#   Copy /usr/local from here to any image that needs node
FROM node:20 AS node_base

########## GOLANG BASE ##########
# To use in a plain Ubuntu image:
#   Copy /usr/local
#   Copy /go
#   Set PATH=/usr/local/go/bin:/go/bin:$PATH
FROM golang:1.21 AS golang_base 

# Install global tools for building/developing Go stuff
COPY golang-tools.txt /go/tools.txt
RUN for PACKAGE in $(cat /go/tools.txt | sed 's/#.*$//'); do \
    if [ -z "$PACKAGE" ]; then continue; fi; \
    go install $PACKAGE || exit 1; \
    done

########## DEVCONTAINER ##############
FROM ubuntu_base AS devcontainer

# Install dependencies
COPY .devcontainer/packages.txt /opt/devcontainer-packages.txt
RUN apt-get update && \
    apt-get satisfy -y $(cat /opt/devcontainer-packages.txt) && \
    apt-get clean

# Include Node
COPY --from=node_base /usr/local /usr/local/

# Include Go
COPY --from=golang_base /usr/local /usr/local/
COPY --from=golang_base /go /go
ENV PATH=/usr/local/go/bin:/go/bin:$PATH

# Set up the developer user
RUN useradd -ms /bin/bash developer
RUN echo "developer ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers.d/developer
WORKDIR /home/developer
USER developer

# Set up developer custom configs
COPY .devcontainer/bashrc.d /home/developer/.bashrc.d

CMD [ "echo", "devcontainer is not runnable like this"]

########## DEFAULT NO-OP TARGET; ALL BUILDS MUST SPECIFY A TARGET ##############
FROM scratch AS default
CMD [ "ha-ha-you-built-without-a-target" ]
