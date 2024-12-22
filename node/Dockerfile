FROM node:lts-alpine

WORKDIR /APP

COPY ./package*.json .

RUN npm install

COPY  . .

WORKDIR /APP/src

EXPOSE 80

CMD [ "node","index.js" ]