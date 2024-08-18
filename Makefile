build:
	docker build -t my_http_server .
run: 
	docker run --name http_server -p 8000:8000 my_http_server
up:
	pytest tests.py
clean:
	docker rm -f http_server

