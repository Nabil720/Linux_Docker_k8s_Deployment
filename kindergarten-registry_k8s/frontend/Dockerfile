# Use the official Node.js image as a base
FROM node:16

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json to the container
COPY package*.json ./

# Install the app dependencies (including react-scripts)
RUN npm install

# Verify installation
RUN npm list react-scripts

# Copy the rest of the app's files into the container
COPY . .

# Expose port 3000 for the frontend
EXPOSE 3000

# Start the app
CMD ["npm", "start"]
