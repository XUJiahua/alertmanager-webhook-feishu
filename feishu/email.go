package feishu

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/xujiahua/alertmanager-webhook-feishu/config"
	"sync"
)

type EmailHelper struct {
	sdk *Sdk
	sync.RWMutex
	// email -> open_id mapping
	cache map[string]string
}

func NewEmailHelper(app *config.App) (*EmailHelper, error) {
	if app.ID == "" || app.Secret == "" {
		return nil, errors.New("appID, appSecret required")
	}
	return &EmailHelper{
		sdk:   NewSDK(app.ID, app.Secret),
		cache: make(map[string]string),
	}, nil
}

// Lookup open_ids by emails
func (o *EmailHelper) Lookup(emails []string) ([]string, error) {
	// open_id
	var res []string
	var unknownEmails []string
	o.RLock()
	for _, email := range emails {
		if openID, ok := o.cache[email]; ok {
			res = append(res, openID)
		} else {
			unknownEmails = append(unknownEmails, email)
		}
	}
	o.RUnlock()

	if len(unknownEmails) != 0 {
		remaining, err := o.batchGetID(unknownEmails)
		if err != nil {
			return nil, err
		}
		res = append(res, remaining...)
	}

	return res, nil
}

func (o *EmailHelper) batchGetID(emails []string) ([]string, error) {
	logrus.Debugf("lookup emails %v from feishu", emails)
	mapping, err := o.sdk.BatchGetID(emails)
	if err != nil {
		return nil, err
	}

	o.Lock()
	defer o.Unlock()
	var openIDs []string
	for email, openID := range mapping {
		o.cache[email] = openID
		openIDs = append(openIDs, openID)
	}

	return openIDs, nil
}
