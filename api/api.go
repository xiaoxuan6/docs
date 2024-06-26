package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/levigross/grequests"
    "github.com/tidwall/gjson"
    "net/http"
    "os"
    "strings"
)

type Response struct {
    Status int         `json:"status"`
    Msg    string      `json:"msg"`
    Data   interface{} `json:"data"`
}

func Api(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    response, err := grequests.Get(fmt.Sprintf("https://api.github.com/gists/%s", os.Getenv("gist_id")), nil)
    if err != nil {
        body, _ := json.Marshal(&Response{
            500,
            err.Error(),
            nil,
        })
        _, _ = w.Write(body)
        return
    }

    item := response.String()
    files := gjson.Get(item, "files").String()

    if len(files) < 1 {
        body, _ := json.Marshal(&Response{
            500,
            "暂无数据",
            nil,
        })
        _, _ = w.Write(body)
        return
    }

    body := map[string]map[string]string{
        "docs.json": map[string]string{},
    }
    _ = json.NewDecoder(bytes.NewBufferString(files)).Decode(&body)
    content := body["docs.json"]["content"]
    content = strings.ReplaceAll(content, `\n`, "")
    content = strings.ReplaceAll(content, `\"`, `"`)

    var bods map[string]map[string]string
    _ = json.Unmarshal([]byte(content), &bods)

    b, _ := json.Marshal(&Response{
        200,
        "ok",
        bods,
    })
    _, _ = w.Write(b)
}
