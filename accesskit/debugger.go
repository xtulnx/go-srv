package accesskit

import (
	"gorm.io/gen/helper"
	"gorm.io/gorm/utils"
	"path"
	"runtime"
	_ "unsafe"
)

// gorm 调试辅助

//go:linkname gormSourceDir gorm.io/gorm/utils.gormSourceDir
var gormSourceDir string

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

type _fGenObject struct {
}

func (d *_fGenObject) TableName() string { return "" }
func (d *_fGenObject) StructName() string {
	_, file, _, _ := runtime.Caller(1)
	// fmt.Println(file)

	//// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = path.Dir(path.Dir(path.Dir(file)))
	return ""
}
func (d *_fGenObject) FileName() string         { return "" }
func (d *_fGenObject) ImportPkgPaths() []string { return nil }
func (d *_fGenObject) Fields() []helper.Field   { return nil }

// 通过辅助类查询
func resetGormLineNum() {
	//var c gen.Config
	//fmt.Println("targetDir", reflect.TypeOf(c).PkgPath())
	_ = helper.CheckObject(&_fGenObject{})
	//fmt.Println("targetDir", gormSourceDir)
	//logkit.ForTask("debug").Info("修改:")
}

// -o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-

// DebugGorm gorm 的调试
//   - 在用 gorm.io/gen 的情况下，校正调试日志的正确路径
func DebugGorm() {
	_ = utils.FileWithLineNum()
	//fmt.Println("sourceDir", gormSourceDir)
	if gormSourceDir != "" {
		gormSourceDir = path.Dir(path.Dir(gormSourceDir))
		//fmt.Println("targetDir", gormSourceDir)
	}
}
