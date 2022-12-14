package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/kswapd/k11s/internal/client"
	"github.com/kswapd/k11s/internal/ui"
)

// PriorityClass presents a priority class viewer.
type PriorityClass struct {
	ResourceViewer
}

// NewPriorityClass returns a new viewer.
func NewPriorityClass(gvr client.GVR) ResourceViewer {
	s := PriorityClass{
		ResourceViewer: NewBrowser(gvr),
	}
	s.AddBindKeysFn(s.bindKeys)

	return &s
}

func (s *PriorityClass) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyU: ui.NewKeyAction("UsedBy", s.refCmd, true),
	})
}

func (s *PriorityClass) refCmd(evt *tcell.EventKey) *tcell.EventKey {
	return scanRefs(evt, s.App(), s.GetTable(), "scheduling.k8s.io/v1/priorityclasses")
}
