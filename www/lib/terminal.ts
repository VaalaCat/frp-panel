export const terminalWebsocketUrl = async (
	clientID: string,
	height: number,
	width: number,
): Promise<string> => {
	const query = new URLSearchParams();
	query.set("height", height.toString());
	query.set("width", width.toString());

	const url = new URL(`${location.protocol}//${location.host}`);
	url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
	if (!url.pathname.endsWith("/")) {
		`${url.pathname}/`;
	}
	url.pathname += `api/v1/pty/${clientID}`;
	url.search = `?${query.toString()}`;

	url.search = `?${query.toString()}`;
	return url.toString();
};
