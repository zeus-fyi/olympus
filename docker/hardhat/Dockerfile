#syntax=docker/dockerfile:1.4
FROM node:16-buster AS src_builder
WORKDIR /app

# Copy the package.json and package-lock.json files into the container
COPY apps/olympus/hardhat/package*.json ./

# Install any needed packages
RUN npm ci

# Copy the rest of the application code into the container
COPY apps/olympus/hardhat/ ./

# Build the React app for production

ENV NODE_ENV=production

EXPOSE 8545
