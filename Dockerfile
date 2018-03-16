FROM node:9.8-alpine

COPY ./src /app

RUN cd /app && \
    npm install

ENTRYPOINT [ "" ]
CMD ["sh", "-c", "cd /app && node app.js"]