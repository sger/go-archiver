package archiver

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// Archiver Interface
type Archiver interface {
	Name() string
	Archive(src, dest string) error
	Restore(src, dest string) error
}

// Name describes the name of the zip file
func (z *Zip) Name() string {
	return "%d.zip"
}

// Archive a directory 2 parameters required the source file and the destination
func (z *Zip) Archive(src, dest string) error {
	// check if directory exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}

	out, err := os.Create(dest)

	if err != nil {
		return err
	}

	defer out.Close()

	w := zip.NewWriter(out)

	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		in, err := os.Open(path)

		if err != nil {
			return err
		}

		defer in.Close()

		f, err := w.Create(path)

		if err != nil {
			return err
		}

		io.Copy(f, in)
		return nil
	})
}

// Restore a directory 2 paramaters required the source file and the destination
func (z *Zip) Restore(src, dest string) error {
	r, err := zip.OpenReader(src)

	if err != nil {
		return err
	}

	defer r.Close()

	var w sync.WaitGroup
	var errs []error

	errChan := make(chan error)

	go func() {
		for err := range errChan {
			errs = append(errs, err)
		}
	}()

	for _, f := range r.File {
		w.Add(1)
		go func(f *zip.File) {
			zippedFile, err := f.Open()
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			toFilename := path.Join(dest, f.Name)
			err = os.MkdirAll(path.Dir(toFilename), 0777)
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			newFile, err := os.Create(toFilename)
			if err != nil {
				zippedFile.Close()
				errChan <- err
				w.Done()
				return
			}
			_, err = io.Copy(newFile, zippedFile)
			newFile.Close()
			zippedFile.Close()
			if err != nil {
				errChan <- err
				w.Done()
				return
			}
			w.Done()
		}(f)
	}
	w.Wait()
	close(errChan)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
