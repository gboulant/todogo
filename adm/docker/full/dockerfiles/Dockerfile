ARG SYSTEM=ubuntu
FROM $SYSTEM

# -------------------------------------------------------------------
# This part is to let you specify the apt proxy configuration. If you
# work with no proxy, just let the default value for HTTP_PROXY
# (void). If you work behind a proxy, then specify the http url of you
# proxy at build process using the docker option --build-arg
ARG HTTP_PROXY=""
ARG APTCONFIG=/etc/apt/apt.conf.d/99proxy
RUN echo "# FILE GENERATED BY gdocker" > $APTCONFIG
RUN echo "Acquire::http::Proxy \"$HTTP_PROXY\";" >> $APTCONFIG
RUN echo "Acquire::https::Proxy \"$HTTP_PROXY\";" >> $APTCONFIG
# -------------------------------------------------------------------

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y sudo && \
    apt-get install -y vim && \
    apt-get install -y git && \
    apt-get install -y make

RUN apt-get install -y golang

ARG USER=admin
ARG PASSWD=adm123

RUN useradd -ms /bin/bash $USER
RUN echo "$USER:$PASSWD" | chpasswd
RUN adduser $USER sudo

USER $USER
WORKDIR /home/$USER

RUN git config --global credential.helper cache && \
    git config --global credential.https://github.com.username $USER && \
    git clone https://gitlab.galuma.net/guiboule/todogo.git && \
    cd todogo && sudo make install
    
VOLUME  ["/home"]
