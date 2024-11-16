PACKAGE_NAME := obeops
DESCRIPTION := Ferramentas para manter servidor Minecraft estável e ativo.
MAINTAINER := Dev Isac Gondim
ARCHITECTURE := all

SHELL := /bin/bash

## Gerenciamento de versões

MAKE               := make --no-print-directory

DESCRIBE           := $(shell git describe --match "v*" --always --tags)
DESCRIBE_PARTS     := $(subst -, ,$(DESCRIBE))

VERSION_TAG        := $(word 1,$(DESCRIBE_PARTS))
COMMITS_SINCE_TAG  := $(word 2,$(DESCRIBE_PARTS))

VERSION            := $(subst v,,$(VERSION_TAG))
VERSION_PARTS      := $(subst ., ,$(VERSION))

MAJOR              := $(word 1,$(VERSION_PARTS))
MINOR              := $(word 2,$(VERSION_PARTS))
MICRO              := $(word 3,$(VERSION_PARTS))

NEXT_MAJOR         := $(shell echo $$(($(MAJOR)+1)))
NEXT_MINOR         := $(shell echo $$(($(MINOR)+1)))
NEXT_MICRO          = $(shell echo $$(($(MICRO)+$(COMMITS_SINCE_TAG))))

ifeq ($(strip $(COMMITS_SINCE_TAG)),)
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(MICRO)
CURRENT_VERSION_MINOR := $(CURRENT_VERSION_MICRO)
CURRENT_VERSION_MAJOR := $(CURRENT_VERSION_MICRO)
else
CURRENT_VERSION_MICRO := $(MAJOR).$(MINOR).$(NEXT_MICRO)
CURRENT_VERSION_MINOR := $(MAJOR).$(NEXT_MINOR).0
CURRENT_VERSION_MAJOR := $(NEXT_MAJOR).0.0
endif

DATE                = $(shell date +'%d.%m.%Y')
TIME                = $(shell date +'%H:%M:%S')
COMMIT             := $(shell git rev-parse HEAD)
AUTHOR             := $(firstword $(subst @, ,$(shell git show --format="%aE" $(COMMIT))))
BRANCH_NAME        := $(shell git rev-parse --abbrev-ref HEAD)

TAG_MESSAGE         = "$(TIME) $(DATE) $(AUTHOR) $(BRANCH_NAME)"
COMMIT_MESSAGE     := $(shell git log --format=%B -n 1 $(COMMIT))

CURRENT_TAG_MICRO  := "v$(CURRENT_VERSION_MICRO)"
CURRENT_TAG_MINOR  := "v$(CURRENT_VERSION_MINOR)"
CURRENT_TAG_MAJOR  := "v$(CURRENT_VERSION_MAJOR)"

# --- Version commands ---

.PHONY: version
version:
	@$(MAKE) version-micro

.PHONY: version-micro
version-micro:
	@echo "$(CURRENT_VERSION_MICRO)"

.PHONY: version-minor
version-minor:
	@echo "$(CURRENT_VERSION_MINOR)"

.PHONY: version-major
version-major:
	@echo "$(CURRENT_VERSION_MAJOR)"

# --- Tag commands ---

.PHONY: tag-micro
tag-micro:
	@echo "$(CURRENT_TAG_MICRO)"

.PHONY: tag-minor
tag-minor:
	@echo "$(CURRENT_TAG_MINOR)"

.PHONY: tag-major
tag-major:
	@echo "$(CURRENT_TAG_MAJOR)"

# -- Meta info ---

.PHONY: tag-message
tag-message:
	@echo "$(TAG_MESSAGE)"

.PHONY: commit-message
commit-message:
	@echo "$(COMMIT_MESSAGE)"

VERSION_APP := $(CURRENT_VERSION_MICRO)

# --- Compile ---
## THIS SCRIPT IS RESPONSIBLE FOR COMPILING THE APPLICATION BINARY ##
.PHONY: build-binary
build-binary:
	rm -rf /tmp/*
	mkdir -p /tmp/build_$(PACKAGE_NAME)/
	mkdir -p /tmp/Build/APPS/ /tmp/Build/bin/$(PACKAGE_NAME)/
	export GOOS=linux GOARCH=amd64
	cp App/shell/* /tmp/Build/bin/$(PACKAGE_NAME)/
	go build -o /tmp/Build/bin/$(PACKAGE_NAME)/$(PACKAGE_NAME) App/bin/bedrock/*.go

## THIS SCRIPT IS RESPONSIBLE FOR COMPILING THE APPLICATION
## IT SEPARATES ALL THE FILES INTO THEIR NECESSARY STRUCTURES AND COMPILES THE PACKAGE
.PHONY: build-package
build-package: build-binary
	echo "Iniciando o processo de compilação"
	mkdir -p /tmp/Build/$(PACKAGE_NAME)/usr/bin/ \
			 /tmp/Build/$(PACKAGE_NAME)/usr/lib/systemd/system \
			 /tmp/Build/$(PACKAGE_NAME)/DEBIAN/ \
			 /tmp/Build/APPS/ \
			 dist/ \
			 /tmp/Build/$(PACKAGE_NAME)/etc/$(PACKAGE_NAME) \
			 /tmp/Build/$(PACKAGE_NAME)/var/log/$(PACKAGE_NAME)/ \
			 /tmp/Build/$(PACKAGE_NAME)/var/ \
			 /tmp/build_$(PACKAGE_NAME)/
	mv /tmp/Build/bin/$(PACKAGE_NAME)/* /tmp/Build/$(PACKAGE_NAME)/usr/bin/
	cp App/dystro/debian/debian.control /tmp/Build/$(PACKAGE_NAME)/DEBIAN/control
	cp App/dystro/debian/install.sh /tmp/Build/$(PACKAGE_NAME)/DEBIAN/postinst
	cp App/systemd/* /tmp/Build/$(PACKAGE_NAME)/usr/lib/systemd/system/
	cp App/logs/* /tmp/Build/$(PACKAGE_NAME)/var/log/$(PACKAGE_NAME)/
	cp App/confs/* /tmp/Build/$(PACKAGE_NAME)/etc/$(PACKAGE_NAME)/
	sed -i "s/AppTemplate/$(PACKAGE_NAME)/; s/x.y.z/$(VERSION_APP)/; s/Dev/$(MAINTAINER)/; s/arc/$(ARCHITECTURE)/; s/DescriptionApp/$(DESCRIPTION)/" /tmp/Build/$(PACKAGE_NAME)/DEBIAN/control
	chmod +x /tmp/Build/$(PACKAGE_NAME)/usr/bin/* /tmp/Build/$(PACKAGE_NAME)/DEBIAN/postinst
	chmod 711 /tmp/Build/$(PACKAGE_NAME)/var/log/* /tmp/Build/$(PACKAGE_NAME)/etc/$(PACKAGE_NAME)/*
	dpkg-deb --build /tmp/Build/$(PACKAGE_NAME)/ /tmp/Build/APPS/obeops_app.deb
	echo "Pacote DEBIAN criado com sucesso!"

.PHONY: build-app
build-app:
	sudo docker build --no-cache -t $(PACKAGE_NAME) .
	sudo docker run -d -e TERM=xterm-256color --name ObsidianOps -p 2308:2308 $(PACKAGE_NAME)
