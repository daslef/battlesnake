import type {
	FieldValues,
	RegisterOptions,
	UseFormRegister,
} from "react-hook-form";

interface ITextEntry {
	name: string;
	label: string;
	placeholder?: string;
	type: string;
	inputMode?: React.HTMLAttributes<HTMLLIElement>["inputMode"];
	register: UseFormRegister<FieldValues>;
	isPending: boolean;
	options: RegisterOptions<FieldValues, string>;
}

export default function TextEntry({
	name,
	label,
	placeholder,
	type,
	register,
	inputMode,
	options,
	isPending,
}: ITextEntry) {
	return (
		<section className="entry" key={`entry-${name}`}>
			<label htmlFor={name} className="entry__label">
				{label}
			</label>
			{type === "textarea" ? (
				<textarea
					placeholder={placeholder ?? ""}
					id={name}
					key={name}
					className="entry__input entry__input--textarea"
					rows={5}
					autoComplete="off"
					spellCheck={false}
					readOnly={isPending}
					{...register(name, options)}
				/>
			) : (
				<input
					type={type}
					placeholder={placeholder ?? ""}
					id={name}
					inputMode={inputMode}
					key={name}
					className="entry__input"
					autoComplete="off"
					spellCheck={false}
					readOnly={isPending}
					{...register(name, options)}
				/>
			)}
		</section>
	);
}
