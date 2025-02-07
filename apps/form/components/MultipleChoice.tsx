import type { FieldValues, UseFormRegister } from "react-hook-form";

interface IMultipleChoice {
	name: string;
	legend: string;
	choices: Array<{ name: string; label: string }>;
	register: UseFormRegister<FieldValues>;
	isPending: boolean;
}

export default function MultipleChoice({
	name,
	legend,
	choices,
	register,
	isPending,
}: IMultipleChoice) {
	return (
		<fieldset className="fieldset">
			<legend className="fieldset__legend">{legend}</legend>
			{choices.map(({ name: choiceName, label }) => (
				<section
					key={`${name}_${choiceName}`}
					className="fieldset__choice choice"
				>
					<input
						type="checkbox"
						className="choice__input"
						key={`${name}_${choiceName}-input`}
						id={`${name}_${choiceName}`}
						disabled={isPending}
						{...register(`${name}.${choiceName}`, {})}
					/>
					<label htmlFor={`${name}_${choiceName}`} className="choice__label">
						{label}
					</label>
				</section>
			))}
		</fieldset>
	);
}
