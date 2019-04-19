package main

const (
	OpenFileEIO = iota
	OpenFileEPERM

	ReadFileDelay
	ReadFileErr

	WriteFileENOSPC
	WriteFileDelay

	MkDirEACCES
	MkDirEPERM

	RmDirEACCESS
	RmDirEPERM

	FsycnDelay
	FsycnEIO

	OpenDirEACCESS
	OpenDirEPERM

)
