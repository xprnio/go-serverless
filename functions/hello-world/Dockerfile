FROM node:21-alpine

WORKDIR /app

ADD package.json /app
RUN npm install

ADD . /app
RUN npm run build

ENTRYPOINT ["npm", "start"]
