package cleaner

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type cleanupEntry struct {
	client client.Client
	object client.Object
}

type CleanerCollector interface {
	Empty() bool
	Add(client.Client, client.Object)
}

type Cleaner interface {
	CleanerCollector
	Run(ctx context.Context) []error
}

func New(timeout time.Duration, delay *time.Duration) Cleaner {
	return &cleaner{
		delay:   delay,
		timeout: timeout,
	}
}

type cleaner struct {
	delay   *time.Duration
	timeout time.Duration
	entries []cleanupEntry
}

func (c *cleaner) Add(client client.Client, object client.Object) {
	c.entries = append(c.entries, cleanupEntry{
		client: client,
		object: object,
	})
}

func (c *cleaner) Empty() bool {
	return len(c.entries) == 0
}

func (c *cleaner) Run(ctx context.Context) []error {
	if c.delay != nil {
		time.Sleep(*c.delay)
	}
	var errs []error
	for i := len(c.entries) - 1; i >= 0; i-- {
		if err := c.delete(ctx, c.entries[i]); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *cleaner) delete(ctx context.Context, entry cleanupEntry) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if err := entry.client.Delete(ctx, entry.object, ctrlclient.PropagationPolicy(metav1.DeletePropagationForeground)); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}
	} else if err := client.WaitForDeletion(ctx, entry.client, entry.object); err != nil {
		return err
	}
	return nil
}
