# Build stage
FROM node:16 as build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage - Node.js দিয়ে serve করুন
FROM node:16-alpine

WORKDIR /app

# শুধু production dependencies install করুন
COPY package*.json ./
RUN npm install --only=production

# Build files copy করুন
COPY --from=build /app/build ./build

# Simple Node.js server file copy করুন
COPY server.js ./server.js

EXPOSE 3000

CMD ["node", "server.js"]