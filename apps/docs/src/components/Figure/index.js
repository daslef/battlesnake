import useBaseUrl from "@docusaurus/useBaseUrl";
import React from "react";

export default function Figure({ src, caption, credit }) {
	return (
		<figure style={{ display: "flex", justifyContent: "center" }}>
			<div>
				<img src={useBaseUrl(src)} alt={caption} />
				<figcaption style={{ textAlign: "center" }}>{caption}</figcaption>
				{credit ? (
					<figcaption style={{ textAlign: "center" }}>
						<i>{credit}</i>
					</figcaption>
				) : (
					""
				)}
			</div>
		</figure>
	);
}
