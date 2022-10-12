build:
	-cd server1 && go build -o server1
	-cd server2 && go build -o server2
clean:
	-cd server1 && rm server1
	-cd server2 && rm server2
