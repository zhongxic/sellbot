package bot

import (
	"github.com/zhongxic/sellbot/internal/service/bot/session"
	"github.com/zhongxic/sellbot/internal/service/process"
)

func analyzeIntention(session *session.Session, loadedProcess *process.Process) []process.IntentionRule {
	// TODO analyze intention
	return make([]process.IntentionRule, 0)
}
