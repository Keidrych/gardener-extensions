// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:generate packr2

package imagevector

import (
	"strings"

	"github.com/gobuffalo/packr/v2"

	"github.com/gardener/gardener/pkg/utils/imagevector"

	"k8s.io/apimachinery/pkg/util/runtime"
)

const (
	// TerraformerImageName is the name of the Terraformer image.
	TerraformerImageName = "terraformer"
)

var imageVector imagevector.ImageVector

func init() {
	box := packr.New("charts", "../../../charts")

	imagesYaml, err := box.FindString("images.yaml")
	runtime.Must(err)

	imageVector, err = imagevector.Read(strings.NewReader(imagesYaml))
	runtime.Must(err)

	imageVector, err = imagevector.WithEnvOverride(imageVector)
	runtime.Must(err)
}

// ImageVector is the image vector that contains al the needed images
func ImageVector() imagevector.ImageVector {
	return imageVector
}

// TerraformerImage retrieves the Terraformer image.
func TerraformerImage() string {
	image, err := imageVector.FindImage(TerraformerImageName, "", "")
	runtime.Must(err)

	return image.String()
}