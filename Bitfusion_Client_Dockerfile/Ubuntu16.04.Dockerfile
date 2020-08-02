FROM ubuntu:16.04

RUN sed -i 's/archive.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    apt-get update && apt-get install -y --no-install-recommends \
        apt-utils \
        build-essential \
        ca-certificates \
        curl \
        kmod

# TODO
#RUN echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ xenial main" > /etc/apt/sources.list && \
#    echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ xenial-updates main" >> /etc/apt/sources.list && \
#    echo "deb [arch=amd64] http://archive.ubuntu.com/ubuntu/ xenial-security main" >> /etc/apt/sources.list && \
#    usermod -o -u 0 -g 0 _apt

#ARG BASE_URL=http://us.download.nvidia.com/XFree86/Linux-x86_64
ARG BASE_URL=https://cn.download.nvidia.com/tesla
ARG DRIVER_VERSION=440.33.01
ENV DRIVER_VERSION=$DRIVER_VERSION

# Install the userspace components and copy the kernel module sources.
RUN cd /tmp && \
    curl -fSslL -O $BASE_URL/$DRIVER_VERSION/NVIDIA-Linux-x86_64-$DRIVER_VERSION.run && \
    sh NVIDIA-Linux-x86_64-$DRIVER_VERSION.run -x && \
    cd NVIDIA-Linux-x86_64-$DRIVER_VERSION* && \
    ./nvidia-installer --silent \
                       --no-kernel-module \
                       --install-compat32-libs \
                       --no-nouveau-check \
                       --no-nvidia-modprobe \
                       --no-rpms \
                       --no-backup \
                       --no-check-for-alternate-installs \
                       --no-libglx-indirect \
                       --no-install-libglvnd \
                       --x-prefix=/tmp/null \
                       --x-module-path=/tmp/null \
                       --x-library-path=/tmp/null \
                       --x-sysconfig-path=/tmp/null && \
    rm -rf /tmp/*

#COPY nvidia-driver /usr/local/bin

#WORKDIR /usr/src/nvidia-$DRIVER_VERSION

#ARG PUBLIC_KEY=empty
#COPY ${PUBLIC_KEY} kernel/pubkey.x509

#ARG PRIVATE_KEY
#ARG KERNEL_VERSION=generic,generic-hwe-16.04

# Compile the kernel modules and generate precompiled packages for use by the nvidia-installer.
#RUN apt-get update && \
#    for version in $(echo $KERNEL_VERSION | tr ',' ' '); do \
#        nvidia-driver update -k $version -t builtin ${PRIVATE_KEY:+"-s ${PRIVATE_KEY}"}; \
#    done && \
#    rm -rf /var/lib/apt/lists/*

#ENTRYPOINT ["nvidia-driver", "init"]

# Start of CUDA base
RUN apt-get install -y --no-install-recommends \
ca-certificates apt-transport-https gnupg-curl && \
    NVIDIA_GPGKEY_SUM=d1be581509378368edeec8c1eb2958702feedf3bc3d17011adbf24efacce4ab5 && \
    NVIDIA_GPGKEY_FPR=ae09fe4bbd223a84b2ccfce3f60f4b3d7fa2af80 && \
    apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1604/x86_64/7fa2af80.pub && \
    apt-key adv --export --no-emit-version -a $NVIDIA_GPGKEY_FPR | tail -n +5 > cudasign.pub && \
    echo "$NVIDIA_GPGKEY_SUM  cudasign.pub" | sha256sum -c --strict - && rm cudasign.pub && \
    echo "deb https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1604/x86_64 /" > /etc/apt/sources.list.d/cuda.list && \
    echo "deb https://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1604/x86_64 /" > /etc/apt/sources.list.d/nvidia-ml.list && \
rm -rf /var/lib/apt/lists/*

ENV CUDA_VERSION 10.2.89

ENV CUDA_PKG_VERSION 10-2=$CUDA_VERSION-1

# For libraries in the cuda-compat-* package: https://docs.nvidia.com/cuda/eula/index.html#attachment-a
RUN apt-get update && apt-get install -y --no-install-recommends \
        cuda-cudart-$CUDA_PKG_VERSION \
cuda-compat-10-2 && \
ln -s cuda-10.2 /usr/local/cuda

# Required for nvidia-docker v1
#RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
#    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

ENV PATH /usr/local/nvidia/bin:/usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64

# nvidia-container-runtime
#ENV NVIDIA_VISIBLE_DEVICES all
#ENV NVIDIA_DRIVER_CAPABILITIES compute,utility
#ENV NVIDIA_REQUIRE_CUDA "cuda>=10.2 brand=tesla,driver>=384,driver<385 brand=tesla,driver>=396,driver<397 brand=tesla,driver>=410,driver<411 brand=tesla,driver>=418,driver<419"

# Start of CUDA runtime
#ARG IMAGE_NAME
#FROM ${IMAGE_NAME}:10.2-base-ubuntu16.04

ENV NCCL_VERSION 2.5.6

RUN apt-get install -y --no-install-recommends \
    cuda-libraries-$CUDA_PKG_VERSION \
cuda-nvtx-$CUDA_PKG_VERSION \
libcublas10=10.2.2.89-1 \
libnccl2=$NCCL_VERSION-1+cuda10.2 && \
    apt-mark hold libnccl2

# Start of CUDA cudnn
ENV CUDNN_VERSION 7.6.5.32
LABEL com.nvidia.cudnn.version="${CUDNN_VERSION}"

RUN apt-get install -y --no-install-recommends \
    libcudnn7=$CUDNN_VERSION-1+cuda10.2 \
&& \
    apt-mark hold libcudnn7

# Start of bitfusion
RUN apt-get install -y --no-install-recommends \
    wget libjson-c2 libjsoncpp1 librdmacm1 libssl-dev libibverbs1 libnuma1 libcapstone3 libnl-3-200 libnl-route-3-200 && \
    rm -rf /var/lib/apt/lists/*

RUN cd /tmp && \
    curl -fSslL -O http://pa-dbc1131.eng.vmware.com/liuqi/bitfusion/ubuntu_16.04.x86_64.tgz && \
    tar zxf ubuntu_16.04.x86_64.tgz && \
    ./install -m binaries
