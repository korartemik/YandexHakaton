package api

type State string

const (
	StateInit                  State = ""
	StateWaitAnswerFromChatGPT State = "WAIT_ANSWER"
	StateCreateReqName         State = "CREATE_REQ_NAME"
	StateDelReqConfirm         State = "DELETE_REQ_CNFRM"
	StateDelReqName            State = "DELETE_REQ_NAME"
	StateViewReqName           State = "VIEW_REQ_NAME"
	StateAddItemReqItem        State = "ADD_ITM_REQ_ITM"
	StateAddItemReqList        State = "ADD_ITM_REQ_LST"
	StateDelItemReqItem        State = "DEL_ITM_REQ_ITM"
	StateDelItemReqList        State = "DEL_ITM_REQ_CNFRM"
)

type StateData struct {
	State    State
	ListName string
	ItemText string
}

func (s *StateData) GetState() State {
	if s == nil {
		return StateInit
	}
	return s.State
}
