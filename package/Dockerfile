FROM registry.suse.com/bci/bci-base:15.6

# nfs-client is needed by the dep https://github.com/longhorn/backupstore to check backup store availability.
RUN zypper -n rm container-suseconnect && \
    zypper -n install curl gzip tar nfs-client && \
    zypper -n clean -a && rm -rf /tmp/* /var/tmp/* /usr/share/doc/packages/* && \
    useradd -M cloudweav && \
    mkdir -p /var/lib/cloudweav/cloudweav && \
    chown -R cloudweav /var/lib/cloudweav/cloudweav /usr/local/bin

WORKDIR /var/lib/cloudweav/cloudweav

ENV CLOUDWEAV_UI_VERSION latest
ENV CLOUDWEAV_UI_PATH /usr/share/cloudweav/cloudweav
# Please update the api-ui-version in pkg/settings/settings.go when updating the version here.
ENV CLOUDWEAV_API_UI_VERSION 1.1.9

ENV CLOUDWEAV_UI_PLUGIN_BUNDLED_VERSION latest

ARG ARCH=amd64
ARG VERSION=dev
ENV CLOUDWEAV_SERVER_VERSION ${VERSION}
ENV TINI_VERSION v0.19.0
ENV TINI_URL_amd64=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini \
    TINI_URL_arm64=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-arm64 \
    TINI_URL_s390x=https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-s390x \
    TINI_URL=TINI_URL_${ARCH}

RUN curl -sLf ${!TINI_URL} > /usr/bin/tini && chmod +x /usr/bin/tini

RUN mkdir -p /usr/share/cloudweav/cloudweav && \
    mkdir -p /usr/share/cloudweav/cloudweav/api-ui && \
    mkdir -p /usr/share/cloudweav/cloudweav/v1/cloudweav/plugin-assets
COPY assests/CLOUDWEAV_UI_VERSION_latest.tar.gz /usr/share/cloudweav/cloudweav/
COPY assests/CLOUDWEAV_API_UI_VERSION_1.1.9.tar.gz /usr/share/cloudweav/cloudweav/api-ui/
COPY assests/CLOUDWEAV_UI_PLUGIN_BUNDLED_VERSION_latest.tar.gz /usr/share/cloudweav/cloudweav/v1/cloudweav/plugin-assets/
RUN mkdir -p /usr/share/cloudweav/cloudweav && \
    cd /usr/share/cloudweav/cloudweav && \
    #curl -sL https://releases.rancher.com/cloudweav-ui/dashboard/${CLOUDWEAV_UI_VERSION}.tar.gz | tar xvzf - --strip-components=2 && \
    tar xvzf CLOUDWEAV_UI_VERSION_latest.tar.gz --strip-components=2 && \
    mkdir -p /usr/share/cloudweav/cloudweav/api-ui && \
    cd /usr/share/cloudweav/cloudweav/api-ui && \
    #curl -sL https://releases.rancher.com/api-ui/${CLOUDWEAV_API_UI_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    tar xvzf CLOUDWEAV_API_UI_VERSION_1.1.9.tar.gz --strip-components=1 && \
    mkdir -p /usr/share/cloudweav/cloudweav/v1/cloudweav/plugin-assets && \
    cd /usr/share/cloudweav/cloudweav/v1/cloudweav/plugin-assets && \
    #curl -sL https://releases.rancher.com/cloudweav-ui/plugin/cloudweav-${CLOUDWEAV_UI_PLUGIN_BUNDLED_VERSION}.tar.gz | tar xvzf - --strip-components=1 && \
    tar xvzf CLOUDWEAV_UI_PLUGIN_BUNDLED_VERSION_latest.tar.gz --strip-components=1 && \
    cd /var/lib/cloudweav/cloudweav

COPY entrypoint.sh cloudweav /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh

VOLUME /var/lib/cloudweav/cloudweav
ENTRYPOINT ["entrypoint.sh"]
