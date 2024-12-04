import { useEffect, useRef } from "react";
import { Alert, AlertDescription, AlertTitle } from "../ui/alert";
import { Terminal } from "lucide-react";

type TerminalAlertsProps = {
    clientID: string;
	status: "loading" | "success" | "error" | undefined;
	onAlertChange: () => void;
};

export const TerminalAlerts = ({
	clientID,
	status,
	onAlertChange,
}: TerminalAlertsProps) => {
	const wrapperRef = useRef<HTMLDivElement>(null);
	useEffect(() => {
		if (!wrapperRef.current) {
			return;
		}
		const observer = new MutationObserver(onAlertChange);
		observer.observe(wrapperRef.current, { childList: true });

		return () => {
			observer.disconnect();
		};
	}, [onAlertChange]);

	return (
		<div ref={wrapperRef} className="absolute top-0 z-10 w-full bg-opacity-100 bg-slate-400">
			{status === "error" ? (
				<ErrorAlert />
			) : null}
		</div>
	);
};

const ErrorAlert = () => {
    return (<Alert className="bg-slate-400 bg-opacity-100">
        <Terminal className="h-4 w-4" />
        <AlertTitle>Heads up!</AlertTitle>
        <AlertDescription>
          You can add components and dependencies to your app using the cli.
        </AlertDescription>
      </Alert>)
}
