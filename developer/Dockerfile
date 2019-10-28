FROM ubuntu:18.04

# Tools
RUN apt-get update && apt-get install -y \
    vim \
    wget \
    git  \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Go installation
RUN cd /tmp && \
    wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.12.9.linux-amd64.tar.gz && \
    rm go1.12.9.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Apache Arrow
RUN cd /tmp && \
    wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && \
    bash Miniconda3-latest-Linux-x86_64.sh -b -p /miniconda && \
    rm Miniconda3-latest-Linux-x86_64.sh
ENV PATH="/miniconda/bin:${PATH}"
RUN conda install -y \
    conda-forge::arrow-cpp=0.14
