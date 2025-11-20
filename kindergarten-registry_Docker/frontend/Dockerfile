# Stage 1: Dependencies stage
FROM node:16 AS deps

WORKDIR /app
COPY package*.json ./
RUN npm ci

# Stage 2: Build stage 
FROM node:16 AS builder

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Stage 3: Production stage
FROM node:16 AS production

WORKDIR /app
COPY --from=deps /app/package*.json ./
RUN npm ci --only=production

COPY --from=builder /app/build ./build
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["npm", "start"]