FROM python:slim

# runtime dependencies
RUN set -eux; \
	apt-get update; \
	apt-get upgrade -y; \
	rm -rf /var/lib/apt/lists/*

WORKDIR /app

ADD requirements.txt /app
RUN pip install -r requirements.txt

ADD . /app
ENTRYPOINT [ "streamlit", "run", "main.py", "--server.port", "8080", "--server.address", "0.0.0.0", "--browser.serverAddress", "0.0.0.0"]
