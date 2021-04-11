/*
Copyright (c) 2019 StackRox Inc.

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

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	microapplication "github.com/sbose78/micro-application/api/v1alpha1"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	microApplicationResource = metav1.GroupVersionResource{Version: "v1alpha1", Group: "argoproj.io", Resource: "microapplications"}
)

func applyUserInformation(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {

	if req.Resource != microApplicationResource {
		log.Printf("expect resource to be %s, found %s", microApplicationResource, req.Resource)
		return nil, nil
	}

	raw := req.Object.Raw
	microApplication := microapplication.MicroApplication{}

	if _, _, err := universalDeserializer.Decode(raw, nil, &microApplication); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	// Create patch operations to apply sensible defaults, if those options are not set explicitly.
	var patches []patchOperation

	key := "generated-creator"
	value := req.UserInfo.Username

	if len(microApplication.Annotations) == 0 {
		//|| microApplication.Annotations["microapplications/generated-creator"] {
		patches = append(patches, patchOperation{
			Op:   "add",
			Path: "/metadata/annotations",
			Value: map[string]string{
				key: value,
			},
		})
	} else {
		updatedAnnotations := microApplication.Annotations
		updatedAnnotations[key] = value
		patches = append(patches, patchOperation{
			Op:    "replace",
			Path:  "/metadata/annotations",
			Value: updatedAnnotations,
		})
	}
	return patches, nil
}

func main() {
	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applyUserInformation))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
