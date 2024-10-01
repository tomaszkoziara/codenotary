FROM node:22-alpine as builder

WORKDIR /app
COPY . /app

ENV REACT_APP_API_PORT=8000

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:22-alpine

EXPOSE 3000

COPY --from=builder /app/build ./build
CMD ["npx", "http-server", "--auto-confirm", "-p", "3000", "./build"]