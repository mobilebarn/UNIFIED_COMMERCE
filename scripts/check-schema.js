const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

async function checkSchema() {
  try {
    const response = await fetch('http://localhost:8003/graphql', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        query: `
          {
            _service {
              sdl
            }
          }
        `,
      }),
    });

    const data = await response.json();
    console.log('Service SDL:', data.data._service.sdl);
  } catch (error) {
    console.error('Error:', error);
  }
}

checkSchema();