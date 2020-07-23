package idv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/cli/util"
)

var (
	errConflictingDataset = errors.New("Error: POST dataset failed, dataset already exists")
	errConflictingCommit  = errors.New("Error: POST commit failed, commit is not child of HEAD. Pull latest changes to resolve")
	errNotFound           = errors.New("Error: Daemon responded with 404, resource not found")
)

// _get takes a url to fire a get request upon and a pointer to an interface to store the result in
// It returns an error on failure of either http.Get, Reading response or Unmarshalling json body
func _get(url string, target interface{}) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	resp, err := http.Get(url)
	util.PanicIfErr(err, "")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	util.PanicIfErr(err, "")

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return errNotFound
	default:
		return fmt.Errorf("Error: GET failed, daemon responded with statuscode %v", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(body), target)
	util.PanicIfErr(err, "")

	return
}

// constructMultiFileRequest creates a new file upload http request with optional extra otherParams
func constructMultiFileRequest(url string, otherParams map[string]string, nameFileMap map[string]string) (request *http.Request, err error) {
	defer util.ReturnErrOnPanic(&err)()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for filename, path := range nameFileMap {
		file, err := os.Open(path)
		util.PanicIfErr(err, "")
		defer file.Close()

		part, err := writer.CreateFormFile(filepath.Base(path), filename)
		util.PanicIfErr(err, "")
		io.Copy(part, file)
	}

	for key, val := range otherParams {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	util.PanicIfErr(err, "")

	request, err = http.NewRequest("POST", url, body)
	util.PanicIfErr(err, "")
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return
}

func _postMultipartForm(url string, filemap map[string]string) (response *http.Response, err error) {
	defer util.ReturnErrOnPanic(&err)()
	request, err := constructMultiFileRequest(url, nil, filemap)
	util.PanicIfErr(err, "")

	client := &http.Client{}
	return client.Do(request)
}

// getBranch pulls a specific branch based on its hash from the daemon
func getBranch(ctl ctl.DataCTL, bhash hash) (branch Branch, err error) {
	err = _get(ctl.DaemonURL+ctl.Name+"/branch/"+bhash.String(), &branch)
	return
}

// getCommit pulls a specific commmit based on its hash
func getCommit(ctl ctl.DataCTL, chash hash) (commit Commit, err error) {
	err = _get(ctl.DaemonURL+ctl.Name+"/commit/"+chash.String(), &commit)
	return
}

// getConfig pulls a config based on the passed dataset name
func getConfig(ctl ctl.DataCTL) (ctlremote ctl.DataCTL, err error) {
	var m map[string]interface{}
	err = _get(ctl.DaemonURL+ctl.Name, &m)
	if err != nil {
		return
	}
	err = ctlremote.ParseFromMap(m)
	return
}

// getVTree pulls the entire version history file: vtree for the given dataset
func getVTree(ctl ctl.DataCTL) (history VTree, err error) {
	err = _get(ctl.DaemonURL+ctl.Name+"/vtree", &history)
	return
}

// getDatasets pulls the list of datasets currently known to the daemon
func getDatasets(ctl ctl.DataCTL) (datasets []string, err error) {
	err = _get(ctl.DaemonURL, &datasets)
	return
}

// postDataset posts the passed dataset to the Daemon.
func postDataset(ctl ctl.DataCTL) (err error) {
	data, err := json.Marshal(ctl)
	if err != nil {
		return
	}

	resp, err := http.Post(ctl.DaemonURL, "application/json", bytes.NewBuffer(data))
	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusConflict:
		return errConflictingDataset
	default:
		return fmt.Errorf("Error: POST failed, daemon responded with statuscode %v", resp.StatusCode)
	}
	return
}

// _performCommit actually sends the commit to the Daemon and parses its response
func _performCommit(url string, filemap map[string]string) (branch Branch, history VTree, err error) {
	defer util.ReturnErrOnPanic(&err)()
	response, err := _postMultipartForm(url, filemap)
	util.PanicIfErr(err, "")
	switch response.StatusCode {
	case http.StatusOK:
		break
	case http.StatusConflict:
		err = errConflictingCommit
		return
	default:
		err = fmt.Errorf("Error: POST multipart form failed, daemon responded with statuscode %v", response.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	util.PanicIfErr(err, "")
	var rawBody struct {
		VTree  VTree  `json:"vtree"`
		Branch Branch `json:"branch"`
	}
	err = json.Unmarshal(body, &rawBody)
	util.PanicIfErr(err, "")

	return rawBody.Branch, rawBody.VTree, err
}

// pushCommit pushes a commit to a branch. returns the updated VTree and Branch
func postCommit(ctl ctl.DataCTL, commit Commit, stagemap Stagemap) (branch Branch, history VTree, err error) {
	defer util.ReturnErrOnPanic(&err)()
	filemap := make(map[string]string)
	for key, val := range stagemap {
		filemap[key] = val
	}
	filemap["commit"] = tempCommitPath
	if !util.FileExists(tempCommitPath) {
		err = errors.New("Error: temporary commit could not be located, can't push")
		return
	}
	return _performCommit(ctl.DaemonURL+ctl.Name+"/commit", filemap)
}

// postBranchedCommit pushes a commit which is the root of a new branch. returns the updated VTree
func postBranchedCommit(ctl ctl.DataCTL, branch Branch, commit Commit, stagemap Stagemap) (updatedBranch Branch, history VTree, err error) {
	defer util.ReturnErrOnPanic(&err)()
	filemap := make(map[string]string)
	for key, val := range stagemap {
		filemap[key] = val
	}
	filemap["commit"] = tempCommitPath
	filemap["branch"] = branch.ToFilePath(true)
	if !util.FileExists(tempCommitPath) {
		err = errors.New("Error: temporary commit could not be located, can't push")
		return
	}

	return _performCommit(ctl.DaemonURL+ctl.Name+"/commit", filemap)
}
