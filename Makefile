all: gateway flights tickets privileges identity-provider

# Creating images
IMAGES?=gateway flights tickets privileges identity-provider
$(IMAGES):
	$(MAKE) docker-push -C ./src/$@
