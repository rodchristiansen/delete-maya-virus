package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/delete-maya-virus/utils"
)

func main() {
	recursive := flag.Bool("r", false, "recursive path traversal")
	createBackup := flag.Bool("b", false, "make a backup copy of processed files")

	flag.Parse()
	argPath := flag.Arg(0)

	if argPath == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err.Error())
		}
		argPath = currentDir
	}

	fi, err := os.Stat(argPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	var scannedMayaFiles []string

	switch mode := fi.Mode(); {
	case mode.IsDir():
		if *recursive {
			files, err := utils.ReturnMayaFilesFromDirRecursively(argPath)
			if err != nil {
				log.Fatal(err.Error())
			}
			scannedMayaFiles = files
		} else {
			files, err := utils.ReturnMayaFilesFromDir(argPath)
			if err != nil {
				log.Fatal(err.Error())
			}
			scannedMayaFiles = files
		}
	case mode.IsRegular():
		fileExtension := filepath.Ext(argPath)
		if fileExtension != utils.MaFileExt {
			log.Println("Only .ma files can be processed!")
			return
		}
		scannedMayaFiles = append(scannedMayaFiles, argPath)
	}

	for _, pathToMayaScannedFile := range scannedMayaFiles {
		mayaFileName := filepath.Base(pathToMayaScannedFile)
		backupPath := pathToMayaScannedFile + ".backup"

		// Attempt to read the .ma file
		mayaFileData, err := ioutil.ReadFile(pathToMayaScannedFile)
		if err != nil {
			log.Printf("WARNING: Can't open maya file '%s': %s. Skipping to next file.", mayaFileName, err.Error())
			continue // Skip to the next file
		}
		log.Println("Checking the file -", mayaFileName)

		isVaccineVirusInFile := utils.CheckVaccineVirus(mayaFileData)
		isMayaMelVirusInFile := utils.CheckMayaMelVirus(mayaFileData)

		if isVaccineVirusInFile {
			log.Println("Vaccine virus found! Processing file... ", mayaFileName)
			if *createBackup {
				err := utils.CreateBuckup(backupPath, mayaFileData)
				if err != nil {
					log.Printf("ERROR: Failed to create backup for '%s': %s. Skipping to next file.", mayaFileName, err.Error())
					continue // Optionally skip to next file or decide how to handle
				}
			}
			mayaFileData = utils.DeleteVaccineVirus(mayaFileData)
		}

		if isMayaMelVirusInFile {
			log.Println("MayaMelUIConfigurationFile virus found! Processing file... ", mayaFileName)
			if *createBackup {
				err := utils.CreateBuckup(backupPath, mayaFileData)
				if err != nil {
					log.Printf("ERROR: Failed to create backup for '%s': %s. Skipping to next file.", mayaFileName, err.Error())
					continue // Optionally skip to next file or decide how to handle
				}
			}
			mayaFileData = utils.DeleteMayaMelVirus(mayaFileData)
		}

		// Attempt to write the modified .ma file
		writeErr := ioutil.WriteFile(pathToMayaScannedFile, mayaFileData, utils.FilePermission)
		if writeErr != nil {
			log.Printf("ERROR: Failed to write to '%s': %s. Skipping to next file.", mayaFileName, writeErr.Error())
			continue // Skip to the next file
		}
	}
	log.Println("All files are checked.")
}