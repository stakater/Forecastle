FROM node:9.8-alpine

COPY ./src /app

COPY stk /usr/local/bin/stk

ENTRYPOINT [ "" ]
CMD ["sh", "-c", "cd /app && node app.js"]
