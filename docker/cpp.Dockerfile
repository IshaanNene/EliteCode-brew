FROM gcc:latest
COPY . /app
WORKDIR /app
RUN g++ *.cpp -o app
CMD ["./app"]
