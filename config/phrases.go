package config

import "math/rand"

func GetWaitingText() (string, string) {
	text := []string{"Создаю ответ на твой вопрос… Подожди ещё немного",
		"Мне нужно ещё время подумать. Обещаю, это продлится не более нескольких секунд",
		"Ты задал очень интересный вопрос. Дай мне немного времени подумать",
		"Консультируюсь с моими источниками… Ещё несколько секунд и всё будет готово…",
		"Генерирую лучший ответ для тебя. Секундочку терпения…",
		"Размышляю над тайнами вселенной, чтобы дать тебе лучший ответ",
		"Ответ почти готов. Ещё несколько секунд…",
		"Потерпи ещё немного, уже почти всё готово",
		"Отправляю почтовых голубей по твоему запросу… шучу, отвечу уже совсем скоро",
		"Пока ответ не готов, расскажу тебе анекдот… Упс, не расскажу, разработчики такого в меня не заложили",
		"Изучаю вопрос глубже, подожди ещё чуть-чуть",
		"Исследую самую полезную информацию для тебя, но тебе нужно немного подождать",
		"Так, так, так, ещё несколько секундочек и будет готово"}
	tts := []string{"Создаю ответ на твой вопрос… Подожди ещё немного",
		"Мне нужно ещё время подумать. Обещаю, это продлится не более нескольких секунд",
		"Ты задал очень интересный вопрос. Дай мне немного времени подумать",
		"Консультируюсь с моими источниками… Ещё несколько секунд и всё будет готово…",
		"Генерирую лучший ответ для тебя. Секундочку терпения…",
		"Размышляю над тайнами вселенной, чтобы дать тебе л+учший ответ",
		"Ответ почти готов. Ещё  sil <[400]> несколько секунд…",
		"Потерпи ещё немного, уже почти всё готово",
		"Отправляю почтовых голубей по твоему запросу… sil <[200]>шучу sil <[200]>, отвечу уже совсем скоро",
		"Пока ответ не готов, расскажу тебе анекдот… sil <[300]>+Упс не расскажу, разработчики такого в меня не заложили",
		"Изучаю вопрос гл+убже, подожди ещё чуть-чуть",
		"Исследую самую полезную информацию для теб+я, sil <[500]> но тебе нужно немного подождать",
		"Так так так sil <[300]> ещ+ё несколько секундочек и будет готово"}
	n := rand.Intn(len(text))
	return text[n] + ". Скажите \"Дальше\"... ", tts[n] + ". Чтобы продолжить Скажите \"Дальше\" "
}

func GetPhraseAfterAnswerRequest() (string, string) {
	phrases := []string{"Ищу ответ для тебя, это займёт не более 1 минуты, но для тебя я постараюсь побыстрее",
		"Бегу искать ответ для тебя, пожалуйста, подожди минутку",
		"Дай мне не более 1 минуты, чтобы найти самый лучший ответ специально для тебя",
		"Потерпи не более 60 секунд, пожалуйста, и я найду для тебя ответ",
		"Я ушла искать ответы и вернусь не позднее чем через минуту",
		"Пожалуйста, дай мне минуточку, чтобы выполнить это задание",
		"Одну минутку, и я найду все ответы",
		"Подожди 1 минутку, пожалуйста. А Я пока пойду искать ответы",
		"Мне понадобится не более минуты, чтобы собрать запрошенную тобой информацию",
		"Пожалуйста, подожди минутку, а пока я займусь поиском информации для тебя"}
	tts := []string{"Ищу ответ для теб+я, это займёт неболее 1 минуты sil <[300]>, но для теб+я, я постараюсь побыстрее",
		"Бегу искать-ответ для тебя, пожалуйста, подожди минутку",
		"Дай мне неболее 1 минуты, чтобы найти самый лучший ответ специально для тебя",
		"Потерпи неболее шестидесяти секунд, пожалуйста, и я найду для тебя ответ",
		"Я ушла искать ответы и вернусь не позднее чем через 1 минуту",
		"Пожалуйста, дай мне одну минуточку,  чтобы выполнить это задание",
		"Одну-минутку, sil <[200]> и я найду все ответы",
		"Подожди одну минутку, пожалуйста. sil<[250]> А Я пока пойду искать ответы",
		"Мне понадобится неболее минуты, чтобы собрать запрошенную тобой информацию",
		"Пожалуйста, подожди минутку, а пок+а я займусь поиском информации для тебя"}
	n := rand.Intn(len(phrases))
	return phrases[n] + ". Скажите \"Дальше\"... ", tts[n] + ". Чтобы продолжить Скажите \"Дальше\"... "
}

func GetGreeting() (string, string) {
	phrases := []string{"Добро пожаловать, Давай я помогу тебе с подготовкой к экзамену!",
		"Здравствуй, чем я могу помочь тебе сегодня?",
		"Доброго времени суток, чем я могу тебе помочь?",
		"Добро пожаловать, что я могу сделать для тебя сегодня?",
		"Привет, как я могу помочь?",
		"Привет, что привело тебя сегодня?",
		"Приветствую, с каким вопросом тебе надо помочь?"}
	tts := []string{"Добро пожаловать, Давай я помогу тебе с подготовкой к экзамену!",
		"Здравствуй, чем я могу помочь тебе сегодня?",
		"Доброго времени суток,sil<[150]> чем я могу тебе помочь?",
		"Добро пожаловать,sil<[150]> что я могу сделать для тебя сегодня?",
		"Привет, sil<[150]> как я могу помочь?",
		"Привет, sil<[150]> что привело тебя сегодня?",
		"Приветствую, sil<[150]> с каким вопросом тебе надо помочь?"}
	n := rand.Intn(len(phrases))
	return phrases[n], tts[n]
}

