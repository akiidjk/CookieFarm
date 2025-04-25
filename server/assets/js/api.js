const HOST = 'http://localhost:8080'

export async function checkConfig() {
  try {
    const data = await getConfig();
    if (data.error) {
      return false;
    }
    if (!data.configured) {
      return false;
    }
  } catch (error) {
    console.error(error);
    return false;
  }
  return true;
}

export async function sendConfig(config) {
  const res = await fetch(HOST + "/api/v1/config", {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      config: config,
    }),
  });

  if (!res.ok) {
    return res.status;
  }
}

export async function getConfig() {
  try {
    const res = await fetch(HOST + "/api/v1/config", {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    });
    const data = await res.json();
    return data;
  } catch (error) {
    console.error(error);
    return { "error": error.message };
  }
}
