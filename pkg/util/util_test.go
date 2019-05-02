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

package util_test

import (
	"context"
	"fmt"
	"time"

	mockclient "github.com/gardener/gardener-extensions/pkg/mock/controller-runtime/client"
	. "github.com/gardener/gardener-extensions/pkg/util"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	var (
		ctrl *gomock.Controller
		c    *mockclient.MockClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		c = mockclient.NewMockClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#GetSecretByRef", func() {
		var (
			ctx = context.TODO()

			name      = "foo"
			namespace = "bar"
		)

		It("should get the secret", func() {
			var (
				objectMeta = metav1.ObjectMeta{
					Name:      name,
					Namespace: namespace,
				}
				data = map[string][]byte{
					"foo": []byte("bar"),
				}
			)

			c.EXPECT().Get(ctx, kutil.Key(namespace, name), gomock.AssignableToTypeOf(&corev1.Secret{})).DoAndReturn(func(_ context.Context, _ client.ObjectKey, secret *corev1.Secret) error {
				secret.ObjectMeta = objectMeta
				secret.Data = data
				return nil
			})

			secret, err := GetSecretByRef(ctx, c, corev1.SecretReference{
				Name:      name,
				Namespace: namespace,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(secret).To(Equal(&corev1.Secret{
				ObjectMeta: objectMeta,
				Data:       data,
			}))
		})

		It("should return the error", func() {
			ctx := context.TODO()

			c.EXPECT().Get(ctx, kutil.Key(namespace, name), gomock.AssignableToTypeOf(&corev1.Secret{})).Return(fmt.Errorf("error"))

			secret, err := GetSecretByRef(ctx, c, corev1.SecretReference{
				Name:      name,
				Namespace: namespace,
			})

			Expect(err).To(HaveOccurred())
			Expect(secret).To(Equal(&corev1.Secret{}))
		})
	})

	Describe("#WaitUntilResourceDeleted", func() {
		var (
			namespace = "bar"
			name      = "foo"
			key       = kutil.Key(namespace, name)
			configMap = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      name,
				},
			}
		)

		It("should wait until the resource is deleted", func() {
			ctx := context.TODO()

			gomock.InOrder(
				c.EXPECT().
					Get(ctx, key, configMap),
				c.EXPECT().
					Get(ctx, key, configMap),
				c.EXPECT().
					Get(ctx, key, configMap).
					Return(apierrors.NewNotFound(schema.GroupResource{}, name)),
			)

			Expect(WaitUntilResourceDeleted(ctx, c, configMap, time.Microsecond)).To(Succeed())
		})

		It("return an unexpected error", func() {
			ctx := context.TODO()

			expectedErr := fmt.Errorf("unexpected")
			c.EXPECT().
				Get(ctx, key, configMap).
				Return(expectedErr)

			err := WaitUntilResourceDeleted(ctx, c, configMap, time.Microsecond)

			Expect(err).To(HaveOccurred())
			Expect(err).To(BeIdenticalTo(expectedErr))
		})
	})
})
