package main

import (
	"hookfs"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// MyHookContext implements hookfs.HookContext
type MyHookContext struct {
	path string
}

// MyHook implements hookfs.Hook
//type MyHook struct{}

type MyHook struct {
	faultType int
	percent int
	delay time.Duration
}

// Init implements hookfs.HookWithInit
func (h *MyHook) Init() error {
	log.WithFields(log.Fields{
		"h": h,
	}).Info("MyInit: initializing")
	return nil
}

// PreOpen implements hookfs.HookOnOpen
func (h *MyHook) PreOpen(path string, flags uint32) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == OpenFileEIO {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpen: returning EIO")
		return true, ctx, syscall.EIO
	}
	return false, ctx, nil
}

// PostOpen implements hookfs.HookOnOpen
func (h *MyHook) PostOpen(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == OpenFileEPERM {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpen: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreRead implements hookfs.HookOnRead
func (h *MyHook) PreRead(path string, length int64, offset int64) ([]byte, bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == ReadFileDelay {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": h.delay,
		}).Info("MyPreRead: sleeping")
		time.Sleep(h.delay)
	}
	return nil, false, ctx, nil
}

// PostRead implements hookfs.HookOnRead
func (h *MyHook) PostRead(realRetCode int32, realBuf []byte, ctx hookfs.HookContext) ([]byte, bool, error) {
	percentage := 0
	if h.faultType == ReadFileErr {
		percentage = h.percent
	}
	if probab(percentage) {
		buf := []byte("Hello FaultFS hooked error Data!\n")
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
			"buf": buf,
		}).Info("MyPostRead: returning injected buffer")
		return buf, true, nil
	}
	return nil, false, nil
}

// PreWrite implements hookfs.HookOnWrite
func (h *MyHook) PreWrite(path string, buf []byte, offset int64) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == WriteFileDelay {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": h.delay,
		}).Info("MyPreWrite: sleeping")
		time.Sleep(h.delay)
	}
	return false, ctx, nil
}

// PostWrite implements hookfs.HookOnWrite
func (h *MyHook) PostWrite(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == WriteFileENOSPC {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostWrite: returning ENOSPC")
		return true, syscall.ENOSPC
	}
	return false, nil
}

// PreMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PreMkdir(path string, mode uint32) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == MkDirEACCES {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreMkdir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PostMkdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == MkDirEPERM {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostMkdir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PreRmdir(path string) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == RmDirEACCESS {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreRmdir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PostRmdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == RmDirEPERM {
		percentage = h.percent
	}
	if probab(percentage) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostRmdir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PreOpenDir(path string) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == OpenDirEACCESS {
		percentage = h.percent
	}
	if probab(percentage) && path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpenDir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PostOpenDir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == OpenDirEPERM {
		percentage = h.percent
	}
	if probab(percentage) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpenDir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreFsync implements hookfs.HookOnFsync
func (h *MyHook) PreFsync(path string, flags uint32) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	percentage := 0
	if h.faultType == FsycnDelay {
		percentage = h.percent
	}
	if probab(percentage) && path != "" {
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": h.delay,
		}).Info("MyPreFsync: sleeping")
		time.Sleep(h.delay)
	}
	return false, ctx, nil
}

// PostFsync implements hookfs.HookOnFsync
func (h *MyHook) PostFsync(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	percentage := 0
	if h.faultType == FsycnEIO {
		percentage = h.percent
	}
	if probab(percentage) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostFsync: returning EIO")
		return true, syscall.EIO
	}
	return false, nil
}
