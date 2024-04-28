# Use an official Node.js runtime as a base image
FROM node:20

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Run Prisma migrations to generate the database schema
RUN npx prisma generate

# Expose the port your app runs on
EXPOSE 9094

# Command to start both Prisma and Express.js application
CMD ["npm", "start"]
