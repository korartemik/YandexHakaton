package api

type Intents struct {
	Confirm          *EmptyObj           `json:"YANDEX.CONFIRM"`
	Reject           *EmptyObj           `json:"YANDEX.REJECT"`
	ButtonNext       *EmptyObj           `json:"YANDEX.BOOK.NAVIGATION.NEXT"`
	Next             *EmptyObj           `json:"next_step"`
	Help             *EmptyObj           `json:"help"`
	Stop             *EmptyObj           `json:"stop"`
	Agree            *EmptyObj           `json:"agree_with_answer"`
	Thx              *EmptyObj           `json:"thx"`
	Bye              *EmptyObj           `json:"goodbye"`
	WhatCanYouDo     *EmptyObj           `json:"what_can_you_do"`
	NextQuestion     *EmptyObj           `json:"next_question"`
	WaitAnswer       *EmptyObj           `json:"wait_answer"`
	CreateAnswer     *IntentCreateAnswer `json:"answer_to_the_question"`
	AskBot           *IntentCreateAnswer `json:"ask_a_bot"`
	CompleteAsk      *IntentCreateAnswer `json:"complete_ask"`
	AskMeQuestion    *IntentCreateAnswer `json:"ask_me_a_question"`
	CompleteQuestion *IntentCreateAnswer `json:"complete_the_question"`
	CompleteTheme    *IntentCreateAnswer `json:"complete_theme"`
}

type IntentCreateAnswer struct {
	Slots IntentCreateListSlots `json:"slots"`
}

type IntentCreateListSlots struct {
	ListName *Slot `json:"question"`
}
