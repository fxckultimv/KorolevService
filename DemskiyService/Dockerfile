FROM python:3.9-slim

WORKDIR /app

COPY external_tickets_api.py .

RUN pip install flask flask-restful

CMD ["python", "external_tickets_api.py"]
