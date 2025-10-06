FROM node:24-alpine

WORKDIR /app

COPY package*.json .
RUN npm ci

COPY bs-config.js esbuild-config-dev.js .

CMD ["npm", "run", "dev"]
