FROM openjdk:latest
COPY . /app
WORKDIR /app
RUN javac *.java
CMD ["java", "Main"]
