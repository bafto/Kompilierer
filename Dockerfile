FROM ubuntu:latest

# install dependencies
RUN apt-get update
RUN apt-get install -y \
	git \
	build-essential \
	llvm-12-dev \
	openssh-server

RUN locale-gen de_DE de_DE.UTF-8

# setup ssh
RUN mkdir /var/run/sshd
RUN echo 'root:rootpassword' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# default ssh port
EXPOSE 22

# install go
FROM golang:1.21

# setup the project
RUN git clone git@github.com:DDP-Projekt/Kompilierer.git
WORKDIR /Kompilierer
RUN go mod tidy
ENV DDPPATH=/Kompilierer/build/DDP

CMD ["/usr/sbin/sshd", "-D"]