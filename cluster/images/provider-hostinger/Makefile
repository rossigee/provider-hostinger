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
	@$(MAKE) BUILDARGS="$(BUILDARGS)" img.build.shared

img.build.shared:
	@rm -rf $(IMAGE_TEMP_DIR)
	@mkdir -p $(IMAGE_TEMP_DIR)
	@cp Dockerfile $(IMAGE_TEMP_DIR)/
	@cp -r ../../../_output/bin/ $(IMAGE_TEMP_DIR)/
	@cp ../../../package/crossplane.yaml $(IMAGE_TEMP_DIR)/package.yaml
	@cp -r ../../../package/crds $(IMAGE_TEMP_DIR)/crds
	@docker buildx build $(BUILDARGS) --build-arg TARGETOS=linux --build-arg TARGETARCH=$(subst linux_,,$(PLATFORMS)) --load -t $(IMAGE) $(IMAGE_TEMP_DIR)

img.publish:
	@docker push $(IMAGE)

img.promote:
	@docker tag $(FROM_IMAGE) $(TO_IMAGE)
	@docker push $(TO_IMAGE)

.PHONY: img.build img.build.shared img.publish img.promote img.clean
img.clean:
	@rm -rf $(IMAGE_TEMP_DIR)
