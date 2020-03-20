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

	"github.com/Mantsje/iterum-cli/idv/ctl"
	"github.com/Mantsje/iterum-cli/util"
)

// DaemonURL is the url at which we can reach the idv/iterum daemon
const DaemonURL = "http://localhost:3000/"

var errConflictingDataset = errors.New("Error: POST dataset failed, dataset already exists")

// _get takes a url to fire a get request upon and a pointer to an interface to store the result in
// It returns an error on failure of either http.Get, Reading response or Unmarshalling json body
func _get(url string, target interface{}) (err error) {
	defer _returnErrOnPanic(&err)()

	resp, err := http.Get(url)
	util.PanicIfErr(err, "")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	util.PanicIfErr(err, "")

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusNotFound:
		return fmt.Errorf("Error: GET failed with code %v, resource not found. Is the provided dataset correct?", resp.StatusCode)
	default:
		return fmt.Errorf("Error: GET failed, daemon responded with statuscode %v", resp.StatusCode)
	}

	err = json.Unmarshal([]byte(body), target)
	util.PanicIfErr(err, "")

	return
}

// constructMultiFileRequest creates a new file upload http request with optional extra otherParams
func constructMultiFileRequest(url string, otherParams map[string]string, nameFileMap map[string]string) (request *http.Request, err error) {
	defer _returnErrOnPanic(&err)()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	for filename, path := range nameFileMap {
		file, err := os.Open(path)
		util.PanicIfErr(err, "")
		defer file.Close()

		part, err := writer.CreateFormFile(filename, filepath.Base(path))
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
	defer _returnErrOnPanic(&err)()
	request, err := constructMultiFileRequest(url, nil, filemap)
	util.PanicIfErr(err, "")

	client := &http.Client{}
	resp, err := client.Do(request)

	return resp, err
}

// getBranch pulls a specific branch based on its hash
func getBranch(bhash hash, dataset string) (branch Branch, err error) {
	err = _get(DaemonURL+dataset+"/branch/"+bhash.String(), &branch)
	return
}

// getCommit pulls a specific commmit based on its hash
func getCommit(chash hash, dataset string) (commit Commit, err error) {
	err = _get(DaemonURL+dataset+"/commit/"+chash.String(), &commit)
	return
}

// getVTree pulls the entire version history file: vtree for the given dataset
func getVTree(dataset string) (history VTree, err error) {
	err = _get(DaemonURL+dataset+"/vtree", &history)
	return
}

// postDataset posts the passed dataset to the Daemon.
func postDataset(ctl ctl.DataCTL) (err error) {
	data, err := json.Marshal(ctl)
	if err != nil {
		return
	}

	resp, err := http.Post(DaemonURL, "application/json", bytes.NewBuffer(data))
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

// pushCommit pushes a commit to a branch. returns the updated VTree and Branch
func pushCommit(commit Commit, stagemap Stagemap) (branch Branch, history VTree, err error) {
	filemap := make(map[string]string)
	for key, val := range stagemap {
		filemap[key] = val
	}
	filemap["commit"] = commit.ToFilePath(false)

	response, err := _postMultipartForm(DaemonURL+"/push/commit", filemap)

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Header)

	fmt.Println(string(body))
	return
}

// pushBranchedCommit pushes a commit which is the root of a new branch. returns the updated VTree
func pushBranchedCommit(branch Branch, commit Commit, stagemap Stagemap) (history VTree, err error) {
	filemap := make(map[string]string)
	for key, val := range stagemap {
		filemap[key] = val
	}
	filemap["commit"] = commit.ToFilePath(false)
	filemap["branch"] = branch.ToFilePath(false)

	response, err := _postMultipartForm(DaemonURL+"/push/branched-commit", filemap)

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Header)

	fmt.Println(string(body))
	return
}
