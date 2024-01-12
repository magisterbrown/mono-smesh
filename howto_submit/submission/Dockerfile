FROM python:latest
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY agent.py .
ENTRYPOINT ["python3", "agent.py"]
