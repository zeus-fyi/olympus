# Start with Ubuntu as the base image
FROM ubuntu:22.04

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Update and install Python and pip
RUN apt-get update && \
    apt-get install -y python3.9 python3-pip ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy your application code
COPY apps/olympus/pandora/ ./

# Install any needed packages specified in requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt

# Expose port 8000
EXPOSE 8000

# Define environment variable
ENV FLASK_APP=pandora.py

# Run app.py when the container launches
CMD ["flask", "run", "--host=0.0.0.0", "--port=8000"]