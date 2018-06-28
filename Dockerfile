FROM node:9.8-alpine

COPY ./src /app

# Update apk repository list in separate layer
# so that install layer does not run everytime
RUN apk update

# Install ansible, boto, aws-cli, and some handy tools
RUN echo "===> Installing Utilities from apk ..."  && \
    apk -v --update --progress add sudo git bash wget openssh groff less python py-pip curl jq unzip nodejs=8.9.3-r1 coreutils python py-pip openssl ca-certificates make sshpass openssh-client


ARG VERSION_URL=https://raw.githubusercontent.com/stakater/stk/add-jenkinsfile-to-stk/.

COPY ./hub /hub
RUN TOKEN=$(cat /hub) \
    && VERSION=$(curl -H 'Authorization: token ${TOKEN}' -H 'Accept: application/vnd.github.v4.raw' -L ${VERSION_URL}) \
    && curl -H 'Authorization: token ${TOKEN}' -H 'Accept: application/vnd.github.v4.raw' \
    -L https://github.com/stakater/stk/releases/download/${VERSION}/stk_${VERSION}_linux_386.tar.gz | tar zxv -C ./temp \
    && cp ./temp/stk /usr/local/bin/stk \
    && rm -rf ./temp/* \
    && rm /hub

RUN cd /app && npm install

ENTRYPOINT [ "" ]
CMD ["sh", "-c", "cd /app && node app.js"]
