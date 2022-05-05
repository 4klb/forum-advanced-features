clean:
	@docker system prune

dbuild:
	@docker image build -f Dockerfile -t forumimage .
	@docker container run -p 8080:8080 --detach --name forum forumimage:latest