package api

type Intents struct {
	Confirm      *EmptyObj           `json:"YANDEX.CONFIRM"`
	Reject       *EmptyObj           `json:"YANDEX.REJECT"`
	ButtonNext   *EmptyObj           `json:"YANDEX.BOOK.NAVIGATION.NEXT"`
	Next         *EmptyObj           `json:"next_step"`
	CreateAnswer *IntentCreateAnswer `json:"answer_to_the_question"`
	ListLists    *EmptyObj           `json:"list_lists"`
	Cancel       *EmptyObj           `json:"cancel"`
}

type IntentCreateAnswer struct {
	Slots IntentCreateListSlots `json:"slots"`
}

type IntentCreateListSlots struct {
	ListName *Slot `json:"question"`
}
