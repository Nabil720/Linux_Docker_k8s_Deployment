# Stage 1: Dependencies
FROM node:16 AS deps

WORKDIR /app
COPY package*.json ./
RUN npm ci

# Stage 2: Build
FROM node:16 AS builder

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Stage 3: Production
FROM node:16 AS production

WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=deps /app/package*.json ./
RUN npm install --only=production
RUN npm install -g serve

EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "3000"]
