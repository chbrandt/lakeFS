package catalog_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/golang/mock/gomock"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/pkg/block"
	"github.com/treeverse/lakefs/pkg/catalog"
	cUtils "github.com/treeverse/lakefs/pkg/catalog/testutils"
	"github.com/treeverse/lakefs/pkg/graveler"
	gUtils "github.com/treeverse/lakefs/pkg/graveler/testutil"
	"github.com/treeverse/lakefs/pkg/testutil"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetStartPos(t *testing.T) {
	/**
	Here are the possible states for where to start iterating, given a prefix, an "after" and a delimiter:
	delim	prefix	after	expected start position
	0		0		0		""
	0		0		1		at after (when iterating we'll skip the exact match anyway)
	0		1		0		at prefix
	0		1		1		at after (when iterating we'll skip the exact match anyway)
	1		0		0		""
	1		0		1		next common prefix after "after"
	1		1		0		at prefix
	1		1		1		next common prefix after "after"
	*/
	cases := []struct {
		Name      string
		Delimiter string
		Prefix    string
		After     string
		Expected  string
	}{
		{"all_empty", "", "", "", ""},
		{"only_after", "", "", "a/b/", "a/b/"},
		{"only_prefix", "", "a/", "", "a/"},
		{"prefix_and_after", "", "a/", "a/b/", "a/b/"},
		{"only_delimiter", "/", "", "", ""},
		{"delimiter_and_after", "/", "", "a/b/", string(graveler.UpperBoundForPrefix([]byte("a/b/")))},
		{"delimiter_and_prefix", "/", "a/", "", "a/"},
		{"delimiter_prefix_and_after", "/", "a/", "a/b/", string(graveler.UpperBoundForPrefix([]byte("a/b/")))},
		{"empty_directory", "/", "", "a/b//", string(graveler.UpperBoundForPrefix([]byte("a/b//")))},
		{"after_before_prefix", "", "c", "a", "c"},
		{"after_before_prefix_with_delim", "/", "c/", "a", "c/"},
	}

	for _, cas := range cases {
		t.Run(cas.Name, func(t *testing.T) {
			got := catalog.GetStartPos(cas.Prefix, cas.After, cas.Delimiter)
			if got != cas.Expected {
				t.Fatalf("expected %s got %s", cas.Expected, got)
			}
		})
	}
}

