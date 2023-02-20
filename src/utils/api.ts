export const API_URL = import.meta.env.VITE_API_URL;

export function fetchStatus() {
  return fetchJsonApi(`${API_URL}/status`, "GET");
}

export function fetchStatusOfService(serviceId: number, days: number) {
  return fetchJsonApi(`${API_URL}/status/${serviceId}?days=${days}`, "GET");
}

export function fetchMaintenance() {
  return fetchJsonApi(`${API_URL}/maintenance`, "GET");
}

export function fetchJsonApi(url, method) {
  return new Promise((res, rej) => {
    fetch(url, {method: method})
      .then(resp => {
        resp
          .json()
          .then(r => {
            res(r);
          })
          .catch(function (err) {
            rej(err);
          });
      })
      .catch(() => {
        rej(null);
      });
  });
}
