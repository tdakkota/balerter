package manager

import (
	"fmt"
	"github.com/balerter/balerter/internal/script/script"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
	"strings"
)

func (m *Manager) getAlertName(L *lua.LState) (string, error) {
	alertName := L.Get(1).String()
	alertName = strings.TrimSpace(alertName)

	if alertName == "" || alertName == "nil" {
		return "", fmt.Errorf("alertName must be specified and not empty")
	}

	return alertName, nil
}

func (m *Manager) getArgs(L *lua.LState) (alertName string, text string, fields []string, err error) {
	alertName, err = m.getAlertName(L)
	if err != nil {
		return
	}

	text = L.Get(2).String()

	flds := L.Get(3)
	if flds.Type() == lua.LTNil {
		return
	}

	L.ToTable(3).ForEach(func(value lua.LValue, value2 lua.LValue) {
		fields = append(fields, value2.String())
	})

	return
}

func (m *Manager) on(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		alertName, text, fields, err := m.getArgs(L)
		if err != nil {
			m.logger.Error("error get args", zap.Error(err))
			L.Push(lua.LString("error: " + err.Error()))
			return 1
		}

		m.activeMx.Lock()
		defer m.activeMx.Unlock()
		info, ok := m.active[alertName]
		if !ok {
			info = &alertInfo{
				ScriptName: s.Name,
			}
			m.active[alertName] = info
			m.sendError(alertName, text, fields...)
		}
		info.Count++

		m.logger.Debug("call alert ON", zap.String("alertName", alertName), zap.Int("count", info.Count), zap.String("scriptName", s.Name))

		return 0
	}
}

func (m *Manager) off(s *script.Script) lua.LGFunction {
	return func(L *lua.LState) int {
		alertName, text, fields, err := m.getArgs(L)
		if err != nil {
			m.logger.Error("error get args", zap.Error(err))
			L.Push(lua.LString("error: " + err.Error()))
			return 1
		}

		m.activeMx.Lock()
		defer m.activeMx.Unlock()
		if _, ok := m.active[alertName]; ok {
			delete(m.active, alertName)
			m.sendSuccess(alertName, text, fields...)
		}

		m.logger.Debug("call alert OFF", zap.String("alertName", alertName), zap.String("scriptName", s.Name))

		return 0
	}
}