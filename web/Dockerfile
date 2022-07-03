
FROM node:16-alpine

COPY . /app/
WORKDIR /app

# RUN apk add --no-cache git python3 make g++
RUN npm install -g pnpm@next-7
RUN pnpm install