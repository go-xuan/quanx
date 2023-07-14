package filex

import (
	"os"
	"strings"
	"sync"
)

type FileList []*File
type File struct {
	Path string
	Info os.FileInfo
}

type FileType uint

const (
	DirAndFile FileType = iota
	OnlyDir
	OnlyFile
	GoType
	ModType
	JavaType
	JarType
	ClassType
	SqlType
	LogType
	ShellType
	BatType
	TxtType
	JsonType
	YmlType
	YamlType
	XmlType
	HtmlType
	PropertiesType
	DocType
	DocxType
	XlsType
	XlsxType
	PptType
	PptxType
	PdfType
	MdType
	Mp3Type
	Mp4Type
	JpgType
	HeicType
	WavType
)

var FileTypeMap sync.Map

func init() {
	FileTypeMap.Store(GoType, ".go")
	FileTypeMap.Store(ModType, ".mod")
	FileTypeMap.Store(JavaType, ".java")
	FileTypeMap.Store(JarType, ".jar")
	FileTypeMap.Store(ClassType, ".class")
	FileTypeMap.Store(SqlType, ".sql")
	FileTypeMap.Store(LogType, ".log")
	FileTypeMap.Store(ShellType, ".sh")
	FileTypeMap.Store(BatType, ".bat")
	FileTypeMap.Store(TxtType, ".txt")
	FileTypeMap.Store(JsonType, ".json")
	FileTypeMap.Store(YmlType, ".yml")
	FileTypeMap.Store(YamlType, ".yaml")
	FileTypeMap.Store(XmlType, ".xml")
	FileTypeMap.Store(HtmlType, ".html")
	FileTypeMap.Store(PropertiesType, ".properties")
	FileTypeMap.Store(DocType, ".doc")
	FileTypeMap.Store(DocxType, ".docx")
	FileTypeMap.Store(XlsType, ".xls")
	FileTypeMap.Store(XlsxType, ".xlsx")
	FileTypeMap.Store(PptType, ".ppt")
	FileTypeMap.Store(PptxType, ".pptx")
	FileTypeMap.Store(PdfType, ".pdf")
	FileTypeMap.Store(MdType, ".md")
	FileTypeMap.Store(Mp3Type, ".mp3")
	FileTypeMap.Store(Mp4Type, ".mp4")
	FileTypeMap.Store(JpgType, ".jpg")
	FileTypeMap.Store(HeicType, ".heic")
	FileTypeMap.Store(WavType, ".wav")
}

// 判断文件类型是否匹配
func FileTypeMatch(fileName string, mode FileType) bool {
	var match = false
	FileTypeMap.Range(func(k, v any) bool {
		if mode == k && strings.HasSuffix(fileName, v.(string)) {
			match = true
			return false
		}
		return true
	})
	return match
}
