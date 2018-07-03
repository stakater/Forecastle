FROM node:9.8-alpine

COPY ./src /app

# Update apk repository list in separate layer
# so that install layer does not run everytime
RUN apk update

# Install curl unzip and some handy tools
RUN echo "===> Installing Utilities from apk ..."  && \
    apk -v --update --progress add curl unzip jq


ARG VERSION_URL=https://raw.githubusercontent.com/stakater/stk/add-jenkinsfile-to-stk/.version
COPY ./hub /hub

RUN TOKEN=$(cat /hub) \
    && VERSION=$(curl -H 'Authorization: token '"${TOKEN}" -H 'Accept: application/vnd.github.v4.raw' -L ${VERSION_URL}) \
    && TAG_URL=https://api.github.com/repos/stakater/stk/releases/tags/${VERSION} \
    && FILE_NAME=stk_${VERSION}_linux_386.tar.gz \
    && ASSET_URL=$(curl -H 'Authorization: token '"${TOKEN}" ${TAG_URL} | jq '.assets[]  | select(.name == "'${FILE_NAME}'") | .url') \
    && echo ${VERSION} \
    && echo ${TAG_URL} \
    && echo ${FILE_NAME} \
    && ASSET_URL="${ASSET_URL%\"}" \
    && ASSET_URL="${ASSET_URL#\"}" \
    && echo "$ASSET_URL" \
    && curl -L -H "Accept: application/octet-stream" ${ASSET_URL}?access_token=${TOKEN} > stk.tar.gz \
    && mkdir ./temp && tar zxvf ./stk.tar.gz -C ./temp \
    && cp ./temp/stk /usr/local/bin/stk \
    && rm -rf ./temp/* \
    && rm /hub && rm ./stk.tar.gz


RUN cd /app && npm install

ENTRYPOINT [ "" ]
CMD ["sh", "-c", "cd /app && node app.js"]
