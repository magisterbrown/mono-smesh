FROM python:3.10
WORKDIR /code
COPY ./requirements.txt .
RUN pip install --no-cache-dir --upgrade -r requirements.txt
COPY . /code
RUN mv ./configs/prod.py ./configs/conf.py

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]

