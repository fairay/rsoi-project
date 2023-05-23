all: gateway flights tickets privileges identity-provider statistics

# Creating images
IMAGES?=gateway flights tickets privileges identity-provider statistics
$(IMAGES):
	$(MAKE) docker-push -C ./src/$@
