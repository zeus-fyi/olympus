# Use an official Python runtime as a parent image
FROM python:3.9
# Set the working directory in the container
WORKDIR /app

COPY apps/olympus/pandora/ .

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Make port 5000 available to the world outside this container
EXPOSE 8000
# Define environment variable
ENV FLASK_APP=pandora.py
# Run app.py when the container launches
CMD ["flask", "run", "--host=0.0.0.0", "--port=8000"]
