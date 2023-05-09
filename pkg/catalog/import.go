package catalog

import (
	"context"
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/kv"
	"github.com/treeverse/lakefs/pkg/logging"
	"google.golang.org/protobuf/proto"
	"os"
	"time"
)

const (
	ImportCanceled = "Canceled"
	statusChanSize = 100
)

type ImportManager interface {
	Set(record EntryRecord) error
	Ingest(it *walkEntryIterator) error
	Status() graveler.ImportStatus
	NewItr() *importIterator
	Close()
}

type Import struct {
	db            *pebble.DB
	dbPath        string
	kvStore       kv.Store
	status        graveler.ImportStatus
	StatusChan    chan graveler.ImportStatus
	logger        logging.Logger
	repoPartition string
	progress      int64
}

func NewImport(ctx context.Context, cancel context.CancelFunc, kvStore kv.Store, repository *graveler.RepositoryRecord, importID string) (*Import, error) {
	status := graveler.ImportStatus{
		ID:        graveler.ImportID(importID),
		UpdatedAt: time.Now(),
	}
	repoPartition := graveler.RepoPartition(repository)
	// Must be set first
	err := kv.SetMsg(ctx, kvStore, repoPartition, []byte(importID), graveler.ProtoFromImportStatus(&status))
	if err != nil {
		return nil, err
	}
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("import_%s", importID))
	if err != nil {
		return nil, err
	}
	importDB, err := pebble.Open(tmpDir, nil)
	if err != nil {
		return nil, err
	}

	i := Import{
		db:            importDB,
		dbPath:        tmpDir,
		kvStore:       kvStore,
		status:        status,
		StatusChan:    make(chan graveler.ImportStatus, statusChanSize),
		logger:        logging.Default().WithField("import_id", importID),
		repoPartition: repoPartition,
	}

	go func() {
		err := i.importStatusAsync(ctx, cancel)
		if err != nil {
			i.logger.WithError(err).Error("failed to update import status")
		}
	}()
	return &i, nil
}

func (i *Import) Set(record EntryRecord) error {
	key := []byte(record.Path)
	value, err := EntryToValue(record.Entry)
	if err != nil {
		return err
	}
	pb := graveler.ProtoFromStagedEntry(key, value)
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return i.db.Set(key, data, &pebble.WriteOptions{
		Sync: false,
	})
}

func (i *Import) Ingest(it *walkEntryIterator) error {
	for it.Next() {
		if err := i.Set(*it.Value()); err != nil {
			return err
		}
		i.progress += 1
	}
	return nil
}

func (i *Import) Status() graveler.ImportStatus {
	return i.status
}

func (i *Import) NewItr() *importIterator {
	return newImportIterator(i.db.NewIter(nil))
}

func (i *Import) Close() {
	close(i.StatusChan)
	_ = os.RemoveAll(i.dbPath)
}

func (i *Import) importStatusAsync(ctx context.Context, cancel context.CancelFunc) error {
	const statusUpdateInterval = 1 * time.Second
	statusData := graveler.ImportStatusData{}
	newStatus := i.status
	timer := time.NewTimer(statusUpdateInterval)
	done := false
	for {
		select {
		case <-ctx.Done():
			done = true
		case s, ok := <-i.StatusChan:
			if ok {
				newStatus = s
			}
		case <-timer.C:
			pred, err := kv.GetMsg(ctx, i.kvStore, i.repoPartition, []byte(i.status.ID), &statusData)
			if err != nil {
				cancel()
				return err
			}
			if statusData.Error != "" || statusData.Completed {
				cancel()
				return nil
			}
			newStatus.UpdatedAt = time.Now()
			newStatus.Progress = i.progress
			// TODO: Remove key after 24 hrs
			err = kv.SetMsgIf(ctx, i.kvStore, i.repoPartition, []byte(i.status.ID), graveler.ProtoFromImportStatus(&newStatus), pred)
			if errors.Is(err, kv.ErrPredicateFailed) {
				i.logger.Warning("Import status update failed")
			} else if err != nil {
				cancel()
				return err
			}
			i.status = newStatus

			// Check if context was canceled, we want to do that only after we update the import state
			if done {
				return nil
			}
			timer.Reset(statusUpdateInterval)
		}
	}
}

type importIterator struct {
	it      *pebble.Iterator
	hasMore bool
	value   *graveler.ValueRecord
	err     error
}

func newImportIterator(it *pebble.Iterator) *importIterator {
	itr := &importIterator{
		it: it,
	}
	itr.hasMore = it.First()
	return itr
}

func (it *importIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if !it.hasMore {
		it.value = nil
		return false
	}
	var value []byte
	value, it.err = it.it.ValueAndErr()
	if it.err != nil {
		return false
	}
	data := graveler.StagedEntryData{}
	it.err = proto.Unmarshal(value, &data)
	if it.err != nil {
		return false
	}
	it.value = &graveler.ValueRecord{
		Key:   data.Key,
		Value: graveler.StagedEntryFromProto(&data),
	}
	it.hasMore = it.it.Next()
	return true
}

func (it *importIterator) SeekGE(id graveler.Key) {
	it.value = nil
	it.hasMore = it.it.SeekGE(id)
}

func (it *importIterator) Value() *graveler.ValueRecord {
	return it.value
}

func (it *importIterator) Err() error {
	return it.err
}

func (it *importIterator) Close() {
	_ = it.it.Close()
}
