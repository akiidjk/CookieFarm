export const HOST = 'http://localhost:8080'


interface Config {
  configured: boolean;
  server: {
    submit_flag_checker_time: number;
    host_flagchecker: string;
    team_token: string;
    max_flag_batch_size: number;
    protocol: string;
  };
  client: {
    base_url_server: string;
    submit_flag_server_time: number;
    services: string[] | null;
    range_ip_teams: string;
    format_ip_teams: string;
    my_team_ip: string;
  };
}

export async function checkConfig(): Promise<boolean> {
  const data = await getConfig();

  if (data.error) {
    console.log('Config not configured 1');
    return false;
  }

  if (!data.configured) {
    console.log('Config not configured 3');
    return false;
  }

  return true;
}




export async function sendConfig(config: Config) {

  const token = useCookie('token')
  if (!token) {
    return false;
  }
  console.log(JSON.stringify({
    config: config,
  }))
  const res = await fetch(HOST + "/api/v1/config", {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + token.value,
    },
    body: JSON.stringify({
      config: config,
    }),
  });

  if (!res.ok) {
    console.log('Config not configured 2');
    return false;
  }

}

export async function getConfig() {
  const token = useCookie('token')
  if (!token) {
    return false;
  }

  const res = await fetch(HOST + "/api/v1/config", {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + token.value,
    },
  });

  const data = await res.json();

  return data;
}
