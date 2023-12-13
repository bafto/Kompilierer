FROM ubuntu:latest

WORKDIR /

# install dependencies
RUN apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y git \
	build-essential \
	llvm-12-dev \
	openssh-server \
	locales

RUN locale-gen de_DE de_DE.UTF-8

# setup ssh
RUN mkdir /var/run/sshd
RUN echo 'root:rootpassword' | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/#PermitUserEnvironment no/PermitUserEnvironment yes/' /etc/ssh/sshd_config

# default ssh port
EXPOSE 22

# install go 
RUN wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
RUN rm -f go1.21.5.linux-amd64.tar.gz
ENV PATH="${PATH}:/usr/local/go/bin"

# setup the project
WORKDIR /usr/src
RUN git clone https://github.com/DDP-Projekt/Kompilierer.git
WORKDIR /usr/src/Kompilierer
RUN git switch dev
RUN go mod tidy
ENV DDPPATH=/usr/src/Kompilierer/build/DDP
ENV PATH="${PATH}:/usr/src/Kompilierer/build/DDP/bin"

# for ssh
RUN echo "export PATH=$PATH" >> /root/.bashrc
RUN echo "export DDPPATH=$DDPPATH" >> /root/.bashrc
RUN echo "cd /usr/src/Kompilierer" >> /root/.bashrc

CMD ["/usr/sbin/sshd", "-D"]