#syntax=docker/dockerfile:1.4
FROM node:16-buster AS src_builder
WORKDIR /app

# Copy the package.json and package-lock.json files into the container
COPY apps/olympus/hestia/assets/package*.json ./

# Install any needed packages
RUN npm ci

# Copy the rest of the application code into the container
COPY apps/olympus/hestia/assets/ ./

# Build the React app for production
RUN npm run build

FROM gcr.io/distroless/nodejs:16

# Set the working directory to /app
WORKDIR /app

# Copy the build output and node_modules from the previous build stage
COPY --from=src_builder /app/build /app/build
COPY --from=src_builder /app/node_modules /app/node_modules

ENV NODE_ENV=production

EXPOSE 3000

CMD ["node", "/app/node_modules/serve/bin/serve.js", "-s", "build", "-l", "3000"]
