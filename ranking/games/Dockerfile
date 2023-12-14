FROM python:3.11
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY play.py .

ENTRYPOINT python3 play.py
