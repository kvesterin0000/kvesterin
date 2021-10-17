package locales

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Locales struct {
	bundle *i18n.Bundle
}

func New(paths ...string) (*Locales, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, p := range paths {
		_, err := bundle.LoadMessageFile(p)
		if err != nil {
			return nil, fmt.Errorf("cannot load message file %w", err)
		}
	}
	return &Locales{bundle: bundle}, nil
}

func (l *Locales) TranslatePage(lang string, msgIDs ...string) (map[string]string, error) {
	var msgs = make(map[string]string) //map[messageID]messageText

	localizer := i18n.NewLocalizer(l.bundle, lang)
	for _, msgID := range msgIDs {
		lc, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: msgID,
		})
		if err != nil {
			return msgs, fmt.Errorf("cannot localize %w", err)
		}
		msgs[msgID] = lc
	}
	return msgs, nil
}
