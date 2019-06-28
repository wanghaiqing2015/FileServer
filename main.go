package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
    "time"
)

var (
    Info *log.Logger
    Warning *log.Logger
    Error * log.Logger
)

func init(){
    //log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
    //log.Println("飞雪无情的博客:","http://www.flysnow.org")
    //log.Printf("飞雪无情的微信公众号：%s\n","flysnow_org")

    errFile,err:=os.OpenFile("errors.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
    if err!=nil{
        log.Fatalln("打开日志文件失败：",err)
    }

    Info = log.New(io.MultiWriter(os.Stderr,errFile),"Info：",log.Ldate | log.Ltime | log.Lshortfile)
    Warning = log.New(io.MultiWriter(os.Stderr,errFile),"Warning：",log.Ldate | log.Ltime | log.Lshortfile)
    Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error：",log.Ldate | log.Ltime | log.Lshortfile)
}

func GetDestFilePath()string{
    var Path string
    if runtime.GOOS == "windows"{
        Path = "C:\\data"
    }else{
        Path = "/opt/data"
    }

    err := Exists(Path)
    if err !=nil{
        os.MkdirAll(Path, os.ModePerm)
    }
    return Path
}

func Exists(path string) error {
    _, err := os.Stat(path)    //os.Stat获取文件信息
    if err != nil {
        if os.IsExist(err) {
            return nil
        }
        return fmt.Errorf("文件不存在")
    }
    return fmt.Errorf("文件不存在")
}

type Test struct {
    Url string `json:"url,omitempty"`
}

type JsonContent struct {
    Status bool `json:"status,omitempty"`
    Msg string `json:"msg,omitempty"`
    Data interface{} `json:"data,omitempty"`
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(32 << 20)  //解析参数，默认是不会解析的

    DataPath := GetDestFilePath()
    Date := time.Now().Format("2006-01-02")
    Path := filepath.Join(DataPath, Date)
    err := Exists(Path)
    if err !=nil{
        os.MkdirAll(Path, os.ModePerm)
    }

    var UrlList []string
    for _,v := range r.MultipartForm.File {
        for i := 0; i < len(v); i++ {
            FilePath := v[i].Filename
            FileName :=filepath.Base(FilePath)

            formFile, _ := v[i].Open()
            defer formFile.Close()

            destFile, err := os.Create( filepath.Join(Path, FileName))
            if err != nil {
                fmt.Println("创建文件失败")
            }
            defer destFile.Close()

            _, err = io.Copy(destFile, formFile)
            if err != nil {
                log.Println("写文件报错: ", err)
                return
            }

            FileUrl := fmt.Sprintf("/%s/%s/%s" , filepath.Base(DataPath), Date, FileName)
            UrlList = append(UrlList, FileUrl)
        }
    }

    var result JsonContent
    result.Status = true
    result.Msg = "文件下载成功！"
    result.Data = UrlList
    ret, _ := json.Marshal(result)

    fmt.Println( string(ret) )

    fmt.Fprintf(w, string(ret) ) //这个写入到w的是输出到客户端的
}

func main() {
    http.Handle("/data/",
        http.StripPrefix("/data/",
            http.FileServer(http.Dir(GetDestFilePath())),
        ),
    )
    http.HandleFunc("/upload/", UploadFile) //设置访问的路由
    ListenAddress := ":9090"
    Info.Println("服务已经运行，监听端口为", ListenAddress)

    err := http.ListenAndServe(ListenAddress, nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

//文件太大的话，是不是会耗内存
//日志分文件，每天生成一个
//做一个界面，非常简单。是一个文件服务器