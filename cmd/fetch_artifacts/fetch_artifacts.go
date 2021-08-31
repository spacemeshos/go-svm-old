package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	branch string
	token  string
	dest   string
)

func init() {
	flag.StringVar(&branch, "branch", "", "svm branch to download the artifacts from it's latest workflow run")
	flag.StringVar(&token, "token", "", "github token for better api privileges")
	flag.StringVar(&dest, "dest", "", "destination folder")
	flag.Parse()
}

type archivesDownloadURL struct {
	linux     string
	macOS     string
	windows   string
	wasmCodec string
}

func main() {
	log.Printf("%v; branch: %v, token: %v, dest: %v\n", filepath.Base(os.Args[0]), branch, token, dest)

	runURL, err := runURL()
	noError(err)

	artifactsURL, err := artifactsURL(runURL)
	noError(err)

	archivesURL, err := archivesURL(artifactsURL)
	noError(err)

	noError(fetch(archivesURL.linux, dest))
	noError(fetch(archivesURL.macOS, dest))
	noError(fetch(archivesURL.windows, dest))
	noError(fetch(archivesURL.wasmCodec, dest))
}

func runURL() (string, error) {
	url := "https://api.github.com/repos/spacemeshos/svm/actions/runs"
	if branch != "" {
		url += fmt.Sprintf("?branch=%v", branch)
	}
	log.Printf("GET %v", url)
	res, err := http.DefaultClient.Do(req("GET", url))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return "", err
	}

	totalCount, ok := body["total_count"].(float64)
	if !ok {
		return "", fmt.Errorf("unexpected response: %s", rawBody)
	}
	if totalCount == 0 {
		return "", fmt.Errorf("no workflow runs found for branch %v", branch)
	}

	workflowRuns := body["workflow_runs"].([]interface{})
	workflowRun := workflowRuns[0].(map[string]interface{})

	return workflowRun["url"].(string), nil
}

func artifactsURL(runURL string) (string, error) {
	log.Printf("GET %v", runURL)
	res, err := http.DefaultClient.Do(req("GET", runURL))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return "", err
	}

	url, ok := body["artifacts_url"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response: %s", rawBody)
	}

	return url, nil
}

func archivesURL(artifactsURL string) (archivesDownloadURL, error) {
	log.Printf("GET %v", artifactsURL)
	res, err := http.DefaultClient.Do(req("GET", artifactsURL))
	if err != nil {
		return archivesDownloadURL{}, err
	}

	defer res.Body.Close()
	rawBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return archivesDownloadURL{}, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return archivesDownloadURL{}, err
	}

	//totalCount, ok := body["total_count"].(float64)
	//if !ok {
	//	return archivesDownloadURL{}, fmt.Errorf("unexpected response: %s", rawBody)
	//}
	//if totalCount != 4 {
	//	return archivesDownloadURL{}, fmt.Errorf("found artifacts listing for %v platforms, expected 4", totalCount)
	//}

	ret := archivesDownloadURL{}
	items := body["artifacts"].([]interface{})
	for _, item := range items {
		artifact := item.(map[string]interface{})
		url := artifact["archive_download_url"].(string)
		switch artifact["name"] {
		case "bins-Linux-release":
			ret.linux = url
		case "bins-macOS-release":
			ret.macOS = url
		case "bins-Windows-release":
			ret.windows = url
		//case "svm_codec.wasm":
		//	ret.wasmCodec = url
		default:
			return archivesDownloadURL{}, fmt.Errorf("invalid artifact tag: %v", artifact["name"])
		}
	}

	return ret, nil
}

func fetch(url, name string) error {
	name, err := download(url)
	if err != nil {
		return err
	}

	if err := unzip(name); err != nil {
		return err
	}

	if err := os.Remove(name); err != nil {
		return err
	}

	return nil
}

func download(url string) (string, error) {
	log.Printf("GET %v", url)
	res, err := http.DefaultClient.Do(req("GET", url))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	name := filepath.Join(dest, "temp.zip")
	file, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	written, err := io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	log.Printf("downloaded %v bytes to %v\n", written, name)

	return name, nil
}

func unzip(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		name := filepath.Join(dest, f.Name)
		file, err := os.Create(name)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		written, err := io.Copy(file, rc)

		file.Close()
		rc.Close()

		if err != nil {
			return err
		}

		log.Printf("unzipped %v bytes to %v\n", written, name)
	}

	return nil
}

func req(method, url string) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	//req.Header.Set("Authorization", fmt.Sprintf("token %v", token))
	return req
}

func noError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
