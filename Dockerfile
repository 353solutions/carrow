FROM rvernica/apache-arrow

RUN apt-get update && apt-get install -y \
    wget \
    vim \ 
    && rm -rf /var/lib/apt/lists/*

RUN cd /tmp && wget https://dl.google.com/go/go1.12.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.12.1.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin