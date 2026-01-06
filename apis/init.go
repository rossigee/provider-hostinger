/*
Copyright 2025 Ross Golder.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apis

import (
	"k8s.io/apimachinery/pkg/runtime"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
	instancev1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	backupv1beta1 "github.com/rossigee/provider-hostinger/apis/backup/v1beta1"
	firewallv1beta1 "github.com/rossigee/provider-hostinger/apis/firewall/v1beta1"
	sshkeyv1beta1 "github.com/rossigee/provider-hostinger/apis/sshkey/v1beta1"
)

// Scheme is the runtime.Scheme for all Hostinger APIs
var Scheme = runtime.NewScheme()

// SchemeBuilder builds the Hostinger scheme
var SchemeBuilder = runtime.NewSchemeBuilder(
	v1beta1.SchemeBuilder.AddToScheme,
	instancev1beta1.SchemeBuilder.AddToScheme,
	backupv1beta1.SchemeBuilder.AddToScheme,
	firewallv1beta1.SchemeBuilder.AddToScheme,
	sshkeyv1beta1.SchemeBuilder.AddToScheme,
)

// AddToScheme adds all Hostinger API types to the scheme
func AddToScheme(s *runtime.Scheme) error {
	return SchemeBuilder.AddToScheme(s)
}

func init() {
	if err := AddToScheme(Scheme); err != nil {
		panic(err)
	}
}
