package idv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mantsje/iterum-cli/util"
)

// DaemonURL is the url at which we can reach the idv/iterum daemon
const DaemonURL = "http://idv-daemon.com/"

// _get takes a url to fire a get request upon and a pointer to an interface to store the result in
// It returns an error on failure of either http.Get, Reading response or Unmarshalling json body
func _get(url string, target interface{}) (err error) {
	defer _returnErrOnPanic(&err)()

	resp, err := http.Get(url)
	util.PanicIfErr(err, "")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	util.PanicIfErr(err, "")

	err = json.Unmarshal([]byte(body), target)
	util.PanicIfErr(err, "")

	return
}

// constructMultiFileRequest creates a new file upload http request with optional extra otherParams
func constructMultiFileRequest(uri string, otherParams map[string]string, nameFileMap map[string]string) (request *http.Request, err error) {
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

	request, err = http.NewRequest("POST", uri, body)
	util.PanicIfErr(err, "")
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return
}

func _postMultipartForm(url string, filemap map[string]string) (response *http.Response, err error) {
	defer _returnErrOnPanic(&err)()
	request, err := constructMultiFileRequest(DaemonURL+"/push/commit", nil, filemap)
	util.PanicIfErr(err, "")

	client := &http.Client{}
	resp, err := client.Do(request)

	return resp, err
}

// pullBranch pulls a specific branch based on its hash
func pullBranch(bhash hash) (branch Branch, err error) {
	err = _get(DaemonURL+"/pull/branch/"+bhash.String(), &branch)
	fmt.Println(branch)
	return
}

// pullCommit pulls a specific commmit based on its hash
func pullCommit(chash hash) (commit Commit, err error) {
	err = _get(DaemonURL+"/pull/commit/"+chash.String(), &commit)
	fmt.Println(commit)
	return
}

// pullVTree pulls the entire version history file: vtree
func pullVTree() (history VTree, err error) {
	err = _get(DaemonURL+"/pull/vtree", &history)
	fmt.Println(history)
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
