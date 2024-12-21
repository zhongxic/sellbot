package bot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/zhongxic/sellbot/internal/service/process"
)

func (s *serviceImpl) Prologue(ctx context.Context, prologueDTO *PrologueDTO) (*InteractiveRespond, error) {
	slog.Info("start process prologue", "traceId", ctx.Value("traceId"))
	loadedProcess, err := s.processManager.Load(prologueDTO.ProcessId, prologueDTO.Test)
	if err != nil {
		return nil, err
	}
	if err := loadedProcess.Validate(); err != nil {
		return nil, err
	}
	if err := validateVariables(prologueDTO.Variables, loadedProcess.Variables); err != nil {
		return nil, err
	}
	return &InteractiveRespond{}, nil
}

func validateVariables(actual map[string]string, expected []process.Variable) error {
	params := actual
	if actual == nil {
		params = make(map[string]string)
	}
	variables := expected
	if expected == nil {
		variables = make([]process.Variable, 0)
	}
	if len(params) != len(variables) {
		return fmt.Errorf("process variables not matched expected [%d] actual [%d]", len(variables), len(params))
	}
	messages := make([]string, 0)
	for _, variable := range variables {
		if _, ok := params[variable.Code]; !ok {
			message := fmt.Sprintf("process variable [%s] is required", variable.Code)
			messages = append(messages, message)
		}
	}
	if len(messages) > 0 {
		return errors.New(strings.Join(messages, ", "))
	}
	return nil
}
