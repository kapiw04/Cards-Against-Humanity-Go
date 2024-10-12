# Dockerfile for React
FROM node:22-alpine as build

WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./
COPY package-lock.json ./
RUN npm install --silent
COPY . ./
EXPOSE 3000
RUN npm run build
RUN npm install -g serve

CMD [ "serve", "-s", "build" ]