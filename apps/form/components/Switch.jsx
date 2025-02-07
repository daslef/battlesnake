import classes from "./Switch.module.css";

export default function Switch({ name, label, register, isPending }) {
	return (
		<label className={classes.toggle} htmlFor={`switch-${name}`}>
			{label}
			<input
				className={classes.toggle__input}
				type="checkbox"
				id={`switch-${name}`}
				disabled={isPending}
				{...register(name, {})}
			/>
			<div className={classes.toggle__fill} />
		</label>
	);
}