func TestCatalog_ListRepositories(t *testing.T) {
	// prepare data tests
	now := time.Now()
	gravelerData := []*graveler.RepositoryRecord{
		{RepositoryID: "repo1", Repository: &graveler.Repository{StorageNamespace: "storage1", CreationDate: now, DefaultBranchID: "main1"}},
		{RepositoryID: "repo2", Repository: &graveler.Repository{StorageNamespace: "storage2", CreationDate: now, DefaultBranchID: "main2"}},
		{RepositoryID: "repo3", Repository: &graveler.Repository{StorageNamespace: "storage3", CreationDate: now, DefaultBranchID: "main3"}},
	}
	type args struct {
		limit int
		after string
	}
	tests := []struct {
		name        string
		args        args
		want        []*catalog.Repository
		wantHasMore bool
		wantErr     bool
	}{
		{
			name: "all",
			args: args{
				limit: -1,
				after: "",
			},
			want: []*catalog.Repository{
				{Name: "repo1", StorageNamespace: "storage1", DefaultBranch: "main1", CreationDate: now},
				{Name: "repo2", StorageNamespace: "storage2", DefaultBranch: "main2", CreationDate: now},
				{Name: "repo3", StorageNamespace: "storage3", DefaultBranch: "main3", CreationDate: now},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "first",
			args: args{
				limit: 1,
				after: "",
			},
			want: []*catalog.Repository{
				{Name: "repo1", StorageNamespace: "storage1", DefaultBranch: "main1", CreationDate: now},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "second",
			args: args{
				limit: 1,
				after: "repo1",
			},
			want: []*catalog.Repository{
				{Name: "repo2", StorageNamespace: "storage2", DefaultBranch: "main2", CreationDate: now},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "last2",
			args: args{
				limit: 10,
				after: "repo1",
			},
			want: []*catalog.Repository{
				{Name: "repo2", StorageNamespace: "storage2", DefaultBranch: "main2", CreationDate: now},
				{Name: "repo3", StorageNamespace: "storage3", DefaultBranch: "main3", CreationDate: now},
			},
			wantHasMore: false,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup Catalog
			gravelerMock := &catalog.FakeGraveler{
				RepositoryIteratorFactory: catalog.NewFakeRepositoryIteratorFactory(gravelerData),
			}
			c := &catalog.Catalog{
				Store: gravelerMock,
			}
			// test method
			ctx := context.Background()
			got, hasMore, err := c.ListRepositories(ctx, tt.args.limit, "", tt.args.after)
			if tt.wantErr && err == nil {
				t.Fatal("ListRepositories err nil, expected error")
			}
			if err != nil {
				return
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("ListRepositories hasMore %t, expected %t", hasMore, tt.wantHasMore)
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error("ListRepositories diff found:", diff)
			}
		})
	}
}

func TestCatalog_BranchExists(t *testing.T) {
	// prepare branch data
	gravelerData := []*graveler.BranchRecord{
		{BranchID: "branch1", Branch: &graveler.Branch{CommitID: "commit1"}},
		{BranchID: "branch2", Branch: &graveler.Branch{CommitID: "commit2"}},
		{BranchID: "branch3", Branch: &graveler.Branch{CommitID: "commit3"}},
	}
	tests := []struct {
		Branch      string
		ShouldExist bool
	}{{"branch1", true}, {"branch2", true}, {"branch-foo", false}}
	for _, tt := range tests {
		t.Run(tt.Branch, func(t *testing.T) {
			// setup Catalog
			gravelerMock := &catalog.FakeGraveler{
				BranchIteratorFactory: gUtils.NewFakeBranchIteratorFactory(gravelerData),
			}
			c := &catalog.Catalog{
				Store: gravelerMock,
			}
			// test method
			ctx := context.Background()
			exists, err := c.BranchExists(ctx, "repo", tt.Branch)
			if err != nil {
				t.Fatal("BranchExists failed:", err)
			}
			if exists != tt.ShouldExist {
				not := ""
				if !tt.ShouldExist {
					not = " not"
				}
				t.Errorf("branch %s should%s exist", tt.Branch, not)
			}
		})
	}
}

func TestCatalog_ListBranches(t *testing.T) {
	// prepare branch data
	gravelerData := []*graveler.BranchRecord{
		{BranchID: "branch1", Branch: &graveler.Branch{CommitID: "commit1"}},
		{BranchID: "branch2", Branch: &graveler.Branch{CommitID: "commit2"}},
		{BranchID: "branch3", Branch: &graveler.Branch{CommitID: "commit3"}},
	}
	type args struct {
		prefix string
		limit  int
		after  string
	}
	tests := []struct {
		name        string
		args        args
		want        []*catalog.Branch
		wantHasMore bool
		wantErr     bool
	}{
		{
			name: "all",
			args: args{limit: -1},
			want: []*catalog.Branch{
				{Name: "branch1", Reference: "commit1"},
				{Name: "branch2", Reference: "commit2"},
				{Name: "branch3", Reference: "commit3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "exact",
			args: args{limit: 3},
			want: []*catalog.Branch{
				{Name: "branch1", Reference: "commit1"},
				{Name: "branch2", Reference: "commit2"},
				{Name: "branch3", Reference: "commit3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "first",
			args: args{limit: 1},
			want: []*catalog.Branch{
				{Name: "branch1", Reference: "commit1"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "second",
			args: args{limit: 1, after: "branch1"},
			want: []*catalog.Branch{
				{Name: "branch2", Reference: "commit2"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "last2",
			args: args{limit: 10, after: "branch1"},
			want: []*catalog.Branch{
				{Name: "branch2", Reference: "commit2"},
				{Name: "branch3", Reference: "commit3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name:        "not found",
			args:        args{limit: 10, after: "zzz"},
			wantHasMore: false,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup Catalog
			gravelerMock := &catalog.FakeGraveler{
				BranchIteratorFactory: gUtils.NewFakeBranchIteratorFactory(gravelerData),
			}
			c := &catalog.Catalog{
				Store: gravelerMock,
			}
			// test method
			ctx := context.Background()
			got, hasMore, err := c.ListBranches(ctx, "repo", tt.args.prefix, tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListBranches() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error("ListBranches() found diff", diff)
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("ListBranches() hasMore = %t, want %t", hasMore, tt.wantHasMore)
			}
		})
	}
}

func TestCatalog_ListTags(t *testing.T) {
	gravelerData := []*graveler.TagRecord{
		{TagID: "t1", CommitID: "c1"},
		{TagID: "t2", CommitID: "c2"},
		{TagID: "t3", CommitID: "c3"},
	}
	type args struct {
		limit int
		after string
	}
	tests := []struct {
		name        string
		args        args
		want        []*catalog.Tag
		wantHasMore bool
		wantErr     bool
	}{
		{
			name: "all",
			args: args{limit: -1},
			want: []*catalog.Tag{
				{ID: "t1", CommitID: "c1"},
				{ID: "t2", CommitID: "c2"},
				{ID: "t3", CommitID: "c3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "exact",
			args: args{limit: 3},
			want: []*catalog.Tag{
				{ID: "t1", CommitID: "c1"},
				{ID: "t2", CommitID: "c2"},
				{ID: "t3", CommitID: "c3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "first",
			args: args{limit: 1},
			want: []*catalog.Tag{
				{ID: "t1", CommitID: "c1"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "second",
			args: args{limit: 1, after: "t1"},
			want: []*catalog.Tag{
				{ID: "t2", CommitID: "c2"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "last2",
			args: args{limit: 10, after: "t1"},
			want: []*catalog.Tag{
				{ID: "t2", CommitID: "c2"},
				{ID: "t3", CommitID: "c3"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name:        "not found",
			args:        args{limit: 10, after: "zzz"},
			wantHasMore: false,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gravelerMock := &catalog.FakeGraveler{
				TagIteratorFactory: catalog.NewFakeTagIteratorFactory(gravelerData),
			}
			c := &catalog.Catalog{
				Store: gravelerMock,
			}
			ctx := context.Background()
			got, hasMore, err := c.ListTags(ctx, "repo", "", tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error("ListTags() found diff", diff)
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("ListTags() hasMore = %t, want %t", hasMore, tt.wantHasMore)
			}
		})
	}
}

func TestCatalog_ListEntries(t *testing.T) {
	// prepare branch data
	now := time.Now()
	gravelerData := []*graveler.ValueRecord{
		{Key: graveler.Key("file1"), Value: catalog.MustEntryToValue(&catalog.Entry{Address: "file1", LastModified: timestamppb.New(now), Size: 1, ETag: "01"})},
		{Key: graveler.Key("file2"), Value: catalog.MustEntryToValue(&catalog.Entry{Address: "file2", LastModified: timestamppb.New(now), Size: 2, ETag: "02"})},
		{Key: graveler.Key("file3"), Value: catalog.MustEntryToValue(&catalog.Entry{Address: "file3", LastModified: timestamppb.New(now), Size: 3, ETag: "03"})},
		{Key: graveler.Key("h/file1"), Value: catalog.MustEntryToValue(&catalog.Entry{Address: "h/file1", LastModified: timestamppb.New(now), Size: 1, ETag: "01"})},
		{Key: graveler.Key("h/file2"), Value: catalog.MustEntryToValue(&catalog.Entry{Address: "h/file2", LastModified: timestamppb.New(now), Size: 2, ETag: "02"})},
	}
	type args struct {
		prefix    string
		after     string
		delimiter string
		limit     int
	}
	tests := []struct {
		name        string
		args        args
		want        []*catalog.DBEntry
		wantHasMore bool
		wantErr     bool
	}{
		{
			name: "all no delimiter",
			args: args{limit: -1},
			want: []*catalog.DBEntry{
				{Path: "file1", PhysicalAddress: "file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
				{Path: "file2", PhysicalAddress: "file2", CreationDate: now, Size: 2, Checksum: "02", ContentType: "application/octet-stream"},
				{Path: "file3", PhysicalAddress: "file3", CreationDate: now, Size: 3, Checksum: "03", ContentType: "application/octet-stream"},
				{Path: "h/file1", PhysicalAddress: "h/file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
				{Path: "h/file2", PhysicalAddress: "h/file2", CreationDate: now, Size: 2, Checksum: "02", ContentType: "application/octet-stream"},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "first no delimiter",
			args: args{limit: 1},
			want: []*catalog.DBEntry{
				{Path: "file1", PhysicalAddress: "file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "second no delimiter",
			args: args{limit: 1, after: "file1"},
			want: []*catalog.DBEntry{
				{Path: "file2", PhysicalAddress: "file2", CreationDate: now, Size: 2, Checksum: "02", ContentType: "application/octet-stream"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "last two no delimiter",
			args: args{limit: 10, after: "file3"},
			want: []*catalog.DBEntry{
				{Path: "h/file1", PhysicalAddress: "h/file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
				{Path: "h/file2", PhysicalAddress: "h/file2", CreationDate: now, Size: 2, Checksum: "02", ContentType: "application/octet-stream"},
			},
			wantHasMore: false,
			wantErr:     false,
		},

		{
			name: "all with delimiter",
			args: args{limit: -1, delimiter: "/"},
			want: []*catalog.DBEntry{
				{Path: "file1", PhysicalAddress: "file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
				{Path: "file2", PhysicalAddress: "file2", CreationDate: now, Size: 2, Checksum: "02", ContentType: "application/octet-stream"},
				{Path: "file3", PhysicalAddress: "file3", CreationDate: now, Size: 3, Checksum: "03", ContentType: "application/octet-stream"},
				{Path: "h/", CommonLevel: true},
			},
			wantHasMore: false,
			wantErr:     false,
		},
		{
			name: "first with delimiter",
			args: args{limit: 1, delimiter: "/"},
			want: []*catalog.DBEntry{
				{Path: "file1", PhysicalAddress: "file1", CreationDate: now, Size: 1, Checksum: "01", ContentType: "application/octet-stream"},
			},
			wantHasMore: true,
			wantErr:     false,
		},
		{
			name: "last with delimiter",
			args: args{limit: 1, after: "file3", delimiter: "/"},
			want: []*catalog.DBEntry{
				{Path: "h/", CommonLevel: true},
			},
			wantHasMore: false,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup Catalog
			gravelerMock := &catalog.FakeGraveler{
				ListIteratorFactory: catalog.NewFakeValueIteratorFactory(gravelerData),
			}
			c := &catalog.Catalog{
				Store: gravelerMock,
			}
			// test method
			ctx := context.Background()
			got, hasMore, err := c.ListEntries(ctx, "repo", "ref", tt.args.prefix, tt.args.after, tt.args.delimiter, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListEntries() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Error("ListEntries() diff found", diff)
			}
			if hasMore != tt.wantHasMore {
				t.Errorf("ListEntries() hasMore = %t, want %t", hasMore, tt.wantHasMore)
			}
		})
	}
}

var (
	repoID     = graveler.RepositoryID("repo1")
	repository = &graveler.RepositoryRecord{
		RepositoryID: repoID,
		Repository: &graveler.Repository{
			StorageNamespace: "mock-sn",
			CreationDate:     time.Now(),
			DefaultBranchID:  "main",
		},
	}
)

func TestCatalog_PrepareGCUncommitted(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		numBranch     int
		numRecords    int
		expectedCalls int
	}{
		{
			name:          "no branches",
			numBranch:     0,
			numRecords:    0,
			expectedCalls: 1,
		},
		{
			name:          "no objects",
			numBranch:     3,
			numRecords:    0,
			expectedCalls: 1,
		},
		{
			name:          "Sanity",
			numBranch:     5,
			numRecords:    3,
			expectedCalls: 1,
		},
		{
			name:          "Tokenized",
			numBranch:     500,
			numRecords:    500,
			expectedCalls: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := createPrepareUncommittedTestScenario(t, tt.numBranch, tt.numRecords, tt.expectedCalls)
			blockAdapter := testutil.NewBlockAdapterByType(t, block.BlockstoreTypeMem)
			c := &catalog.Catalog{
				Store:                    g.Sut,
				BlockAdapter:             blockAdapter,
				GCMaxUncommittedFileSize: 500 * 1024,
			}
			var result *catalog.PrepareGCUncommittedInfo

			result, err := c.PrepareGCUncommitted(ctx, repoID.String(), nil)
			require.NoError(t, err)

			for result.Mark != nil {
				runID := result.Metadata.RunId
				require.Equal(t, runID, result.Mark.RunID)
				result, err = c.PrepareGCUncommitted(ctx, repoID.String(), result.Mark)
				require.NoError(t, err)
				require.Equal(t, runID, result.Metadata.RunId)
			}
			verifyData(t, ctx, tt.numBranch, tt.numRecords, result.Metadata.RunId, c)
		})
	}
}

func createPrepareUncommittedTestScenario(t *testing.T, numBranches, numRecords, expectedCalls int) *gUtils.GravelerTest {
	t.Helper()
	testFolder, err := os.MkdirTemp("", xid.New().String())
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(testFolder)
	})

	test := gUtils.InitGravelerTest(t)
	records := make([][]*graveler.ValueRecord, numBranches)
	var branches []*graveler.BranchRecord
	for i := 0; i < numBranches; i++ {
		branchID := graveler.BranchID(fmt.Sprintf("branch%04d", i))
		token := graveler.StagingToken(fmt.Sprintf("%s_st%04d", branchID, i))
		branches = append(branches, &graveler.BranchRecord{BranchID: branchID, Branch: &graveler.Branch{StagingToken: token}})

		records[i] = make([]*graveler.ValueRecord, numRecords+2)
		for j := 0; j < numRecords; j++ {
			e := catalog.Entry{
				Address:      fmt.Sprintf("%s_record%04d", branchID, j),
				LastModified: timestamppb.New(time.Now()),
				Size:         0,
				ETag:         "",
				Metadata:     nil,
				AddressType:  catalog.Entry_RELATIVE,
				ContentType:  "",
			}
			v, err := proto.Marshal(&e)
			require.NoError(t, err)
			records[i][j] = &graveler.ValueRecord{
				Key: []byte(e.Address),
				Value: &graveler.Value{
					Identity: []byte("dont care"),
					Data:     v,
				},
			}
		}
		// Add tombstone
		records[i][numRecords] = &graveler.ValueRecord{
			Key:   []byte("AtombstoneFile"),
			Value: nil,
		}
		// Add external address
		e := catalog.Entry{
			Address:      fmt.Sprintf("external/address/object"),
			LastModified: timestamppb.New(time.Now()),
			Size:         0,
			ETag:         "",
			Metadata:     nil,
			AddressType:  catalog.Entry_FULL,
			ContentType:  "",
		}
		v, err := proto.Marshal(&e)
		require.NoError(t, err)
		records[i][numRecords+1] = &graveler.ValueRecord{
			Key: []byte(e.Address),
			Value: &graveler.Value{
				Identity: []byte("dont care"),
				Data:     v,
			},
		}
	}
	test.GarbageCollectionManager.EXPECT().NewID().Return("TestRunID")
	test.RefManager.EXPECT().GetRepository(gomock.Any(), repoID).Times(expectedCalls).Return(repository, nil)
	test.RefManager.EXPECT().ListBranches(gomock.Any(), gomock.Any()).Times(expectedCalls).Return(gUtils.NewFakeBranchIterator(branches), nil)

	for i := 0; i < len(branches); i++ {
		test.StagingManager.EXPECT().List(gomock.Any(), branches[i].StagingToken, nil, gomock.Any()).AnyTimes().Return(cUtils.NewFakeValueIterator(records[i]), nil)
	}
	if numRecords*numBranches > 0 {
		test.GarbageCollectionManager.EXPECT().GetUncommittedLocation(gomock.Any(), gomock.Any()).Times(expectedCalls).DoAndReturn(func(runID string, sn graveler.StorageNamespace) (string, error) {
			return fmt.Sprintf("%s/retention/gc/uncommitted/%s/uncommitted/", "_lakefs", runID), nil
		})
	}

	return test
}

func verifyData(t *testing.T, ctx context.Context, numBranches, numRecords int, runID string, c *catalog.Catalog) {
	totalCount := 0
	location := fmt.Sprintf("%s/retention/gc/uncommitted/%s/uncommitted/", "_lakefs", runID)

	err := c.BlockAdapter.Walk(context.Background(), block.WalkOpts{
		StorageNamespace: repository.Repository.StorageNamespace.String(),
		Prefix:           location,
	}, func(id string) error {
		fd, err := c.BlockAdapter.Get(ctx, block.ObjectPointer{
			Identifier:     id,
			IdentifierType: block.IdentifierTypeFull,
		}, 0)
		if err != nil {
			return err
		}
		filename := path.Join(os.TempDir(), xid.New().String())
		fr, err := local.NewLocalFileWriter(filename)
		if err != nil {
			return err
		}
		defer func() {
			fr.Close()
			os.Remove(filename)
		}()

		_, err = io.Copy(fr, fd)
		if err != nil {
			return err
		}

		_, err = fr.Seek(0, 0)
		if err != nil {
			return err
		}

		pr, err := reader.NewParquetReader(fr, new(catalog.UncommittedParquetObject), 4)
		if err != nil {
			return err
		}

		count := pr.GetNumRows()
		u := make([]*catalog.UncommittedParquetObject, count)
		err = pr.Read(&u)
		pr.ReadStop()
		fr.Close()
		for i, record := range u {
			branchID := (i + totalCount) / numRecords
			recordID := (i + totalCount) % numRecords
			require.Equal(t, fmt.Sprintf("branch%04d_record%04d", branchID, recordID), record.PhysicalAddress)
		}
		totalCount += int(count)

		return nil
	})

	require.NoError(t, err)
	require.Equal(t, numBranches*numRecords, totalCount)
}
