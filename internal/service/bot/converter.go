package bot

import (
	"strings"

	"github.com/zhongxic/sellbot/internal/service/bot/matcher"
	"github.com/zhongxic/sellbot/internal/service/process"
)

func makeRespond(matchContext *matcher.Context, answerDTO AnswerDTO, intentionRules []process.IntentionRule) *InteractiveRespond {
	interactiveRespond := &InteractiveRespond{}
	interactiveRespond.SessionId = matchContext.Session.Id
	interactiveRespond.Hits.Sentence = matchContext.Sentence
	if len(matchContext.Segments) == 0 {
		interactiveRespond.Hits.Segments = make([]string, 0)
	} else {
		interactiveRespond.Hits.Segments = matchContext.Segments
	}
	interactiveRespond.Hits.HitPaths = convertMatchedPathListToHitPathDTOList(matchContext.MatchedPaths)
	interactiveRespond.Answer.Text = replaceVariables(answerDTO.Text, matchContext.Session.Variables)
	interactiveRespond.Answer.Audio = answerDTO.Audio
	interactiveRespond.Answer.Ended = answerDTO.Ended
	interactiveRespond.Answer.Agent = answerDTO.Agent
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

func replaceVariables(text string, variables map[string]string) string {
	if len(variables) == 0 {
		return text
	}
	replacements := make([]string, 0)
	for code, value := range variables {
		replacements = append(replacements, code, value)
	}
	return strings.NewReplacer(replacements...).Replace(text)
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
