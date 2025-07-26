FROM python:latest
COPY . /app
WORKDIR /app
CMD ["python", "*.py"]
