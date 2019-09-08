FROM ubuntu

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y git && \
    apt-get install -y make

RUN apt-get install -y golang

RUN git clone https://gitlab.galuma.net/guiboule/todogo.git && \
    cd todogo && make install
