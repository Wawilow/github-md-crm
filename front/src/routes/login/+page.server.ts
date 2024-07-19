import fetch from 'node-fetch';


export const actions = {
	default: async ({ request }) => {
	  const formData = await request.formData();
	  const data = { email: formData.get('email') };
  
	  await fetch('https://mailcoach.app/api/…', {
		method: 'POST',
		body: JSON.stringify(data),
	  });
	},
  };
