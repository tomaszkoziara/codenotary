FROM node:22-alpine as builder

WORKDIR /app
COPY . /app

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:22-alpine

ENV HTTP_PORT=3000
EXPOSE $HTTP_PORT

COPY --from=builder /app/build .
CMD ["npx", "http-server", "--auto-confirm", "-p", "3000", "./build"]