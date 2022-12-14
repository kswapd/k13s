package view

import (
	"context"

	"github.com/gdamore/tcell/v2"
	"github.com/kswapd/k11s/internal/client"
	"github.com/kswapd/k11s/internal/dao"
	"github.com/kswapd/k11s/internal/ui"
)

// ServiceAccount represents a serviceaccount viewer.
type ServiceAccount struct {
	ResourceViewer
}

// NewServiceAccount returns a new viewer.
func NewServiceAccount(gvr client.GVR) ResourceViewer {
	s := ServiceAccount{
		ResourceViewer: NewBrowser(gvr),
	}
	s.AddBindKeysFn(s.bindKeys)

	return &s
}

func (s *ServiceAccount) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyU: ui.NewKeyAction("UsedBy", s.refCmd, true),
	})
}

func (s *ServiceAccount) refCmd(evt *tcell.EventKey) *tcell.EventKey {
	return scanSARefs(evt, s.App(), s.GetTable(), "v1/serviceaccounts")
}

func scanSARefs(evt *tcell.EventKey, a *App, t *Table, gvr string) *tcell.EventKey {
	path := t.GetSelectedItem()
	if path == "" {
		return evt
	}

	ctx := context.Background()
	refs, err := dao.ScanForSARefs(refContext(gvr, path, true)(ctx), a.factory)
	if err != nil {
		a.Flash().Err(err)
		return nil
	}
	if len(refs) == 0 {
		a.Flash().Warnf("No references found at this time for %s::%s. Check again later!", gvr, path)
		return nil
	}
	a.Flash().Infof("Viewing references for %s::%s", gvr, path)
	view := NewReference(client.NewGVR("references"))
	view.SetContextFn(refContext(gvr, path, false))
	if err := a.inject(view); err != nil {
		a.Flash().Err(err)
	}

	return nil
}
