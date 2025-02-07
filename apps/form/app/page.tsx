"use client";

import { type SubmitHandler, useForm } from "react-hook-form";

import MultipleChoice from "../components/MultipleChoice";
import Switch from "../components/Switch";
import TextEntry from "../components/TextEntry";

import { textFields, trackFields } from "../data";

type FormValues = {
	[fieldName: string]: string | boolean;
};

export default function Page() {
	const {
		register,
		watch,
		handleSubmit,
		formState: { isSubmitting, isSubmitSuccessful },
	} = useForm<FormValues>();

	const withTeam = watch("with_team", false);

	const isPending = isSubmitSuccessful || isSubmitting;

	const onSubmit: SubmitHandler<FormValues> = (data) => {
		const { track, ...payload } = data;

		payload.track = Object.entries(track)
			.filter(([_, value]) => value)
			.map(([key, _]) => key.split("_")[1])
			.join(", ");

		fetch("https://api.saint-hubs.tech/hackathon", {
			method: "POST",
			headers: {
				"Content-Type": "application/json;charset=utf-8",
			},
			body: JSON.stringify(payload),
		})
			.then((result) => {
				if (!result.ok) {
					throw new Error("Error Occured");
				}
			})
			.catch((error) => {
				console.log(error);
			});
	};

	return (
		<>
			<header className="header">
				<img
					src="/images/logo.svg"
					alt="IThub Saint-Petersburg Logo"
					className="header__logo"
				/>

				<hgroup className="header__heading heading">
					<h2 className="heading__top">ХАКАТОН</h2>
					<h3 className="heading__bottom">07.02 - 09.02</h3>
				</hgroup>

				<section className="header__description description description--greeting">
					<h4 className="description__subheading subheading subheading--greeting">
						Привет, Хабс!
					</h4>
					<p className="description__text description__text--bottom text description__text--greeting">
						Будем рады видеть тебя на первом событии нового года! Осенний
						хакатон раскрыл новые таланты и подтвердил уровень опытных
						участников. В этот раз мы предлагаем тебе сразу три трека: для
						программистов, дизайнеров и разработчиков игр.
					</p>
				</section>

				<section className="header__description description">
					<h4 className="description__subheading subheading">Game Design</h4>
					<p className="description__text description__text--bottom text">
						Если ты хочешь отдохнуть после хардового Game Jam или просто
						кайфуешь от геймдизайна, то тебе сюда. Именно ты будешь предлагать
						новые механики и режимы для игры.
					</p>
				</section>

				<section className="header__description description">
					<h4 className="description__subheading subheading">Visual & Sound</h4>
					<p className="description__text description__text--bottom text">
						Подойдёт для тех, кому нравится колдовать над красивой картинкой и
						качественным звуком. Ты будешь вносить предложения по улучшению
						визуального и саунд-экспириенса.
					</p>
				</section>

				<section className="header__description description">
					<h4 className="description__subheading subheading">Разработка AI</h4>
					<p className="description__text description__text--bottom text">
						Твоей задачей будет обучить игровую модель, а значит не обойтись без
						логики и алгоритмов! Но не переживай, для участия достаточно
						базового понимания условий, циклов, структур данных и функций.
						Доступные языки для разработки решения: C#, Python, JavaScript.
					</p>
				</section>
			</header>

			<main className="survey">
				<h4 className="survey__subheading subheading">
					Скорее заполняй анкету &lt;3
				</h4>

				<form onSubmit={handleSubmit(onSubmit)} className="survey__form form">
					{textFields.map(
						({ name, label, options, placeholder, type, inputMode }) => (
							<TextEntry
								key={name}
								name={name}
								label={label}
								register={register}
								options={options}
								placeholder={placeholder}
								type={type ?? "text"}
								inputMode={inputMode ?? "text"}
								isPending={isPending}
							/>
						),
					)}

					<MultipleChoice
						name="track"
						legend="В каком треке ты хотел(а) бы участвовать?"
						choices={trackFields}
						register={register}
						isPending={isPending}
					/>

					<section className="disclaimer">
						<p className="description__text disclaimer__description__text">
							<strong>Обрати внимание!</strong>
						</p>
						<p className="description__text disclaimer__description__text">
							Можно участвовать индивидуально либо в команде до трёх человек.
							Всем участникам команды нужно зарегистрироваться.
						</p>

						<Switch
							name="with_team"
							label="Участвуешь в команде?"
							register={register}
							isPending={isPending}
						/>
					</section>

					{withTeam && (
						<TextEntry
							key="team_name"
							name="team_name"
							label="Название команды"
							register={register}
							placeholder=""
							type="text"
							inputMode="text"
							isPending={isPending}
							options={{ required: true }}
						/>
					)}

					<button
						type="submit"
						className="button button--submit"
						disabled={isPending}
					>
						{isSubmitSuccessful ? "Заявка принята!" : "Отправить"}
					</button>
					<p className="form__submit__agreement">
						Отправляя форму, я подтверждаю согласие на обработку персональных
						данных в соответствии с{" "}
						<a
							href="https://ithub.ru/design/images/privacy_policy.pdf"
							target="_blank"
							rel="noreferrer"
						>
							условиями политики конфиденциальности
						</a>
					</p>
				</form>
			</main>
		</>
	);
}
