# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# Copy package files
COPY package*.json ./


# Install dependencies
RUN npm ci

# Copy source code
COPY . .


# Build the application
RUN npm run build

# Production stage
FROM nginx:stable-alpine

# Copy built assets from build stage
COPY --from=build /app/.next/standalone/ /usr/share/nginx/html/
COPY --from=build /app/.next/static /usr/share/nginx/html/.next/static
COPY --from=build /app/public /usr/share/nginx/html/

# Copy nginx config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
