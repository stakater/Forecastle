FROM node:9.8-alpine

COPY ./src /app

ENTRYPOINT [ "" ]
CMD ["sh", "-c", "cd /app && node app.js"]
