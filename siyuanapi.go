package api

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"path"
	"strings"
)

type SiYuanApi struct {
	token  string
	url    string
	client *req.Client
}

func NewSiYuanApi(token, url string) (*SiYuanApi, error) {
	a := &SiYuanApi{
		token:  token,
		url:    strings.TrimSuffix(url, "/"),
		client: req.C(),
	}
	a.client.SetBaseURL(a.url)
	a.client.SetCommonHeader("Authorization", fmt.Sprintf("token %s", a.token))
	_, err := a.SystemVersion()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func NewSiYuanApiWithDefaultUrl(token string) (*SiYuanApi, error) {
	return NewSiYuanApi(token, "http://127.0.0.1:6806")
}

func (r SiYuanApi) jsonReq(endpoint string, data interface{}) (*req.Response, error) {
	return r.client.R().SetBodyJsonMarshal(data).Post(endpoint)
}

// 笔记本

// NotebookLsNotebooks 列出笔记本
func (r SiYuanApi) NotebookLsNotebooks() ([]NotebookLsInfo, error) {
	resp, err := r.jsonReq("/api/notebook/lsNotebooks", nil)
	if err != nil {
		return nil, err
	}
	ret := struct {
		Response
		Data struct {
			Notebooks []NotebookLsInfo
		}
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret.Data.Notebooks, nil
}

func (r SiYuanApi) openOrCloseNotebook(op string, notebook string) error {
	_, err := r.jsonReq(fmt.Sprintf("/api/notebook/%sNotebook", op), map[string]string{
		"notebook": notebook,
	})
	if err != nil {
		return err
	}
	return nil
}

// NotebookOpenNotebook 打开笔记本
func (r SiYuanApi) NotebookOpenNotebook(notebook string) error {
	return r.openOrCloseNotebook("open", notebook)
}

// NotebookCloseNotebook 关闭笔记本
func (r SiYuanApi) NotebookCloseNotebook(notebook string) error {
	return r.openOrCloseNotebook("close", notebook)
}

// NotebookRenameNotebook 重命名笔记本
func (r SiYuanApi) NotebookRenameNotebook(notebook, newName string) error {
	_, err := r.jsonReq("/api/notebook/renameNotebook", map[string]string{
		"notebook": notebook,
		"name":     newName,
	})
	if err != nil {
		return err
	}
	return nil
}

// NotebookCreateNotebook 创建笔记本
func (r SiYuanApi) NotebookCreateNotebook(name string) (*NotebookLsInfo, error) {
	resp, err := r.jsonReq("/api/notebook/createNotebook", map[string]string{
		"name": name,
	})
	if err != nil {
		return nil, err
	}
	ret := struct {
		Response
		Data struct {
			Notebook NotebookLsInfo `json:"notebook"`
		}
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return &ret.Data.Notebook, nil
}

// NotebookRemoveNotebook 删除笔记本
func (r SiYuanApi) NotebookRemoveNotebook(notebook string) error {
	_, err := r.jsonReq("/api/notebook/removeNotebook", map[string]string{
		"notebook": notebook,
	})
	if err != nil {
		return err
	}
	return nil
}

// 获取笔记本配置
// 保存笔记本配置

// 文档

// FiletreeCreateDocWithMd 通过markdown创建文档
func (r SiYuanApi) FiletreeCreateDocWithMd(notebook, path, markdown string) (string, error) {
	resp, err := r.jsonReq("/api/filetree/createDocWithMd", map[string]string{
		"notebook": notebook,
		"path":     path,
		"markdown": markdown,
	})
	if err != nil {
		return "", err
	}
	ret := struct {
		Response
		Data string
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return "", err
	}
	return ret.Data, nil
}

// 重命名文档
// 删除文档
// 移动文档
// 根据路径获取人类可读路径
// 根据ID获取人类可读路径

// 资源文件

// AssetUpload 上传资源文件
func (r SiYuanApi) AssetUpload(filepath string) (string, error) {
	resp, err := r.client.R().SetFormData(map[string]string{
		"assetsDirPath": "/assets/",
	}).SetFile("file[]", filepath).Post("/api/asset/upload")
	if err != nil {
		return "", err
	}
	res := struct {
		Response
		Data struct {
			ErrFiles []string          `json:"errFiles"`
			SuccMap  map[string]string `json:"succMap"`
		} `json:"data"`
	}{}
	err = resp.Unmarshal(&res)
	if err != nil {
		return "", err
	}
	if len(res.Data.ErrFiles) > 0 {
		return "", errors.New("upload failed")
	}
	return res.Data.SuccMap[path.Base(filepath)], nil
}

// AssetUploadBytes 上传资源文件（数据来自内存）
func (r SiYuanApi) AssetUploadBytes(data []byte, filename string) (string, error) {
	resp, err := r.client.R().SetFormData(map[string]string{
		"assetsDirPath": "/assets/",
	}).SetFileBytes("file[]", filename, data).Post("/api/asset/upload")
	if err != nil {
		return "", err
	}
	res := struct {
		Response
		Data struct {
			ErrFiles []string          `json:"errFiles"`
			SuccMap  map[string]string `json:"succMap"`
		} `json:"data"`
	}{}
	err = resp.Unmarshal(&res)
	if err != nil {
		return "", err
	}
	if len(res.Data.ErrFiles) > 0 {
		return "", errors.New("upload failed")
	}
	return res.Data.SuccMap[filename], nil
}

// 块

// BlockInsertBlock 插入块
func (r SiYuanApi) BlockInsertBlock(dataType, data, previousId string) ([]OperationInfo, error) {
	resp, err := r.jsonReq("/api/block/insertBlock", map[string]string{
		"dataType":   dataType,
		"data":       data,
		"previousID": previousId,
	})
	if err != nil {
		return nil, err
	}
	ret := responseBlockOperations{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret.Data[0].DoOperations, nil
}

// BlockPrependBlock 插入前置子块
func (r SiYuanApi) BlockPrependBlock(dataType, data, parentId string) ([]OperationInfo, error) {
	resp, err := r.jsonReq("/api/block/prependBlock", map[string]string{
		"data":     data,
		"dataType": dataType,
		"parentID": parentId,
	})
	if err != nil {
		return nil, err
	}
	ret := responseBlockOperations{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret.Data[0].DoOperations, nil
}

// BlockAppendBlock 插入后置子块
func (r SiYuanApi) BlockAppendBlock(dataType, data, parentId string) ([]OperationInfo, error) {
	resp, err := r.jsonReq("/api/block/appendBlock", map[string]string{
		"data":     data,
		"dataType": dataType,
		"parentID": parentId,
	})
	if err != nil {
		return nil, err
	}
	ret := responseBlockOperations{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, err
	}
	return ret.Data[0].DoOperations, nil
}

// 更新块
// 删除块

// BlockGetBlockKramdown 获取块kramdown源码
func (r SiYuanApi) BlockGetBlockKramdown(id string) (string, error) {
	resp, err := r.jsonReq("/api/block/getBlockKramdown", map[string]string{
		"id": id,
	})
	if err != nil {
		return "", err
	}
	ret := struct {
		id       string
		kramdown string
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return "", err
	}
	return ret.kramdown, nil
}

// 属性

// 设置块属性
// 获取块属性

// SQL

// 执行SQL查询

// 模板

// 渲染模板

// 文件

// FileGetFile 获取文件
func (r SiYuanApi) FileGetFile(path string) (string, error) {
	resp, err := r.jsonReq("/api/file/getFile", map[string]string{
		"path": path,
	})
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

// 写入文件

// 导出

// 导出markdown文本

// 通知

// NotificationPushMsg 推送消息
func (r SiYuanApi) NotificationPushMsg(message string) (string, error) {
	resp, err := r.jsonReq("/api/notification/pushMsg", struct {
		Msg     string `json:"msg"`
		Timeout uint64 `json:"timeout"`
	}{
		Msg:     message,
		Timeout: 5000,
	})
	if err != nil {
		return "", err
	}
	res := struct {
		Response
		Data struct {
			Id string `json:"id"`
		} `json:"data"`
	}{}
	err = resp.Unmarshal(&res)
	if err != nil {
		return "", err
	}
	return res.Data.Id, nil
}

// 推送报错消息

// 系统

// 获取启动进度

// SystemVersion 获取系统版本
func (r SiYuanApi) SystemVersion() (string, error) {
	resp, err := r.jsonReq("/api/system/version", nil)
	if err != nil {
		return "", err
	}
	info := struct {
		Response
		Data string `json:"data"`
	}{}
	err = resp.UnmarshalJson(&info)
	if err != nil {
		return "", err
	}
	// fmt.Print(info.Data)
	return info.Data, nil
}

// 获取系统当前时间
