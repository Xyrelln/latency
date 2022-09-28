package fs

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/leaanthony/slicer"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

// GetTimeStamp return timestamp string
func GetTimeStamp() string {
	return time.Now().Format("20060102150405.000")
}

func GetExecuteRoot() (string, error) {
	p, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return filepath.Dir(p), nil
}

func CreateWorkDir() (string, string) {
	root, _ := GetExecuteRoot()
	timestamp := GetTimeStamp()
	workDir := filepath.Join(root, "cache", timestamp)
	videoDir := filepath.Join(workDir, "video")
	imagesDir := filepath.Join(workDir, "images")

	if _, err := os.Stat(videoDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(videoDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(imagesDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(imagesDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	return videoDir, imagesDir
}

func ClearCacheDir() {
	root, _ := GetExecuteRoot()
	workDir := filepath.Join(root, "cache")
	go os.RemoveAll(workDir)
}

// isWindowsDrivePath returns true if the file path is of the form used by
// Windows. We check if the path begins with a drive letter, followed by a ":".
// For example: C:/x/y/z.
func IsWindowsDrivePath(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[0])) && path[1] == ':'
}

func IsWindowsDrivePathURI(path string) bool {
	if len(path) < 3 {
		return false
	}
	return unicode.IsLetter(rune(path[1])) && path[2] == ':'
}

func GetImageFiles(pathname string, s []string) ([]string, error) {
	// var imgs []string
	rd, err := os.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "png") {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	// sorted
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j] // filename as 0001   0002
	})

	return s, nil
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}


// RelativeToCwd returns an absolute path based on the cwd
// and the given relative path
func RelativeToCwd(relativePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, relativePath), nil
}

// Mkdir will create the given directory
func Mkdir(dirname string) error {
	return os.Mkdir(dirname, 0755)
}

// MkDirs creates the given nested directories.
// Returns error on failure
func MkDirs(fullPath string, mode ...os.FileMode) error {
	var perms os.FileMode
	perms = 0755
	if len(mode) == 1 {
		perms = mode[0]
	}
	return os.MkdirAll(fullPath, perms)
}

// MoveFile attempts to move the source file to the target
// Target is a fully qualified path to a file *name*, not a
// directory
func MoveFile(source string, target string) error {
	return os.Rename(source, target)
}

// DeleteFile will delete the given file
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// CopyFile from source to target
func CopyFile(source string, target string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

// DirExists - Returns true if the given path resolves to a directory on the filesystem
func DirExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

// FileExists returns a boolean value indicating whether
// the given file exists
func FileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}

// RelativePath returns a qualified path created by joining the
// directory of the calling file and the given relative path.
//
// Example: RelativePath("..") in *this* file would give you '/path/to/wails2/v2/internal`
func RelativePath(relativepath string, optionalpaths ...string) string {
	_, thisFile, _, _ := runtime.Caller(1)
	localDir := filepath.Dir(thisFile)

	// If we have optional paths, join them to the relativepath
	if len(optionalpaths) > 0 {
		paths := []string{relativepath}
		paths = append(paths, optionalpaths...)
		relativepath = filepath.Join(paths...)
	}
	result, err := filepath.Abs(filepath.Join(localDir, relativepath))
	if err != nil {
		// I'm allowing this for 1 reason only: It's fatal if the path
		// supplied is wrong as it's only used internally in Wails. If we get
		// that path wrong, we should know about it immediately. The other reason is
		// that it cuts down a ton of unnecassary error handling.
		log.Error(err)
	}
	return result
}

// MustLoadString attempts to load a string and will abort with a fatal message if
// something goes wrong
func MustLoadString(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("FATAL: Unable to load file '%s': %s\n", filename, err.Error())
		os.Exit(1)
	}
	return *(*string)(unsafe.Pointer(&data))
}

// MD5File returns the md5sum of the given file
func MD5File(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}


// fatal will print the optional messages and die
func fatal(message ...string) {
	if len(message) > 0 {
		print("FATAL:")
		for text := range message {
			print(text)
		}
	}
	os.Exit(1)
}

// GetSubdirectories returns a list of subdirectories for the given root directory
func GetSubdirectories(rootDir string) (*slicer.StringSlicer, error) {
	var result slicer.StringSlicer

	// Iterate root dir
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If we have a directory, save it
		if info.IsDir() {
			result.Add(path)
		}
		return nil
	})
	return &result, err
}

func DirIsEmpty(dir string) (bool, error) {

	// CREDIT: https://stackoverflow.com/a/30708914/8325411
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
// Credit: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = MkDirs(dst)
	if err != nil {
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// SetPermissions recursively sets file permissions on a directory
func SetPermissions(dir string, perm os.FileMode) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, perm)
	})
}

// CopyDirExtended recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist. It ignores any files or
// directories that are given through the ignore parameter.
// Symlinks are ignored and skipped.
// Credit: https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyDirExtended(src string, dst string, ignore []string) (err error) {

	ignoreList := slicer.String(ignore)
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = MkDirs(dst)
	if err != nil {
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if ignoreList.Contains(entry.Name()) {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func FindPathToFile(fsys fs.FS, file string) (string, error) {
	stat, _ := fs.Stat(fsys, file)
	if stat != nil {
		return ".", nil
	}
	var indexFiles slicer.StringSlicer
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, file) {
			indexFiles.Add(path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if indexFiles.Length() > 1 {
		selected := indexFiles.AsSlice()[0]
		for _, f := range indexFiles.AsSlice() {
			if len(f) < len(selected) {
				selected = f
			}
		}
		path, _ := filepath.Split(selected)
		return path, nil
	}
	if indexFiles.Length() > 0 {
		path, _ := filepath.Split(indexFiles.AsSlice()[0])
		return path, nil
	}
	return "", fmt.Errorf("no index.html found")
}
