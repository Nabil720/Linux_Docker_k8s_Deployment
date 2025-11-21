# Build stage
FROM node:16 as build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Production stage - Node.js 
FROM node:16-alpine

WORKDIR /app

# শুধু production dependencies install 
COPY package*.json ./
RUN npm install --only=production

# Build files copy 
COPY --from=build /app/build ./build

# Simple Node.js server file copy 
COPY server.js ./server.js

EXPOSE 3000

CMD ["node", "server.js"]