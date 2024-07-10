import { error } from '@sveltejs/kit';

export function load({ params }) {
    const rep = "123";

    if (!rep) throw error(404);

    return {
        rep
    };
}

