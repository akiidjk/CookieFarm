FROM python:3.13-slim

WORKDIR /app
COPY service.py .

RUN pip install flask
CMD ["python3", "./service.py"]
