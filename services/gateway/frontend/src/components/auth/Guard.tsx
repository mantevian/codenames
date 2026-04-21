import { useEffect, useState } from "preact/hooks";
import { useLocation } from "preact-iso";

export default function AuthGuard({ path, children }: { path?: string, children?: any }) {
	const { path: locationPath, route } = useLocation();
	const [isValid, setIsValid] = useState(null);

	useEffect(() => {
		fetch("/api/v1/validate_token", {
			method: "POST",
			headers: {
				"Authorization": `Bearer ${localStorage.getItem("token")}`
			}
		})
			.then(res => res.json())
			.then(json => {
				if (json.success) {
					setIsValid(true);
				} else {
					setIsValid(false);
					route("/");
				}
			})
			.catch(() => {
				setIsValid(false);
				route("/");
			});
	}, [locationPath]);

	if (isValid === null) {
		return <div>Checking authentication...</div>;
	}

	if (isValid) {
		return children;
	}

	return null;
}
