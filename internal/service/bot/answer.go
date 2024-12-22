package bot

import (
	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
)

func makeAnswer(matchContext *matcher.Context) AnswerDTO {
	// TODO 实现 answer 逻辑
	return AnswerDTO{}
}

func makeRespond(matchContext *matcher.Context, answerDTO AnswerDTO, intentionRules []process.IntentionRule) *InteractiveRespond {
	interactiveRespond := &InteractiveRespond{}
	interactiveRespond.SessionId = matchContext.Session.SessionId
	hitsDTO := HitsDTO{}
	hitsDTO.Sentence = matchContext.Sentence
	if len(matchContext.Segments) == 0 {
		hitsDTO.Segments = make([]string, 0)
	} else {
		hitsDTO.Segments = matchContext.Segments
	}
	hitsDTO.HitPaths = convertMatchedPathListToHitPathDTOList(matchContext.MatchedPaths)
	interactiveRespond.Hits = hitsDTO
	interactiveRespond.Answer = answerDTO
	interactiveRespond.Intentions = convertIntentionRuleListToIntentionDTOList(intentionRules)
	return interactiveRespond
}

func convertMatchedPathListToHitPathDTOList(matchedPathList []matcher.MatchedPath) []HitPathDTO {
	if len(matchedPathList) == 0 {
		return make([]HitPathDTO, 0)
	}
	hitPathDTOList := make([]HitPathDTO, len(matchedPathList))
	for i, matchedPath := range matchedPathList {
		hitPathDTOList[i] = convertMatchedPathToHitPathDTO(matchedPath)
	}
	return hitPathDTOList
}

func convertMatchedPathToHitPathDTO(matchedPath matcher.MatchedPath) HitPathDTO {
	hitPathDTO := HitPathDTO{
		Domain: matchedPath.Domain,
		Branch: matchedPath.Branch,
	}
	if len(matchedPath.MatchedWords) == 0 {
		hitPathDTO.MatchedWords = make([]string, 0)
	} else {
		hitPathDTO.MatchedWords = matchedPath.MatchedWords
	}
	return hitPathDTO
}

func convertIntentionRuleListToIntentionDTOList(intentionRules []process.IntentionRule) []IntentionDTO {
	if len(intentionRules) == 0 {
		return make([]IntentionDTO, 0)
	}
	intentionDTOList := make([]IntentionDTO, len(intentionRules))
	for i, intentionRule := range intentionRules {
		intentionDTOList[i] = convertIntentionRuleToIntentionDTO(intentionRule)
	}
	return intentionDTOList
}

func convertIntentionRuleToIntentionDTO(intentionRule process.IntentionRule) IntentionDTO {
	return IntentionDTO{
		Code:        intentionRule.Code,
		DisplayName: intentionRule.DisplayName,
		Reason:      intentionRule.Reason,
	}
}
