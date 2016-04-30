package archiver_test

import (
	"os"
	"testing"

	"github.com/sger/archiver"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) {
	os.MkdirAll("test/output", 0777)
}

func teardown(t *testing.T) {
	os.RemoveAll("test/output")
}

func TestZipArchive(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := archiver.GetInstance().Archive("test/files", "test/output/files.zip")
	require.NoError(t, err)

	err = archiver.GetInstance().Restore("test/output/files.zip", "test/output/restored")
	require.NoError(t, err)
}

type call struct {
	Src  string
	Dest string
}

type TestArchiver struct {
	Archives []*call
	Restores []*call
}

var _ archiver.Archiver = (*TestArchiver)(nil)

func (a *TestArchiver) Name() string {
	return "%d.zip"
}

func (a *TestArchiver) Archive(src, dest string) error {
	a.Archives = append(a.Archives, &call{Src: src, Dest: dest})
	return nil
}

func (a *TestArchiver) Restore(src, dest string) error {
	a.Restores = append(a.Restores, &call{Src: src, Dest: dest})
	return nil
}
