drop TABLE answers

CREATE TABLE `answers`
(
    `id` string,
    `answer` UTF8,
    PRIMARY KEY (`id`)
);

CREATE TABLE `themes`
(
    `id` String,
    `theme` UTF8,
    PRIMARY KEY (`id`)
);

CREATE TABLE `questions`
(
    `id` String,
    `question` UTF8,
    `answer` UTF8,
    PRIMARY KEY (`id`)
);