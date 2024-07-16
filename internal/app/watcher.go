package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/flohansen/kubernetes-configmap-watch-example/internal/product"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ProductRepository interface {
	Migrate(ctx context.Context) error
	Upsert(ctx context.Context, products []product.Model) error
}

type Watcher struct {
	repo ProductRepository
}

func NewWatcher(repo ProductRepository) *Watcher {
	return &Watcher{repo}
}

func (w *Watcher) Run(ctx context.Context) error {
	if err := w.repo.Migrate(ctx); err != nil {
		return errors.Wrap(err, "migration")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "in cluster config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "new kubernetes client")
	}

	watcher, err := clientset.
		CoreV1().
		ConfigMaps("default").
		Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=test-product-data",
		})
	if err != nil {
		return errors.Wrap(err, "config map watch")
	}

	for event := range watcher.ResultChan() {
		cm, ok := event.Object.(*corev1.ConfigMap)
		if !ok {
			continue
		}

		b, ok := cm.Data["products.json"]
		if !ok {
			continue
		}

		var products []product.Model
		if err := json.Unmarshal([]byte(b), &products); err != nil {
			return errors.Wrap(err, "json unmarshal products")
		}

		if err := w.repo.Upsert(ctx, products); err != nil {
			return errors.Wrap(err, "upsert products")
		}

		log.Printf("upserted %d products", len(products))
	}

	return nil
}
