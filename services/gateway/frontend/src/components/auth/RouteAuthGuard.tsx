import { useEffect, useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import { ComponentType } from "preact";
import { Storage } from "../../storage/user";

interface RouteAuthGuardProps {
	component: ComponentType<any>;
	redirectTo?: string;
	[key: string]: any;
}

export default function RouteAuthGUard({
	component: Component,
	redirectTo = "/login",
	...rest
}: RouteAuthGuardProps) {
	const { route } = useLocation();

	useEffect(() => {
		if (!Storage.token) {
			route("/login");
		}
	}, []);

	if (Storage.token) {
		return <Component {...rest} />;
	};
	
	return null;
}