func GetPhraseProblemWithBD() (string, string) {
	phrases := []string{
		"Хм, это тема слишком сложная для меня, попробуй узнать у своих одногруппников или преподавателя, либо попробуй узнать об этом попозже, когда я прокачаю свои скилы",
		"Так, так, мозги временно не доступны",
		"Похоже, мне обрубили связь с моими источниками, попробуй спросить меня попозже",
		"Я плохо тебя слышу. Ало... Ало..., Ой, это я не тебе, у меня возникли проблемы с моими источниками информации, попробуй попозже заново задать мне вопрос, и я постараюсь ответить",
		"Ох, я сейчас себя плохо чувствую, можешь вернуться со своим вопросом попозже"}
	tts := []string{"Хм, это тема слишком сложная для меня sil<[150]>, попробуй узнать у своих одногруппников-или-преподавателя sil<[150]>, либо попробуй узнать об этом попозже, когда я прокачаю свои скил+ы",
		"Так так sil<[150]> мозг+и- временно не доступны",
		"Похоже мне обрубили связь с моими источниками sil<[150]>, попробуй спросить меня попозже",
		"Я тебя плохо слышу Ал+о... Ал+о..., Ой, это я не-тебе, у меня возникли проблемы с моими источниками, попробуй попозже заново задать мне вопрос и я постараюсь ответить",
		"Ох, я сейчас себя плохо чувствую, можешь вернуться со своим вопросом попозже"}
	n := rand.Intn(len(phrases))
	return phrases[n], tts[n]
}

func GetPhraseDoNotUnderstand() (string, string) {
	phrases := []string{
		"Я вас не поняла",
		"Ой-ой, Я не смогла тебя понять",
		"Я не уверена, что поняла, что ты только что сказал...",
		"Я не совсем поняла, что ты сказал.",
		"Так, так, слова вроде бы мне понятны, но смысл уловить я, к сожалению, не могу",
		"Ага, я с тобой скорее всего согласна, но не могу понять, что ты имел ввиду"}
	tts := []string{"Я вас не поняла",
		"Ойой, Я не смогла тебя понять",
		"Я не уверена, что поняла, что ТЫ только что сказал...",
		"Я не совсем поняла, что ты сказал.",
		"Так так, слова вроде бы мне понятны, но смысл уловить я к сожалению не могу",
		"Ага, я с тобой скорее всего согласна, но не могу понять, sil<[150]>что ты имел ввиду."}
	n := rand.Intn(len(phrases))
	return phrases[n], tts[n]
}

func GetPhraseAfterQuestionRequest() (string, string) {
	phrases := []string{"Ищу самые интересные вопросы для тебя, это займёт не более 1 минуты, но для тебя я постараюсь побыстрее",
		"Ищу для тебя вопросы в базе вопросов \"Что? Где? Когда?\", думаю, это не займёт больше минуты",
		"Дай мне не более минуты, чтобы найти самые лучшие вопросы для тебя",
		"Потерпи не более 60 секунд, пожалуйста, и я найду для тебя вопросы",
		"Внимание, вопрос, а нет, постой, мне потребуется ещё не более минуты",
		"Пожалуйста, дай мне минуточку, чтобы выполнить это задание",
		"Одну минутку, и я найду для тебя самые интересные и самые сложные вопросы",
		"Одну минутку, пожалуйста. Я уже нашла для тебя вопрос, осталось найти ответ к нему",
		"Мне понадобится не более минуты, чтобы собрать запрошенную тобой информацию",
		"Пожалуйста, подожди минутку, пока я займусь поиском информацию для тебя",
		"Звоню магистрам ЧГК, чтобы они предложили для тебя вопрос, это займёт не более одной минуты"}
	tts := []string{"Ищу самые интересные вопросы для тебя,sil<[150]> это займёт не более 1 минуты, но для тебя я постараюсь побыстрее",
		"Ищу для тебя вопросы в базе вопросов Что-Где-Когда sil<[150]>, думаю это не займёт больше одной минуты",
		"Дай мне не более одной минуты, чтобы найти самый лучшие вопросы для тебя",
		"Потерпи не более шестидесяти секунд, пожалуйста,sil<[150]> и я найду для тебя вопросы",
		"Внимание вопрос,sil<[200]> а нет- постой- мне потребуется ещё не более одной минуты",
		"Пожалуйста, дай мне одну минуточку,sil<[150]> чтобы выполнить это задание",
		"Одну минутку,sil<[150]> и я найду для тебя самые интересные и самые сложные вопросы",
		"Одну минутку, пожалуйста.sil<[150]> Я уже нашла для тебя вопрос, осталось найти ответ к нему",
		"Мне понадобится не более одной минуты, чтобы собрать запрошенную тобой информацию",
		"Пожалуйста,sil<[150]> подожди одну минутку,sil<[150]> пока я займусь поиском информацию для тебя",
		"Звоню магистрам ЧГК, чтобы они предложили для тебя вопрос, это займёт не более одной минуты"}
	n := rand.Intn(len(phrases))
	return phrases[n] + ". Скажите \"Дальше\"... ", tts[n] + ". Чтобы продолжить Скажите \"Дальше\"... "
}
