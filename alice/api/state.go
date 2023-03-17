package api

type State string

const (
	StateInit                  State = ""
	StateWaitAnswerFromChatGPT State = "WAIT_ANSWER"
	StateCreateReqQuestion     State = "CREATE_REQ_QUESTION"
	StateCreateReqAskBot       State = "CREATE_REQ_ASK"
	StateCreateReqTheme        State = "CREATE_REQ_THEME"
	StateNextQuestion          State = "NEXT_QUESTION"
	StateWaitAnswerFromClient  State = "WAIT_CLIENT_ANSWER"
)

type StateData struct {
	State          State
	AnswerID       string
	NextQuestionID string
	ThemeID        string
}

func (s *StateData) GetState() State {
	if s == nil {
		return StateInit
	}
	return s.State
}
