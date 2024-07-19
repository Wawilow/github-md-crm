import { error } from '@sveltejs/kit';
import fetch from 'node-fetch';


export function load({ params }) {
	const rep = params.rep;

	if (!rep) throw error(404);

	return {
		rep
	};
}
