FROM node:18-alpine AS builder
WORKDIR /app

# Accept build arguments
ARG VITE_WS_URL
ARG VITE_WS_PORT

# Set environment variables
ENV VITE_WS_URL=$VITE_WS_URL
ENV VITE_WS_PORT=$VITE_WS_PORT

COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
