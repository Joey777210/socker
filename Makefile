#Go paramemters
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get

BINARY = "Socker"
VERSION = 1.0

default:
	@echo "Socker v1.0"
build:
	@$(GOBUILD)
	@$(GOBUILD) ../SockerMQTTWatcher
	sudo cp ../ubuntu.tar /root
get:
	@$(GOGET)
	@$(GOGET) ../SockerMQTTWatcher
