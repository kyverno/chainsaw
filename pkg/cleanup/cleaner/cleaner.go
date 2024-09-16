package cleaner

import (
	"context"
	"time"

	"github.com/kyverno/chainsaw/pkg/client"
	"github.com/kyverno/chainsaw/pkg/model"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Run(ctx context.Context, stepReport *model.StepReport) []error
}

func New(timeout time.Duration, delay *time.Duration, propagation metav1.DeletionPropagation) Cleaner {
	return &cleaner{
		delay:       delay,
		timeout:     timeout,
		propagation: propagation,
	}
}

type cleaner struct {
	delay       *time.Duration
	timeout     time.Duration
	propagation metav1.DeletionPropagation
	entries     []cleanupEntry
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

func (c *cleaner) Run(ctx context.Context, stepReport *model.StepReport) []error {
	if c.delay != nil {
		time.Sleep(*c.delay)
	}
	var errs []error
	for i := len(c.entries) - 1; i >= 0; i-- {
		report := model.OperationReport{
			Type:      model.OperationTypeDelete,
			StartTime: time.Now(),
		}
		if report.Err = c.delete(ctx, c.entries[i]); report.Err != nil {
			errs = append(errs, report.Err)
		}
		report.EndTime = time.Now()
		if stepReport != nil {
			stepReport.Add(&report)
		}
	}
	return errs
}

func (c *cleaner) delete(ctx context.Context, entry cleanupEntry) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	if err := entry.client.Delete(ctx, entry.object, client.PropagationPolicy(c.propagation)); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}
	} else if err := client.WaitForDeletion(ctx, entry.client, entry.object); err != nil {
		return err
	}
	return nil
}
