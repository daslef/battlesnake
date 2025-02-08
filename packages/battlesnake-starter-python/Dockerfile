FROM --platform=linux/amd64 python:3.12

COPY . /usr/app
WORKDIR /usr/app

RUN pip install --upgrade pip && pip install -r requirements.txt

CMD [ "python", "main.py" ]