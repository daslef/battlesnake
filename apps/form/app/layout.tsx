import "@/app/globals.css";

export default function RootLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	return (
		<html lang="ru">
			<head>
				<meta charSet="UTF-8" />
				<link rel="icon" type="image/svg+xml" href="/images/favicon.svg" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<title>BattleSnakes / IThub Spb</title>
			</head>
			<body>{children}</body>
		</html>
	);
}
