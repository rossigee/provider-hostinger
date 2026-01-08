# Copyright 2025 Ross Golder. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Image configuration
PLATFORMS ?= linux_amd64
IMAGE ?= localhost/provider-hostinger:latest

# Docker build arguments
BUILDARGS := --build-arg TARGETOS=linux
BUILDARGS += --build-arg TARGETARCH=$(subst linux_,,$(PLATFORMS))

# Build targets
img.build:
	@docker build $(BUILDARGS) -t $(IMAGE) -f Dockerfile ../../../

img.publish:
	@docker push $(IMAGE)

img.promote:
	@docker tag $(FROM_IMAGE) $(TO_IMAGE)
	@docker push $(TO_IMAGE)

.PHONY: img.build img.publish img.promote
