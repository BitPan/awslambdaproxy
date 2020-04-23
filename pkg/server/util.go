package server

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

const (
	getIPUrl = "http://checkip.amazonaws.com/"
)

func bidirectionalCopy(src io.ReadWriteCloser, dst io.ReadWriteCloser) {
	defer dst.Close()
	defer src.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(dst, src)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		io.Copy(src, dst)
		wg.Done()
	}()
	wg.Wait()
}

func getPublicIP() (string, error) {
	resp, err := http.Get(getIPUrl)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get IP address from "+getIPUrl)
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read IP address from "+getIPUrl)
	}
	return string(bytes.TrimSpace(buf)), nil
}

func getSessionAWS() (*session.Session, error) {
	sess, err := session.NewSession(aws.NewConfig())
	if err != nil {
		return nil, err
	}
	if _, err = sess.Config.Credentials.Get(); err != nil {
		return nil, err
	}
	return sess, nil
}
