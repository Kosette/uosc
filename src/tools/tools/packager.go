package tools

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"uosc/bins/src/tools/lib"

	"k8s.io/apimachinery/pkg/util/sets"
)

func Packager(args []string) {
	// Display help.
	if len(args) > 0 && sets.New("--help", "-h").Has(args[0]) {
		fmt.Printf(`Packages uosc release files into 'release/' directory, while ensuring binaries inside the zip file are marked as executable even when packaged on windows (otherwise this could've just be a simple .ps1/.sh file).`)
		os.Exit(0)
	}

	cwd := must(os.Getwd())
	releaseRoot := filepath.Join(cwd, "release")
	releaseArchiveSrcDstMap := map[string]string{
		filepath.Join(cwd, "src/fonts"): "",
		filepath.Join(cwd, "src/uosc"):  "scripts/",
	}
	releaseConfigPath := filepath.Join(releaseRoot, "uosc.conf")
	releaseArchivePath := filepath.Join(releaseRoot, "uosc.zip")
	sourceConfigPath := filepath.Join(cwd, "src/uosc.conf")

	// Cleanup old release.
	check(os.RemoveAll(releaseRoot))

	// Package new release
	stats := must(lib.ZipFilesWithHeaders(releaseArchiveSrcDstMap, releaseArchivePath, nil))

	// Copy config to release folder for convenience.
	configFileSrc := must(os.Open(sourceConfigPath))
	configFileDst := must(os.Create(releaseConfigPath))
	confSize := must(io.Copy(configFileDst, configFileSrc))

	fmt.Printf(
		"Packaging into: %s\n- uosc.zip:  %.2f MB, %d files\n- uosc.conf: %.1f KB",
		filepath.ToSlash(must(filepath.Rel(cwd, releaseRoot)))+"/",
		float64(stats.CompressedBytes)/1024/1024,
		stats.FilesNum,
		float64(confSize)/1024,
	)
}
