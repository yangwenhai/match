package misc

// 提供用于处理memprofile和cpuprofile的helper函数
import (
	"fmt"
	"os"
	"runtime/pprof"
)

type Profiler interface {
	Start(f string) error
	Stop()
}

type iCPUProfiler struct {
	running bool
}

func (p *iCPUProfiler) Start(f string) error {
	p.running = false
	if f == "" {
		return fmt.Errorf("empty file name")
	}

	file, err := os.Create(f)
	if err != nil {
		return err
	}

	err = pprof.StartCPUProfile(file)
	if err != nil {
		return err
	}

	p.running = true

	return nil
}

func (p *iCPUProfiler) Stop() {
	if !p.running {
		return
	}

	pprof.StopCPUProfile()
}

type iMEMProfiler struct {
	file *os.File
}

func (p *iMEMProfiler) Start(f string) error {
	if f == "" {
		return fmt.Errorf("invalid file name")
	}

	file, err := os.Create(f)
	if err != nil {
		return err
	}

	p.file = file
	return nil
}

func (p *iMEMProfiler) Stop() {
	if p.file == nil {
		return
	}

	err := pprof.WriteHeapProfile(p.file)
	if err != nil {
		fmt.Printf("memory profile failed:%s", err.Error())
	}

	p.file.Close()
	p.file = nil
}

func NewCPUProfiler() Profiler {
	return &iCPUProfiler{}
}

func NewMEMProfiler() Profiler {
	return &iMEMProfiler{}
}

/* vim: set ts=4 sw=4 sts=4 tw=100 noet: */
