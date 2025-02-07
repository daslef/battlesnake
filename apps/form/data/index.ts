export const trackFields: Array<{ label: string; name: string }> = [
	{ label: "Разработка AI", name: "track_AI" },
	{ label: "Visual & Sound", name: "lang_VisualSound" },
	{ label: "Game Design", name: "lang_GameDesign" },
];

interface ITextField {
	name: string;
	label: string;
	placeholder?: string;
	type?: string;
	inputMode?: React.HTMLAttributes<HTMLLIElement>["inputMode"];
	options: object;
}

export const textFields: Array<ITextField> = [
	{
		name: "name",
		label: "Имя и фамилия",
		type: "text",
		placeholder: "Иван Иванов",
		options: {
			required: true,
			pattern: /[А-ЯЁ][а-яё]+ [А-ЯЁ][а-яё]+.*/i,
		},
	},
	{
		name: "group",
		label: "Полный номер группы",
		placeholder: "РИ-9.24.1, ИТ-11.23, Д-9.24.5",
		type: "text",
		options: { required: true },
	},
	{
		name: "telegram",
		label: "Твой ник в telegram для связи",
		placeholder: "@hubs",
		options: {
			required: true,
		},
	},
];
