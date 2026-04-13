FROM node:22-alpine

WORKDIR /app

COPY package*.json .
RUN npm install

COPY . .

WORKDIR /app/src/services/gateway

EXPOSE 9001 5173

CMD ["npm", "run", "dev"]
