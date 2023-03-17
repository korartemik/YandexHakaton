package stateful

import (
	aliceapi "awesomeProject1/alice/api"
	"awesomeProject1/errors"
	"context"
)

func (h *Handler) helper(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	answer := "Максимальное время ответа 60 секунд. Если вам приходит ответ, в котором говорится, что вам надо попробовать перезадать ваш вопрос попозже, значит случились какие-то технические проблемы с базой данных." +
		"Если в ответе говорится, что вопрос слишком сложный, то случилась непредвиденная проблема со сторонним сервисом." +
		"Иначе если после 1 минуты вам так и не пришёл ответ, то случилась проблема с нашими внутренними сервисами, которую мы решим в ближайшее время." +
		"Я надеюсь, что все проблемы, которые могут возникнуть решатся, и вы сможете воспользоваться нашим навыком"
	if req.Request.NLU.Intents.Help != nil {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: answer, TTS: answer}}, nil
	}
	return nil, nil
}

func (h *Handler) whatCanIDo(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	if req.Request.NLU.Intents.WhatCanYouDo != nil {
		var buttons []*aliceapi.Button
		buttons = append(buttons, &aliceapi.Button{Title: "Что ты умеешь?"})
		buttons = append(buttons, &aliceapi.Button{Title: "Расскажи про"})
		buttons = append(buttons, &aliceapi.Button{Title: "Задай мне вопрос по теме"})
		buttons = append(buttons, &aliceapi.Button{Title: "Узнай о"})
		answer := "Вы можете попросить меня \"рассказать про что-то\" и я подробно с примерами расскажу про это. " +
			"Также Вы можете попросить меня \"узнать о чём-то\" и я узнаю про это и в короткой форме выдам вам ответ." +
			"И самое главное Вы можете попросить меня задать вам вопросы по вашей теме, и я с радостью предложу вам интересные вопросы с ответами"
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: answer, TTS: answer,
			Buttons: buttons}}, nil
	}
	return nil, nil
}

func (h *Handler) stop(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	if req.Request.NLU.Intents.Stop != nil {
		var buttons []*aliceapi.Button
		buttons = append(buttons, &aliceapi.Button{Title: "Что ты умеешь?"})
		buttons = append(buttons, &aliceapi.Button{Title: "Расскажи про"})
		buttons = append(buttons, &aliceapi.Button{Title: "Задай мне вопрос по теме"})
		buttons = append(buttons, &aliceapi.Button{Title: "Узнай о"})
		stopAnswer := "Ладно ладно, прекращаю... Могу ли я чем-то ещё помочь?"
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: stopAnswer, TTS: stopAnswer, Buttons: buttons}}, nil
	}
	return nil, nil
}

func (h *Handler) agree(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	if req.Request.NLU.Intents.Agree != nil {
		var buttons []*aliceapi.Button
		buttons = append(buttons, &aliceapi.Button{Title: "Расскажи про"})
		buttons = append(buttons, &aliceapi.Button{Title: "Задай мне вопрос по теме"})
		buttons = append(buttons, &aliceapi.Button{Title: "Узнай о"})
		agreeAnswer := "Ну а как ты хотел, это тебе не мозги, а компьютер"
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: agreeAnswer, TTS: agreeAnswer, Buttons: buttons}}, nil
	}
	return nil, nil
}

func (h *Handler) bye(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}

	if req.Request.NLU.Intents.Bye != nil {
		var buttons []*aliceapi.Button
		agreeAnswer := "Надеюсь я смогла помочь узнать тебе, что-то новое. Удачи тебе на экзамене и в других тових делах."
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: agreeAnswer, TTS: agreeAnswer, Buttons: buttons, EndSession: true}}, nil
	}
	return nil, nil
}

func (h *Handler) thx(ctx context.Context, req *aliceapi.Request) (*aliceapi.Response, errors.Err) {
	if req.Request.Type != aliceapi.RequestTypeSimple {
		return nil, nil
	}
	var buttons []*aliceapi.Button
	buttons = append(buttons, &aliceapi.Button{Title: "Расскажи про"})
	buttons = append(buttons, &aliceapi.Button{Title: "Задай мне вопрос по теме"})
	buttons = append(buttons, &aliceapi.Button{Title: "Узнай о"})
	text := "Да ну как бы ну не за что, слово в карман конечно не плоложишь... Но все равно очень приятно"
	tts := "Да ну как бы ну не за что,sil<[150]> слово в карман конечно не положишь sil<[150]>... Но все равно очень приятно"
	if req.Request.NLU.Intents.Thx != nil {
		return &aliceapi.Response{Response: &aliceapi.Resp{Text: text, TTS: tts, Buttons: buttons}}, nil
	}
	return nil, nil
}
