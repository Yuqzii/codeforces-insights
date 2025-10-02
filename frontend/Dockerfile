FROM node:24-alpine

WORKDIR /app

COPY package*.json .
RUN npm ci

CMD ["npm", "run", "dev"]
